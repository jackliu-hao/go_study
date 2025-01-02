package failover

import (
	"context"
	"jikeshijian_go/webbook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	// 你的服务商
	svcs []sms.Service
	//连续超时的个数
	cnt int32

	idx int32
	//阈值，连续超时，超过这个数值就会切换
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}

}
func (t TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, phoneNumbers ...string) error {

	//这段Go代码的功能是从两个原子变量 t.idx 和 t.cnt 中分别读取当前的值，并将它们赋值给局部变量 idx 和 cnt。
	//使用 atomic.LoadInt32 确保了读取操作的原子性，避免竞态条件。
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		//切换服务商 , 拿到新的服务商下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		//这段Go代码使用了原子操作 atomic.CompareAndSwapInt32 来确保线程安全地更新变量 t.idx 的值。
		//检查 t.idx 是否等于 idx。
		//如果相等，则将 t.idx 更新为 newIdx，并返回 true。
		//如果不相等，则不做任何修改，并返回 false。
		// 如果两个人同时判断，只会有一个人会执行成功，另一个会失败，也就是走到else分支
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 如果到这里，说明 成功往后移了一位 , 将当前服务商的失败次数置为0
			atomic.StoreInt32(&t.cnt, 0)

		} else {
			//出现并发了，已经有人切换成功了
		}
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, phoneNumbers...)
	switch err {
	case context.DeadlineExceeded:
		//	超时
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		//	连续状态被打断
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		//这里可以考虑切换到下一个
		return err
	}

}
