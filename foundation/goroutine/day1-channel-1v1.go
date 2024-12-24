package main

import "time"
import "fmt"

// 第一种用法：用作信号传递
// 1 对 1

// 作为信号量机制
type signal struct{}

func worker() {
	fmt.Println("worker is working")
	time.Sleep(time.Second * 1)
}

// 返回值是channel , 参数需要传递一个函数
func spawn(f func()) <-chan signal {

	c := make(chan signal)

	go func() {
		fmt.Println("worker  start to work")
		f()
		// 将信号量写入channel
		c <- signal{}
	}()
	return c
}

func main1() {
	fmt.Println("start a worker...")
	c := spawn(worker)
	// 阻塞以下主线程，这样才可以等待worker执行完毕
	<-c
	fmt.Println("worker work done!")

}
