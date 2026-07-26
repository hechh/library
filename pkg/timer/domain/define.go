package domain

type ITask interface {
	IsEnable() bool   // 是否可用
	GetTTL() int64    // 获取有效时长
	GetExpire() int64 // 任务触发时间点
	Refresh(int64)    // 刷新任务触发时间点
	Call()            // 执行任务
	String() string   // 格式化
}

type ITimer interface {
	Init(*Config) error
	Close()
	Register(ITask) error
}

type Config struct {
	Size               int   `yaml:"size"`                  // 时间轮数量
	MinPeriodBitNumber int64 `yaml:"min_period_bit_number"` // 定时器最小周期
}
