package domain

import "time"

// User BO (business object)
type User struct {
	Id       int64
	Email    string
	NickName string
	Birthday string
	AboutMe  string
	Password string
	Phone    string
	// 组合并不是很好的方式，因为可能还存在顶顶字段，可能存在同名字段
	WechatInfo WechatInfo
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
