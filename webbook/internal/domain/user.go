package domain

import "time"

// User BO (business object)
type User struct {
	Id        int64
	Email     string
	NickName  string
	Birthday  string
	AboutMe   string
	Password  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
