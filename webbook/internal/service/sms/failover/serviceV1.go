package failover

import (
	"context"
	"errors"
	"jikeshijian_go/webbook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailoverSMSServiceV1 struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSServiceV1(svcs []sms.Service) sms.Service {
	return &FailoverSMSServiceV1{
		svcs: svcs,
	}
}

// Send 轮询，但是不让他从他每次都是从0开始循环
func (f FailoverSMSServiceV1) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	// 取下一个节点作为起始节点
	//这段Go代码的功能是使用原子操作对一个无符号64位整数进行递增，并返回递增后的值。
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, phoneNumbers...)
		switch {
		case err == nil:
			return nil
			//	超时或者取消直接return
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return err
		default:
			// 输出日志
			log.Printf("sms service %s failed,err:%v", svc, err)
		}
	}
	return errors.New("全部服务商都失败了")
}
