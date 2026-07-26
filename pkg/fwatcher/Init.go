package fwatcher

import (
	"fmt"

	"github.com/hechh/library/pkg/fwatcher/internal/registry"
)

var (
	object *Fwatcher
)

func SetObject(oj *Fwatcher) {
	object = oj
}

// Register 注册配置解析函数
func RegisterParser[T any](sheet string, parseFunc func(*T) error) {
	registry.Register(sheet, parseFunc)
}

// RegisterChange 注册配置变更回调函数
func RegisterChange(sheet string, changeFunc func()) {
	registry.RegisterChange(sheet, changeFunc)
}

func Put(key string, msg []byte) error {
	if object == nil {
		return fmt.Errorf("fwatcher未初始化")
	}
	return object.Put(key, msg)
}

func Update(key string, msg []byte) error {
	if object == nil {
		return fmt.Errorf("fwatcher未初始化")
	}
	return object.Put(key, msg)
}

func Delete(key string) error {
	if object == nil {
		return fmt.Errorf("fwatcher未初始化")
	}
	return object.Delete(key)
}
