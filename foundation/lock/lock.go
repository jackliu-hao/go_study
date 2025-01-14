package lock

import "sync"

// LockDemo 不需要初始化lock
type LockDemo struct {
	// 读写锁
	lock sync.Mutex
}

func (l *LockDemo) PanicDemo() {
	l.lock.Lock()
	// 如果panic不会释放锁
	panic("panic")
	l.lock.Unlock()
}

func (l *LockDemo) DeferDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}

// NoponiterDemo 会报错，因为是值传递，会拷贝一份，此时存在两把锁
func (l LockDemo) NoponiterDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}

type LockDemoV1 struct {
	// 读写锁
	lock *sync.Mutex
}

// NoponiterDemo 不会报错，因为lock 是指针，但是需要初始化lock，因为指针类型默认是null
func (l LockDemoV1) NoponiterDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}
