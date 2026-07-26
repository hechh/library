package fwatcher

import (
	"fmt"
	"path/filepath"

	"github.com/hechh/library/fwatcher/internal/registry"
	"github.com/hechh/library/fwatcher/internal/watcher"
	"github.com/hechh/library/mlog"
)

type EtcdConfig struct {
	PrefixTopic string   `yaml:"prefix_topic,omitempty"` // etcd 前缀主题
	Endpoints   []string `yaml:"endpoints,omitempty"`    // etcd 节点地址列表
	KeepAlive   int64    `yaml:"keep_alive,omitempty"`   // 保活时间（秒）
}

type Config struct {
	DataPath string      `yaml:"data_path,omitempty"` // 数据文件目录
	XlsxPath string      `yaml:"xlsx_path,omitempty"` // Excel 文件目录
	Ext      string      `yaml:"ext,omitempty"`       // 文件扩展名
	Etcd     *EtcdConfig `yaml:"etcd,omitempty"`      // etcd 配置
}

// 配置同步接口
type ISync interface {
	Init(*Config) error
	Close()
	Put(string, []byte) error
	Update(string, []byte) error
	Delete(string) error
	Watch(func(string, []byte)) error
}

// Fwatcher 文件监听器
type Fwatcher struct {
	local   *watcher.Watcher // 本地监听
	sync    ISync            // 远程同步
	newFunc func() ISync
}

func NewFwatcher[T ISync](f func() T) *Fwatcher {
	return &Fwatcher{
		newFunc: func() ISync {
			if f != nil {
				return f()
			}
			return nil
		},
	}
}

func (d *Fwatcher) Init(cfg *Config) error {
	// 初始化watcher
	d.local = watcher.NewWatcher(cfg.DataPath, cfg.XlsxPath, cfg.Ext)
	if err := d.local.Init(); err != nil {
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
