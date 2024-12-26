package main

func main() {

	//  ==== 使用ioc 初始化
	server := InitWebServerIOC()
	server.Run(":8081")

}
