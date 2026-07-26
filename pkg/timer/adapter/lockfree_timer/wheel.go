package lockfree_timer

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/hechh/library/base/queue"
	"github.com/hechh/library/pkg/timer/domain"
)

type Wheel struct {
	size    int64                        // bucket数量,必须是2的N次方
	mask    int64                        // size-1
	offset  int64                        // cursor的偏移量
	cursor  int64                        // 当前处理的bucket指针，存储时间戳
	buckets []*queue.Queue[domain.ITask] // bucket数据
}

func NewWheel(size int64, offset int64) *Wheel {
	size = 1 << (bits.Len64(uint64(size)) - 1)
	buckets := make([]*queue.Queue[domain.ITask], size)
	for j := int64(0); j < size; j++ {
		buckets[j] = queue.NewQueue[domain.ITask]()
	}
	return &Wheel{
		offset:  offset,
		size:    size,
		mask:    size - 1,
		buckets: buckets,
	}
}

func (d *Wheel) Refresh(now int64) {
	atomic.StoreInt64(&d.cursor, now)
}

func (d *Wheel) Get(now int64) *queue.Queue[domain.ITask] {
	diff := (now - atomic.LoadInt64(&d.cursor)) >> d.offset
	if diff <= 0 || diff > d.size {
		return nil
	}
	return d.buckets[(now>>d.offset)&d.mask]
}

func (d *Wheel) String() string {
	return fmt.Sprintf("size=%d, mask=%d, offset=%d, cursor=%d", d.size, d.mask, d.offset, d.cursor)
}
