package main

import "fmt"
import "sync"

type counter struct {
	c chan int
	i int
}

/**
 问题：在这个函数中，如果执行到最后的 return，那这个函数中的 goroutine 是不是也结束了？
 答案：不是。NewCounter 函数中的 goroutine 会在 NewCounter 函数返回后继续运行。这是因为 goroutine 是独立于其创建者的生命周期的，只要程序没有结束，
 		goroutine 就会继续执行其任务。在这个例子中，goroutine 会不断递增计数器的值并通过通道发送新的计数值，直到程序结束或通道被关闭。
*/
// NewCounter 创建一个新的计数器实例。
// 该函数初始化一个计数器对象，并启动一个goroutine来持续增加计数器的值。
// 计数器的值通过一个内部channel进行通信，以实现线程安全的计数操作。
// 返回值: *counter 返回一个指向counter类型的指针，用于后续的计数操作。
func NewCounter() *counter {
    // 创建一个counter实例，初始化一个用于传递计数器值的channel。
	cter := &counter{
		c: make(chan int),
	}

    // 启动一个goroutine，在其中不断递增计数器的值，并通过channel发送新的计数值。
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()

    // 返回初始化后的counter实例指针。
	return cter
}

func (cter *counter) Increase() int {
	return <-cter.c
}

func main() {
	// 创建一个计数器实例
	cter := NewCounter()

	// 初始化等待组，用于同步goroutine
	var wg sync.WaitGroup

	// 启动10个goroutine来增加计数器的值
	for i := 0; i < 10; i++ {
		// 为每个goroutine添加到等待组
		wg.Add(1)
		
		// 启动一个goroutine，传入当前的循环变量i
		go func(i int) {
			// 增加计数器的值并获取当前值
			v := cter.Increase()
			
			// 打印goroutine编号和当前计数器的值
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
			
			// 通知等待组，此goroutine完成任务
			wg.Done()
		}(i)
	}

	// 等待所有goroutine完成任务
	wg.Wait()
}