package util

import (
	"ffCommon/log/log"
	"fmt"
	"sync/atomic"
	"time"
)

// Worker 一个工作者状态管理, 使用方法:
//	进入可使用状态后(Ready), 与Worker绑定的对象才可被使用, EnterWork方法, 才会返回true
//	开始关闭后, 与Worker绑定的对象, 不能立即关闭或释放, 直到所有使用者全部使用完毕, 即WaitWorkEnd返回
//	关闭或释放与Worker绑定的对象, 依赖于WaitWorkEnd返回
//	每一次EnterWork的调用, 都需要对应调用一次LeaveWork
type Worker struct {
	status int32 // 运行状态  0初始状态,不可使用 >0可使用
	work   int32 // 当前有多少对EnterWork调用

	onceClose Once // 一次Reset期间, Close一次
}

//
func (w *Worker) String() string {
	return fmt.Sprintf("status[%v] work[%v]", w.status, w.work)
}

// Reset 重置回初始状态
func (w *Worker) Reset() {
	w.onceClose.Reset()

	atomic.StoreInt32(&w.status, 0)
	atomic.StoreInt32(&w.work, 0)
}

// Ready 可供使用
func (w *Worker) Ready() {
	atomic.StoreInt32(&w.status, 1)
}

// Close 开始关闭流程
func (w *Worker) Close() {
	w.onceClose.Do(func() {
		if !atomic.CompareAndSwapInt32(&w.status, 0, 0) {
			// Ready过
			atomic.AddInt32(&w.status, -1)
		}
	})
}

// EnterWork 进入工作, 返回能否进行工作. 不论能否工作, 都应该执行LeaveWork
func (w *Worker) EnterWork() bool {
	atomic.AddInt32(&w.work, 1)

	if atomic.CompareAndSwapInt32(&w.status, 0, 0) {
		return false
	}

	atomic.AddInt32(&w.status, 1)
	return true
}

// LeaveWork 离开工作. work为EnterWork的返回值
func (w *Worker) LeaveWork(work bool) {
	atomic.AddInt32(&w.work, -1)

	if work {
		atomic.AddInt32(&w.status, -1)
	}
}

// WaitWorkEnd 等待所有使用者使用完毕, 最短等待1秒
func (w *Worker) WaitWorkEnd(maxWaitSecond int) {
	waitSecond := 0
	for {
		if atomic.CompareAndSwapInt32(&w.work, 0, 0) && atomic.CompareAndSwapInt32(&w.status, 0, 0) {
			break
		}

		// 等待1秒
		<-time.After(time.Second)

		waitSecond++
		if waitSecond > maxWaitSecond {
			log.FatalLogger.Printf("Worker.WaitWorkEnd spend too much time maxWaitSecond[%v]: %v",
				maxWaitSecond, w)
			break
		}
	}
}
