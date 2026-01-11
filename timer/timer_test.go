package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/hechh/library/async"
	"github.com/hechh/library/mlog"
)

func TestTimer(t *testing.T) {
	async.SetExcept(mlog.Infof)
	timer := NewTimer(4, 5)
	taskId := uint64(123)
	for i := 0; i < 2; i++ {
		err := timer.Register(&taskId, 1*time.Second, -1, func() {
			fmt.Println("-->", i, time.Now().Unix())
		})
		if err != nil {
			t.Log("Register failed", err)
			return
		}
	}
	time.Sleep(4 * time.Second)
	//select {}
}
