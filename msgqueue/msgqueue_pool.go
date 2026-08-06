package msgqueue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/hechh/library/base/queue"
	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/mlog"
)

type MsgQueuePool[T ITask] struct {
	*Attribute                 // 基础
	tasks      *queue.Queue[T] // 任务队列
	notifyCh   chan struct{}   // 通知
	exitCh     chan struct{}   // 退出
	taskCh     chan T          // 任务抢占队列
	updateTime int64           // 更新时间
	lockTime   int64           // 更新时间
	w1         sync.WaitGroup  // 等待 run goroutine 退出
	w2         sync.WaitGroup  // 等待 run goroutine 退出
	startWg    sync.WaitGroup  // 启动状态
}

func NewMsgQueuePool[T ITask](opts ...Option) *MsgQueuePool[T] {
	attr := new(Attribute)
	for _, opt := range opts {
		opt(attr)
	}
	if attr.id <= 0 {
		attr.id = GenId()
	}
	return &MsgQueuePool[T]{
		Attribute: attr,
		tasks:     queue.NewQueue[T](),
		notifyCh:  make(chan struct{}, 1),
		exitCh:    make(chan struct{}),
	}
}

func (d *MsgQueuePool[T]) Start() bool {
	if d.IsStopped() {
		d.startWg.Add(1)
		d.w1.Add(1)
		safe.SafeGo(mlog.Fatalf, d.run)
		d.startWg.Wait()
	}
	return d.IsRunning()
}

func (d *MsgQueuePool[T]) Stop() {
	if !d.IsStopped() {
		close(d.exitCh)
		d.Stopped()
		d.OnDelete()
		d.Waiting()
	}
}

func (d *MsgQueuePool[T]) Wait() {
	if d.IsWaiting() {
		id := d.GetId()
		d.SetId(0)
		d.w1.Wait()
		close(d.taskCh)
		d.w2.Wait()
		mlog.Infof("%s(%d)关闭成功", d.name, id)
	}
}

func (d *MsgQueuePool[T]) Push(t T) (flag bool) {
	flag = d.IsRunning()
	if flag {
		d.tasks.Push(t, func() {
			//mlog.Tracef("%s任务数量：%d", d.name, d.tasks.GetCount())
			select {
			case d.notifyCh <- struct{}{}:
			default:
			}
		})
	}
	return flag
}

func (d *MsgQueuePool[T]) start() {
	// 启动成功
	d.taskCh = make(chan T, 5*d.GetSize())
	for range d.GetSize() {
		d.w2.Add(1)
		safe.SafeGo(mlog.Fatalf, func() {
			defer d.w2.Done()
			for task := range d.taskCh {
				if task.Do() {
					atomic.StoreInt64(&d.updateTime, time.Now().Unix())
				}
			}
		})
	}
}

func (d *MsgQueuePool[T]) run() {
	defer d.w1.Done()

	// 先抢占锁
	if err := d.OnLock(); err != nil {
		mlog.Errorf("%s抢占全局锁失败", d.name)
		return
	}

	d.start()
	d.Running()
	d.startWg.Done()

	// 保活全局锁
	tt := time.NewTicker(time.Second)
	defer func() {
		tt.Stop()    // 关闭定时器
		d.Stop()     // 发送停止消息
		d.handle()   // 处理剩余请求
		d.OnUnlock() // 释放锁
	}()
	d.updateTime = time.Now().Unix()
	d.lockTime = d.updateTime
	expire := d.lockSecond * 2 / 3

	// 循环处理任务
	for {
		select {
		case <-d.notifyCh:
			d.handle()
		case <-d.exitCh:
			return
		case tnow := <-tt.C:
			if d.idleSecond > 0 && tnow.Unix()-atomic.LoadInt64(&d.updateTime) > d.idleSecond {
				return
			}
			if d.lockSecond > 0 && tnow.Unix()-d.lockTime >= expire {
				if err := d.OnLock(); err != nil {
					return
				}
				d.lockTime = tnow.Unix()
			}
		}
	}
}

func (d *MsgQueuePool[T]) handle() {
	for range 100 {
		f, ok := d.tasks.Pop()
		if !ok {
			return
		}
		d.taskCh <- f
	}
	if d.tasks.GetCount() > 0 {
		select {
		case d.notifyCh <- struct{}{}:
		default:
		}
	}
}
