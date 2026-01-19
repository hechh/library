package database

import (
	"fmt"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hechh/library/yaml"
	_ "github.com/lib/pq"
)

type Client struct {
	engine     *xorm.EngineGroup
	driverName string
	dsn        []string // 数据库连接字符串
	dbname     string   // 数据库名称
	isAlive    int32    // 连接是否正常
}

func NewClient(driver string, cfg *yaml.DbConfig) *Client {
	dsn := []string{}
	switch driver {
	case MysqlDriver:
		// 主节点
		dsn = append(dsn, fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=3s&parseTime=true&charset=utf8mb4", cfg.User, cfg.Password, cfg.Host, cfg.DbName))
		// 从节点配置
		for _, scfg := range cfg.Slave {
			dsn = append(dsn, fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=3s&parseTime=true&charset=utf8mb4",
				scfg.User,
				scfg.Password,
				scfg.Host,
				scfg.DbName),
			)
		}
	case PostgreSqlDriver:
		// 主节点
		dsn = append(dsn, fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.DbName))
		// 从节点配置
		for _, scfg := range cfg.Slave {
			dsn = append(dsn, fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
				scfg.User,
				scfg.Password,
				scfg.Host,
				scfg.DbName),
			)
		}
	}
	return &Client{driverName: driver, dsn: dsn, dbname: cfg.DbName}
}

func (o *Client) Connect(tables ...interface{}) error {
	eng, err := xorm.NewEngineGroup(o.driverName, o.dsn)
	if err != nil {
		return err
	}
	eng.SetMaxIdleConns(10)
	eng.SetMaxOpenConns(200)
	if len(tables) > 0 {
		if err := eng.Sync2(tables...); err != nil {
			return err
		}
	}

	// 查看连接是否联通
	if err := eng.Ping(); err != nil {
		eng.Close()
		return err
	}
	if o.engine != nil {
		o.engine.Close()
	}
	o.engine = eng
	atomic.StoreInt32(&o.isAlive, 1)
	return nil
}

func (o *Client) Close() {
	o.engine.Close()
}

// 检测连接是否联通
func (o *Client) Ping() error {
	return o.engine.Ping()
}

func (o *Client) IsAlive() bool {
	return atomic.LoadInt32(&o.isAlive) > 0
}

func (o *Client) NewSession() *xorm.Session {
	return o.engine.NewSession()
}

func (o *Client) GetEngine() *xorm.Engine {
	return o.engine.Engine
}
