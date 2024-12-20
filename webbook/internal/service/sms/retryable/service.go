package retryable

import (
	"context"
	"jikeshijian_go/webbook/internal/service/sms"
)

type Service struct {
	svc sms.Service
	// 重试
	retryCount int
}

func (s Service) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	//TODO implement me
	err := s.svc.Send(ctx, tpl, args, phoneNumbers...)
	for err != nil && s.retryCount < 10 {
		err = s.svc.Send(ctx, tpl, args, phoneNumbers...)
		s.retryCount++
	}
	return err
}
