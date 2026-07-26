package fwatcher

import (
	"fmt"
	"path/filepath"

	"github.com/hechh/library/pkg/fwatcher/domain"
	"github.com/hechh/library/pkg/fwatcher/internal/registry"
	"github.com/hechh/library/pkg/fwatcher/internal/watcher"
	"github.com/hechh/library/pkg/mlog"
)

// Fwatcher 文件监听器
type Fwatcher struct {
	local   *watcher.Watcher // 本地监听
	sync    domain.ISync     // 远程同步
	newFunc func() domain.ISync
}

func NewFwatcher[T domain.ISync](f func() T) *Fwatcher {
	return &Fwatcher{
		local: watcher.NewWatcher(),
		newFunc: func() domain.ISync {
			if f != nil {
				return f()
			}
			return nil
		},
	}
}

func (d *Fwatcher) Init(cfg *domain.Config) error {
	// 初始化watcher
	if err := d.local.Init(cfg); err != nil {
		return err
	}

	d.sync = d.newFunc()
	if cfg.Etcd == nil || d.sync == nil {
		return nil
	}

	// 初始化配置同步
	if err := d.sync.Init(cfg); err != nil {
		return err
	}

	abspath, err := filepath.Abs(cfg.DataPath)
	if err != nil {
		return err
	}

	return d.sync.Watch(func(sheet string, body []byte) {
		filename := filepath.Join(abspath, sheet+cfg.Ext)
		if err := registry.Save(sheet, filename, body); err != nil {
			mlog.Errorf("收到到同步配置，但是保存失败 error=%v", err)
		}
	})
}

func (d *Fwatcher) Close() {
	if d.sync != nil {
		d.sync.Close()
	}
	d.local.Close()
}

func (d *Fwatcher) Put(key string, msg []byte) error {
	if d.sync == nil {
		return fmt.Errorf("配置同步服务未初始化")
	}
	return d.sync.Put(key, msg)
}

func (d *Fwatcher) Update(key string, msg []byte) error {
	if d.sync == nil {
		return fmt.Errorf("配置同步服务未初始化")
	}
	return d.sync.Update(key, msg)
}

func (d *Fwatcher) Delete(key string) error {
	if d.sync == nil {
		return fmt.Errorf("配置同步服务未初始化")
	}
	return d.sync.Delete(key)
}
