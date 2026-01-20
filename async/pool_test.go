package async

import "testing"

func TestPool(t *testing.T) {
	aa := NewAsyncPool(50)
	aa.Start()
	aa.Push(func() {
		t.Log("-----1------")
	})
	aa.Push(func() {
		t.Log("-----2------")
	})
	aa.Stop()
}
