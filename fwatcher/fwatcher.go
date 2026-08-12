package fwatcher

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hechh/library/base/fileutil"
	"github.com/hechh/library/fwatcher/internal/registry"
	"github.com/hechh/library/fwatcher/internal/watcher"
	"github.com/hechh/library/mlog"
)

type EtcdConfig struct {
	Prefix    string   `yaml:"prefix,omitempty"`     // etcd 前缀主题
	Endpoints []string `yaml:"endpoints,omitempty"`  // etcd 节点地址列表
	KeepAlive int64    `yaml:"keep_alive,omitempty"` // 保活时间（秒）
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
	// 确保数据目录存在（远程同步写文件与本地监听都依赖它）
	abspath, err := filepath.Abs(cfg.DataPath)
	if err != nil {
		return err
	}
	if err := fileutil.EnsureDir(abspath); err != nil {
		return err
	}

	// 第一步：先初始化远程同步，把远程最新配置同步到本地
	d.sync = d.newFunc()
	if cfg.Etcd != nil && d.sync != nil {
		if err := d.sync.Init(cfg); err != nil {
			return err
		}
		// Watch 内部会先同步拉取一次远程全量配置（阻塞完成），再异步监听后续变更
		if err := d.sync.Watch(func(sheet string, body []byte) {
			// 回调的 sheet 是 etcd 完整 key(prefix/表名)，剥离前缀得到表名
			name := strings.TrimPrefix(sheet, cfg.Etcd.Prefix+"/")
			filename := filepath.Join(abspath, name+cfg.Ext)
			if err := registry.Save(name, filename, body); err != nil {
				mlog.Warnf("收到同步配置，但是保存失败 error=%v", err)
			}
		}); err != nil {
			return err
		}
	}

	// 第二步：再启动本地监听（此时本地文件已是远程同步后的最新版本）
	d.local = watcher.NewWatcher(cfg.DataPath, cfg.XlsxPath, cfg.Ext)
	if err := d.local.Init(); err != nil {
		return err
	}

	return nil
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
