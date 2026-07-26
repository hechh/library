package timer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hechh/library/timer/domain"
)

type Timer struct {
	inner  domain.ITimer
	newFun func() domain.ITimer
}

func NewTimer[T domain.ITimer](f func() T) *Timer {
	return &Timer{
		newFun: func() domain.ITimer { return f() },
	}
}

func (d *Timer) Register(task domain.ITask) error {
	if d.inner == nil {
		return fmt.Errorf("定时器未初始化")
	}
	return d.inner.Register(task)
}

func (d *Timer) Init(cfg *domain.Config) error {
	obj := d.newFun()
	if err := obj.Init(cfg); err != nil {
		return err
	}
	d.inner = obj
	return nil
}

func (d *Timer) Close() {
	if d.inner != nil {
		d.inner.Close()
	}
}

type Task struct {
	f      func()  // 定时任务
	id     *uint64 // 唯一任务id
	ttlMs  int64   // 有效时长
	times  int32   // 任务执行次数
	expire int64   // 任务执行时间
}

func NewTask(id *uint64, ttl time.Duration, times int32, f func()) *Task {
	ttlMs := int64(ttl / time.Millisecond)
	return &Task{
		f:      f,
		id:     id,
		ttlMs:  ttlMs,
		times:  times,
		expire: ttlMs + time.Now().UnixMilli(),
	}
}

func (d *Task) String() string {
	return fmt.Sprintf("id=%d, ttlMs=%d, times=%d, expire=%d", *d.id, d.ttlMs, d.times, d.expire)
}

func (d *Task) IsEnable() bool {
	return d.id != nil && atomic.LoadUint64(d.id) > 0 && atomic.LoadInt32(&d.times) != 0
}

func (d *Task) GetTTL() int64 {
	return d.ttlMs
}

func (d *Task) GetExpire() int64 {
	return atomic.LoadInt64(&d.expire)
}

func (d *Task) Refresh(now int64) {
	atomic.StoreInt64(&d.expire, d.ttlMs+now)
}

func (d *Task) Call() {
	d.f()
	if atomic.LoadInt32(&d.times) > 0 {
		atomic.AddInt32(&d.times, -1)
	}
}
