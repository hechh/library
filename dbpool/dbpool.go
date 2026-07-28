package dbpool

import (
	"sync"
	"time"

	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/dbpool/internal/registry"
	"github.com/hechh/library/mlog"
	"xorm.io/xorm"
)

// DbConfig 数据库分片配置
type DbConfig struct {
	ShardsId uint32              `yaml:"shards_id,omitempty"` // 分片ID
	DbName   string              `yaml:"dbname,omitempty"`    // 数据库名称
	Db       uint32              `yaml:"db,omitempty"`        // 数据库编号
	User     string              `yaml:"user,omitempty"`      // 数据库用户名
	Password string              `yaml:"password,omitempty"`  // 数据库密码
	Ip       string              `yaml:"ip,omitempty"`        // 数据库IP地址
	Port     uint32              `yaml:"port,omitempty"`      // 数据库端口
	Slaves   map[int32]*DbConfig `yaml:"slaves,omitempty"`    // 从库配置列表
}

// Config 数据库连接池配置
type Config struct {
	UidModValue uint64      `yaml:"uid_mod_value,omitempty"` // 用户ID取模值
	ShardsSize  int64       `yaml:"shards_size,omitempty"`   // 分片数量
	Globals     []*DbConfig `yaml:"globals,omitempty"`       // 全局数据库配置列表
	Shards      []*DbConfig `yaml:"shards,omitempty"`        // 分片数据库配置列表
}

type IClient interface {
	Init(*DbConfig, ...any) error
	Close() error
	Connect() error
	Ping() error
	IsAlive() bool
	Engine() *xorm.EngineGroup
	NewSession() *xorm.Session
}

type DbPool struct {
	uidModValue uint64             // 开始的用户ID
	shardsSize  uint64             // 分片数量
	newFunc     func() IClient     // new函数
	shards      map[uint32]IClient // 分片数据库连接
	globals     map[string]IClient // 全局数据库连接
	exitCh      chan struct{}      // 关闭通道
	closeOnce   sync.Once          // 确保 Close 只执行一次
}

func NewDbPool[T IClient](f func() T) *DbPool {
	return &DbPool{
		newFunc: func() IClient { return f() },
		shards:  make(map[uint32]IClient),
		globals: make(map[string]IClient),
		exitCh:  make(chan struct{}),
	}
}

func (d *DbPool) Init(cfg *Config) error {
	// 初始化全局数据库
	for _, dbcfg := range cfg.Globals {
		cli := d.newFunc()
		if err := cli.Init(dbcfg, registry.GetGlobalTables(dbcfg.DbName)...); err != nil {
			d.Close()
			return err
		}
		d.globals[dbcfg.DbName] = cli
	}

	// 初始化分片数据库
	d.uidModValue = cfg.UidModValue
	d.shardsSize = uint64(cfg.ShardsSize)
	for _, dbcfg := range cfg.Shards {
		cli := d.newFunc()
		if err := cli.Init(dbcfg, registry.GetShardsTables()...); err != nil {
			d.Close()
			return err
		}
		d.shards[dbcfg.ShardsId] = cli
	}

	go d.check()
	return nil
}

func (d *DbPool) Close() {
	d.closeOnce.Do(func() {
		close(d.exitCh)
		for name, cli := range d.globals {
			if err := cli.Close(); err != nil {
				mlog.Errorf("关闭全局库[%s]失败: %v", name, err)
			}
		}

		// 关闭分片库
		for id, cli := range d.shards {
			if err := cli.Close(); err != nil {
				mlog.Errorf("关闭分片库[%d]失败: %v", id, err)
			}
		}
	})
}

func (d *DbPool) check() {
	tt := time.NewTicker(30 * time.Second)
	defer tt.Stop()
	for {
		select {
		case <-tt.C:
			for _, cli := range d.shards {
				d.checkClient(cli)
			}
			for _, cli := range d.globals {
				d.checkClient(cli)
			}
		case <-d.exitCh:
			return
		}
	}
}

func (d *DbPool) checkClient(cli IClient) {
	if err := cli.Ping(); err != nil {
		mlog.Errorf("数据库连接异常断开: %v", err)
		if err := safe.Retry(3, 3*time.Second, cli.Connect); err != nil {
			mlog.Errorf("数据库重连全部失败, error:%v", err)
		} else {
			mlog.Infof("数据库重连成功")
		}
	}
}

func (d *DbPool) GetByName(name string) IClient {
	if val, ok := d.globals[name]; ok && val.IsAlive() {
		return val
	}
	return nil
}

func (d *DbPool) GetById(id uint32) IClient {
	if val, ok := d.shards[id]; ok && val.IsAlive() {
		return val
	}
	return nil
}

func (d *DbPool) GetByUid(uid uint64) IClient {
	if val, ok := d.shards[d.GetShardsId(uid)]; ok && val.IsAlive() {
		return val
	}
	return nil
}

func (d *DbPool) GetShardsId(uid uint64) uint32 {
	if d.shardsSize == 0 {
		return 0
	}
	if d.uidModValue == 0 {
		return uint32(uid/d.shardsSize) + 1
	}
	return uint32((uid%d.uidModValue)/d.shardsSize) + 1
}
