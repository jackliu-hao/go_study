package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"jikeshijian_go/webbook/internal/repository/dao"
	"jikeshijian_go/webbook/pkg/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {

	type Config struct {
		DSN string `yaml:"DSN"`
	}
	var cfg Config

	err := viper.UnmarshalKey("db.mysql", &cfg)
	if err != nil {
		panic(err)
	}
	//dns := viper.GetString("db.mysql.dsn")
	// 初始化db
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(goormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询
			SlowThreshold: 0,
			Colorful:      true,
			//为true 会显示占位符的方式
			//ParameterizedQueries: true,
			LogLevel: glogger.Info,
		}),
	})
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

type goormLoggerFunc func(msg string, fields ...logger.Field)

func (g goormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{Key: "args", Val: i})
}

// DoSomething 单方法的接口
type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
