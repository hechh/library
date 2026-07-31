package redispool

// DbConfig 数据库分片配置
type DbConfig struct {
	ShardsId uint32              `yaml:"shards_id,omitempty"` // 分片ID
	DbName   string              `yaml:"dbname,omitempty"`    // 数据库名称
	Db       uint32              `yaml:"db,omitempty"`        // 数据库编号
	User     string              `yaml:"user,omitempty"`      // 数据库用户名
	Password string              `yaml:"password,omitempty"`  // 数据库密码
	Ip       string              `yaml:"ip,omitempty"`        // 数据库IP地址
	Port     uint32              `yaml:"port,omitempty"`      // 数据库端口
	Prefix   string              `yaml:"prefix,omitempty"`    // key 前缀
	Slaves   map[int32]*DbConfig `yaml:"slaves,omitempty"`    // 从库配置列表
}

type Config struct {
	UidModValue uint64      `yaml:"uid_mod_value,omitempty"` // 用户ID取模值
	ShardsSize  int64       `yaml:"shards_size,omitempty"`   // 分片数量
	Globals     []*DbConfig `yaml:"globals,omitempty"`       // 全局 Redis 配置列表
	Shards      []*DbConfig `yaml:"shards,omitempty"`        // 分片 Redis 配置列表
}

type RedisPool struct {
	uidModValue uint64             // 开始 Uid 值，用于路由
	shardsSize  uint64             // 分片数量
	newFunc     func() IClient     // new函数
	globals     map[string]IClient // 全局数据库连接池
	shards      map[uint32]IClient // 分片数据库连接池
}

func NewRedisPool[T IClient](f func() T) *RedisPool {
	return &RedisPool{
		newFunc: func() IClient { return f() },
		globals: make(map[string]IClient),
		shards:  make(map[uint32]IClient),
	}
}

func (d *RedisPool) Init(cfg *Config) error {
	// 初始化全局数据库
	for _, dbCfg := range cfg.Globals {
		cli := d.newFunc()
		if err := cli.Init(dbCfg); err != nil {
			d.Close()
			return err
		}
		d.globals[dbCfg.DbName] = cli
	}

	d.uidModValue = cfg.UidModValue
	d.shardsSize = uint64(cfg.ShardsSize)

	// 初始化分片数据库
	for _, dbCfg := range cfg.Shards {
		cli := d.newFunc()
		if err := cli.Init(dbCfg); err != nil {
			d.Close()
			return err
		}
		d.shards[dbCfg.ShardsId] = cli
	}
	return nil
}

func (d *RedisPool) Close() {
	for _, cli := range d.globals {
		cli.Close()
	}
	for _, cli := range d.shards {
		cli.Close()
	}
}

func (d *RedisPool) GetByName(name string) IClient {
	return d.globals[name]
}

func (d *RedisPool) GetById(id uint32) IClient {
	return d.shards[id]
}

func (d *RedisPool) GetByUid(uid uint64) IClient {
	return d.shards[d.GetShardsId(uid)]
}

func (d *RedisPool) GetShardsId(uid uint64) uint32 {
	if d.uidModValue == 0 {
		return uint32(uid/d.shardsSize) + 1
	}
	return uint32((uid%d.uidModValue)/d.shardsSize) + 1
}
