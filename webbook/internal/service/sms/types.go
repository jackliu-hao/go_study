package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error
}
