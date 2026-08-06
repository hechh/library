package msgqueue

import (
	"sync"
	"time"

	"github.com/hechh/library/base/queue"
	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/mlog"
)

type MsgQueue[T ITask] struct {
	*Attribute                 // 基础
	tasks      *queue.Queue[T] // 任务队列
	notifyCh   chan struct{}   // 通知
	exitCh     chan struct{}   // 退出
	updateTime int64           // 更新时间
	lockTime   int64           // 全局锁保活
	wg         sync.WaitGroup  // 等待任务完成
	startWg    sync.WaitGroup  // 启动状态
}

func NewMsgQueue[T ITask](opts ...Option) *MsgQueue[T] {
	attr := new(Attribute)
	for _, opt := range opts {
		opt(attr)
	}
	if attr.id <= 0 {
		attr.id = GenId()
	}
	return &MsgQueue[T]{
		Attribute: attr,
		tasks:     queue.NewQueue[T](),
		notifyCh:  make(chan struct{}, 1),
		exitCh:    make(chan struct{}),
	}
}

func (d *MsgQueue[T]) Start() bool {
	if d.IsStopped() {
		d.startWg.Add(1)
		d.wg.Add(1)
		safe.SafeGo(mlog.Fatalf, d.run)
		d.startWg.Wait()
	}
	return d.IsRunning()
}

func (d *MsgQueue[T]) Stop() {
	if d.IsRunning() {
		close(d.exitCh)
		d.Stopped()
		d.OnDelete()
		d.Waiting()
	}
}

func (d *MsgQueue[T]) Wait() {
	if d.IsWaiting() {
		id := d.GetId()
		d.SetId(0)
		d.wg.Wait()
		mlog.Infof("%s(%d)关闭成功", d.name, id)
		d.Stopped()
	}
}

func (d *MsgQueue[T]) Push(t T) (flag bool) {
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

func (d *MsgQueue[T]) run() {
	defer d.wg.Done()

	// 先抢占锁
	if err := d.OnLock(); err != nil {
		mlog.Errorf("%s抢占全局锁失败", d.name)
		return
	}

	// 启动成功
	d.Running()
	d.startWg.Done()

	// 保活全局锁
	tt := time.NewTicker(time.Second)
	defer func() {
		tt.Stop()    // 停止定时器
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
			if d.idleSecond > 0 && tnow.Unix()-d.updateTime > d.idleSecond {
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

func (d *MsgQueue[T]) handle() {
	for range 100 {
		f, ok := d.tasks.Pop()
		if !ok {
			return
		}
		if f.Do() {
			d.updateTime = time.Now().Unix()
		}
	}
	if d.tasks.GetCount() > 0 {
		select {
		case d.notifyCh <- struct{}{}:
		default:
		}
	}
}
