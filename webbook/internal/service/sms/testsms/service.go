package testsms

import (
	"context"
	"fmt"
)

type SMSService struct {
	appId *string // appId
}

func NewService(app string) *SMSService {
	return &SMSService{
		appId: &app,
	}
}

func (ss SMSService) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {
	//TODO implement me
	fmt.Println("send sms to ", phoneNumbers)
	fmt.Println("tpl is is ", tpl)
	fmt.Println("args is ", args)
	return nil
}
