package main

import (
	"fmt"
	"sync"
	"time"
)

func workerN(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

// spawnGroup
// 返回一个channel
func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	//sync.WaitGroup 是 Go 语言标准库中用于等待一组 Goroutine 完成的同步原语。
	//通过 Add 方法增加计数器，
	//通过 Done 方法减少计数器，
	//当计数器为零时，所有阻塞在 Wait 方法上的 Goroutine 将被释放。
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal{}
	}()
	return c
}

func main2() {
	// 启动一组工作线程...
	fmt.Println("start a group of workers...")

	// 创建一个通道用于接收信号，以控制工作线程组的启动和停止
	groupSignal := make(chan signal)

	// 启动一个由workerN组成的工作线程组，每个工作线程将运行5次
	// groupSignal用于同步控制工作线程组的启动和停止
	c := spawnGroup(workerN, 5, groupSignal)

	// 等待5秒钟，模拟准备阶段，如资源初始化等
	// time.Sleep(5 * time.Second)

	// 通知工作线程组开始工作...
	fmt.Println("the group of workers start to work...")

	// 关闭信号通道，通知所有工作线程停止工作
	close(groupSignal)

	// 等待工作线程组完成所有任务
	<-c

	// 所有工作线程完成任务后，打印完成信息
	fmt.Println("the group of workers work done!")

}
