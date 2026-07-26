package dbpool

import (
	"sync"
	"time"

	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/pkg/dbpool/domain"
	"github.com/hechh/library/pkg/dbpool/internal/registry"
	"github.com/hechh/library/pkg/mlog"
)

type DbPool struct {
	uidModValue uint64                    // 开始的用户ID
	shardsSize  uint64                    // 分片数量
	newFunc     func() domain.IClient     // new函数
	shards      map[uint32]domain.IClient // 分片数据库连接
	globals     map[string]domain.IClient // 全局数据库连接
	exitCh      chan struct{}             // 关闭通道
	closeOnce   sync.Once                 // 确保 Close 只执行一次
}

func NewDbPool[T domain.IClient](f func() T) *DbPool {
	return &DbPool{
		newFunc: func() domain.IClient { return f() },
		shards:  make(map[uint32]domain.IClient),
		globals: make(map[string]domain.IClient),
		exitCh:  make(chan struct{}),
	}
}

func (d *DbPool) Init(cfg *domain.Config) error {
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

func (d *DbPool) checkClient(cli domain.IClient) {
	if err := cli.Ping(); err != nil {
		mlog.Errorf("数据库连接异常断开: %v", err)
		if err := safe.Retry(3, 3*time.Second, cli.Connect); err != nil {
			mlog.Errorf("数据库重连全部失败, error:%v", err)
		} else {
			mlog.Infof("数据库重连成功")
		}
	}
}

func (d *DbPool) GetByName(name string) domain.IClient {
	if val, ok := d.globals[name]; ok && val.IsAlive() {
		return val
	}
	return nil
}

func (d *DbPool) GetById(id uint32) domain.IClient {
	if val, ok := d.shards[id]; ok && val.IsAlive() {
		return val
	}
	return nil
}

func (d *DbPool) GetByUid(uid uint64) domain.IClient {
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
