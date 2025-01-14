package retryable

import (
	"context"
	"errors"
	"jikeshijian_go/webbook/internal/service/sms"
)

type Service struct {
	svc sms.Service
	// 重试
	retryCount int
}

func (s Service) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, phoneNumbers...)
	cnt := 1
	if err != nil {
		return err
	}
	for err != nil && cnt <= s.retryCount {
		err = s.svc.Send(ctx, tpl, args, phoneNumbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("短信服务功能重试都失败了")
}
