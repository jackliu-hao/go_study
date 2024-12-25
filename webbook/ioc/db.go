package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jikeshijian_go/webbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	// 初始化db
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}
	// 初始化db
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
