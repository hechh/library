package msgqueue

import (
	"sync/atomic"
	"time"
)

const (
	STOPPED_STATUS = 0 // 已停止
	WAITING_STATUS = 1 // 等待启动中
	RUNNING_STATUS = 2 // 运行中
)

type ITask interface {
	Do() bool
}

type Attribute struct {
	name       string                            // 名字
	id         uint64                            // 唯一id
	status     int32                             // 状态
	size       int                               // 协程池大小
	idleSecond int64                             // 闲置时间（秒）
	deleteFunc func(uint64)                      // 删除函数
	lockSecond int64                             // lock有效时长
	lockFunc   func(uint64, time.Duration) error // 全局任务锁
	unlockFunc func(uint64) error                // 全局任务解锁函数
}

var (
	msgqueueId uint64
)

func GenId() uint64 {
	return atomic.AddUint64(&msgqueueId, 1)
}

func (d *Attribute) GetIdleTime() int64 {
	return d.idleSecond
}

func (d *Attribute) GetSize() int {
	return d.size
}

func (d *Attribute) IsRunning() bool {
	return atomic.LoadInt32(&d.status) == RUNNING_STATUS
}

func (d *Attribute) Running() {
	atomic.StoreInt32(&d.status, RUNNING_STATUS)
}

func (d *Attribute) IsWaiting() bool {
	return atomic.LoadInt32(&d.status) == WAITING_STATUS
}

func (d *Attribute) Waiting() {
	atomic.StoreInt32(&d.status, WAITING_STATUS)
}

func (d *Attribute) IsStopped() bool {
	return atomic.LoadInt32(&d.status) == STOPPED_STATUS
}

func (d *Attribute) Stopped() {
	atomic.StoreInt32(&d.status, STOPPED_STATUS)
}

func (d *Attribute) GetName() string {
	return d.name
}

func (d *Attribute) GetIdPointer() *uint64 {
	return &d.id
}

func (d *Attribute) GetId() uint64 {
	return atomic.LoadUint64(&d.id)
}

func (d *Attribute) SetId(val uint64) {
	atomic.StoreUint64(&d.id, val)
}

func (d *Attribute) OnLock() error {
	if d.lockFunc != nil {
		return d.lockFunc(d.GetId(), time.Duration(d.lockSecond)*time.Second)
	}
	return nil
}

func (d *Attribute) OnUnlock() error {
	if d.unlockFunc != nil {
		return d.unlockFunc(d.GetId())
	}
	return nil
}

func (d *Attribute) OnDelete() {
	if d.deleteFunc != nil {
		d.deleteFunc(d.GetId())
	}
}

func (d *Attribute) ToOptions() (rets []Option) {
	rets = append(
		rets,
		WithName(d.GetName()),
		WithId(d.GetId()),
		WithSize(d.GetSize()),
		WithIdleTime(d.GetIdleTime()),
		WithDeleter(d.deleteFunc),
		WithLocker(d.lockSecond, d.lockFunc, d.unlockFunc),
	)
	return
}

type Option func(*Attribute)

func WithName(name string) Option {
	return func(opt *Attribute) {
		opt.name = name
	}
}

func WithId(id uint64) Option {
	return func(opt *Attribute) {
		opt.id = id
	}
}

func WithSize(size int) Option {
	return func(opt *Attribute) {
		opt.size = size
	}
}

func WithIdleTime(idle int64) Option {
	return func(opt *Attribute) {
		opt.idleSecond = idle
	}
}

func WithDeleter(f func(uint64)) Option {
	return func(opt *Attribute) {
		opt.deleteFunc = f
	}
}

func WithLocker(expire int64, lock func(uint64, time.Duration) error, unlock func(uint64) error) Option {
	return func(opt *Attribute) {
		opt.lockSecond = expire
		opt.lockFunc = lock
		opt.unlockFunc = unlock
	}
}
