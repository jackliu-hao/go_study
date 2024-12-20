package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SMSService struct {
	appId     *string // appId
	signature *string // 签名
	client    *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *SMSService {
	return &SMSService{
		appId:     &appId,
		signature: &signName,
		client:    client,
	}
}

func (ss SMSService) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = ss.appId
	req.SignName = ss.signature
	req.TemplateId = &tpl
	// 需要把可变数据转成切片
	req.PhoneNumberSet = ss.toPtrSlice(phoneNumbers)
	req.TemplateParamSet = ss.toPtrSlice(args)
	sendSms, err := ss.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range sendSms.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("发送短信失败: %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *SMSService) toPtrSlice(data []string) []*string {
	// 使用Map函数对字符串切片进行转换，返回一个指向原字符串的指针切片
	// 该操作不会创建新的字符串实例，仅将原有字符串的引用映射到新的切片中
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}
