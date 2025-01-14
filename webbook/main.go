package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	InitLogger()
	InitViper()
	//  ==== 使用ioc 初始化
	server := InitWebServerIOC()
	server.Run(":8081")

}

// InitLogger 初始化日志
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 如果不使用replace，那什么都不会打印出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("日志完成")
}

// 参数指定
func InitViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigType(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func InitViper() {
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go .yaml 之类的后缀
	viper.SetConfigName("dev.yaml")
	// 告诉viper使用的配置文件是哪种
	viper.SetConfigType("yaml")
	// 当前工作目录下的config子目录 , 工作目录是最上级
	// 可以有多个配置文件路径
	viper.AddConfigPath("./config")

	// 读取配置到viper中
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
