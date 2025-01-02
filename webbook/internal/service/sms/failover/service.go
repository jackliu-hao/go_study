package failover

import (
	"context"
	"errors"
	"jikeshijian_go/webbook/internal/service/sms"
	"log"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

// 轮询
func (f FailoverSMSService) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, phoneNumbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 输出日志，做好监控
		log.Println("sms service send error:", err)
	}
	// 如果这里都崩了，那说明是自己的问题，不是别人的问题
	return errors.New("全部服务商都发送失败了")
}
