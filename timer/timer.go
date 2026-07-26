package timer

import (
	"fmt"
)

type Config struct {
	Size               int   `yaml:"size"`                  // 时间轮数量
	MinPeriodBitNumber int64 `yaml:"min_period_bit_number"` // 定时器最小周期
}

type ITimer interface {
	Init(*Config) error
	Close()
	Register(ITask) error
}

type Timer struct {
	inner  ITimer
	newFun func() ITimer
}

func NewTimer[T ITimer](f func() T) *Timer {
	return &Timer{
		newFun: func() ITimer { return f() },
	}
}

func (d *Timer) Register(task ITask) error {
	if d.inner == nil {
		return fmt.Errorf("定时器未初始化")
	}
	return d.inner.Register(task)
}

func (d *Timer) Init(cfg *Config) error {
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
