package async

import (
	"sync"
	"sync/atomic"
)

type AsyncPool struct {
	sync.WaitGroup
	size   int            // 协程数量
	id     uint64         // 唯一id
	status int32          // 状态
	queue  *Queue[func()] // 任务队列
	list   chan func()    // 任务抢占队列
	notify chan struct{}  // 通知
	exit   chan struct{}  // 退出
}

func NewAsyncPool(size int) *AsyncPool {
	return &AsyncPool{
		size:   size,
		queue:  NewQueue[func()](),
		list:   make(chan func(), 50),
		notify: make(chan struct{}, 1),
		exit:   make(chan struct{}),
	}
}

func (d *AsyncPool) GetIdPointer() *uint64 {
	return &d.id
}

func (d *AsyncPool) GetId() uint64 {
	return atomic.LoadUint64(&d.id)
}

func (d *AsyncPool) SetId(id uint64) {
	atomic.StoreUint64(&d.id, id)
}

func (d *AsyncPool) Start() {
	if atomic.CompareAndSwapInt32(&d.status, 0, 1) {
		d.Add(1)
		go d.run()
		for i := 0; i < d.size; i++ {
			go d.handle()
		}
	}
}

func (d *AsyncPool) Stop() {
	if atomic.CompareAndSwapInt32(&d.status, 1, 0) {
		close(d.exit)
		d.Wait()
		atomic.StoreUint64(&d.id, 0)
		d.Add(d.size)
		close(d.list)
		d.Wait()
	}
}

func (d *AsyncPool) Push(f func()) {
	if atomic.CompareAndSwapInt32(&d.status, 1, 1) {
		d.queue.Push(f)
		select {
		case d.notify <- struct{}{}:
		default:
		}
	}
}

func (d *AsyncPool) handle() {
	for f := range d.list {
		Recover(f)
	}
	d.Done()
}

func (d *AsyncPool) run() {
	defer func() {
		for f := d.queue.Pop(); f != nil; f = d.queue.Pop() {
			d.list <- f
		}
		d.Done()
	}()
	for {
		select {
		case <-d.notify:
			for f := d.queue.Pop(); f != nil; f = d.queue.Pop() {
				d.list <- f
			}
		case <-d.exit:
			return
		}
	}
}
