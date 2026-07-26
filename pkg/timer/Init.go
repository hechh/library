package timer

import (
	"fmt"

	"github.com/hechh/library/pkg/timer/domain"
)

var (
	object *Timer
)

func SetObject(oj *Timer) {
	object = oj
}

// 注册定时任务
func Register(task domain.ITask) error {
	if object != nil {
		return object.Register(task)
	}
	return fmt.Errorf("timer not initialized")
}

func Close() {
	if object != nil {
		object.Close()
	}
}
