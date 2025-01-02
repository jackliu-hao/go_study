package ratelimit

import (
	"context"
	"fmt"
	"jikeshijian_go/webbook/internal/service/sms"
	"jikeshijian_go/webbook/pkg/ratelimit"
)

const (
	RateLimitKey = "tencent:sms"
)

type RatelimitSmsService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSmsService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSmsService) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	// 在这里添加一些代码

	limited, err := s.limiter.Limited(ctx, RateLimitKey)

	if err != nil {
		//可以限流： 如果下游业务很坑，不限流直接把下游服务打崩，属于保守策略。
		//不限流： 下游服务很强，并且业务可用性要求比较高，需要尽量容错
		return fmt.Errorf("短信服务判断是否限流出现错误,%w", err)
	}
	if limited {
		return fmt.Errorf("短信服务限流了")
	}

	err = s.svc.Send(ctx, tpl, args, phoneNumbers...)

	// 在这里也添加一些代码
	return err

}
