package service

import (
	"context"
	"fmt"
	"jikeshijian_go/webbook/internal/repository"
	"jikeshijian_go/webbook/internal/service/sms"
	"math/rand"
)

var (
	ErrSetCodeTooManyTimes    = repository.ErrSetCodeTooManyTimes
	ErrCodeVerifyFailed       = repository.ErrCodeVerifyFailed
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
	// 模板ID
	tplId string
}

// NewCodeService 构造函数
func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service, tplId string) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
		tplId:  tplId,
	}
}

// Send 发送验证码
func (svc *CodeService) Send(ctx context.Context,
	// 区别使用业务场景
	biz string,
	phone string) error {
	// 生成验证码 塞进redis
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, svc.tplId, []string{code}, phone)
	if err != nil {
		// 这里怎么半？
		// 意味着， redis存在这个验证码，但是没有发送成功，用户收不到
		// 能否删除掉这个验证码？
		return err
	}
	return err
}

func (svc *CodeService) generateCode() string {
	// 0-999999
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}

// Verify 验证验证码
func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {

	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
