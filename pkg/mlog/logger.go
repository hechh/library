package mlog

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hechh/library/base/datetime"
)

const (
	CACHE_FLUSH_INTERVAL_MS = 500         // 缓存刷新间隔
	CACHE_BUFFER_SIZE       = 1024 * 1024 // 缓存大小
)

type Config struct {
	Mode   string `yaml:"mode,omitempty"`   // 日志模式（debug/release）
	Path   string `yaml:"path,omitempty"`   // 日志文件路径
	Level  string `yaml:"level,omitempty"`  // 日志级别
	Format string `yaml:"format,omitempty"` // 日志格式（text/json）
	Name   string `yaml:"name,omitempty"`   // 日志名称
}

type Logger struct {
	level    atomic.Int32
	writer   atomic.Pointer[GroupWriter]
	dataPool sync.Pool
	caller   atomic.Bool
	format   atomic.Int32
}

// NewLogger 创建新的日志记录器
func NewLogger() *Logger {
	ret := &Logger{
		dataPool: sync.Pool{
			New: func() any {
				buf := &bytes.Buffer{}
				buf.Grow(256)
				return buf
			},
		},
	}
	ret.level.Store(LOG_TRACE)
	ret.writer.Store(&GroupWriter{
		list: []IWriter{
			&StdoutWriter{},
		},
	})
	return ret
}

func (l *Logger) Init(cfg *Config) error {
	l.level.Store(Name2Level(cfg.Level))
	l.format.Store(Name2Format(cfg.Format))
	switch strings.ToLower(cfg.Mode) {
	case "debug":
		l.EnableCaller()
		l.writer.Store(&GroupWriter{
			list: []IWriter{&StdoutWriter{}},
		})
	case "release":
		l.writer.Store(&GroupWriter{
			list: []IWriter{
				NewRotateWriter(cfg.Path, cfg.Name, CACHE_BUFFER_SIZE, time.Duration(CACHE_FLUSH_INTERVAL_MS)*time.Millisecond, RollingByHour),
			},
		})
	default:
		l.EnableCaller()
		l.writer.Store(&GroupWriter{
			list: []IWriter{
				&StdoutWriter{},
				NewRotateWriter(cfg.Path, cfg.Name, CACHE_BUFFER_SIZE, time.Duration(CACHE_FLUSH_INTERVAL_MS)*time.Millisecond, RollingByDay),
			},
		})
	}
	return nil
}

func (l *Logger) Close() {
	l.writer.Load().Close()
}

func (l *Logger) SetLevel(level int32) {
	l.level.Store(level)
}

func (l *Logger) GetLevel() int32 {
	return l.level.Load()
}

func (l *Logger) SetFormat(f int32) {
	l.format.Store(f)
}

func (l *Logger) DisableCaller() {
	l.caller.Store(false)
}

func (l *Logger) EnableCaller() {
	l.caller.Store(true)
}

func (l *Logger) get() *bytes.Buffer {
	buf := l.dataPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (l *Logger) put(buf *bytes.Buffer) {
	l.dataPool.Put(buf)
}

func (l *Logger) Output(skip int, level int32, args ...any) {
	if l.level.Load() > level {
		return
	}
	switch l.format.Load() {
	case FORMAT_JSON:
		l.outputj(skip+1, level, args...)
	default:
		l.output(skip+1, level, args...)
	}
}

func (l *Logger) Outputf(skip int, level int32, format string, args ...any) {
	if l.level.Load() > level {
		return
	}
	switch l.format.Load() {
	case FORMAT_JSON:
		l.outputjf(skip+1, level, format, args...)
	default:
		l.outputf(skip+1, level, format, args...)
	}
}

// writeHeader 写入日志公共头部（时间戳、级别、标签、调用者信息），返回 buffer 和当前时间
func (l *Logger) writeHeader(skip int, level int32) (buff *bytes.Buffer, now time.Time) {
	buff = l.get()
	now = datetime.Now()
	buff.WriteByte('[')
	buff.Write(now.AppendFormat(nil, "2006-01-02 15:04:05.000"))
	buff.WriteByte(']')
	buff.WriteByte(' ')
	buff.WriteByte('[')
	buff.WriteString(Level2Name(level))
	buff.WriteByte(']')

	// 获取调用信息（生产环境可跳过，避免栈回溯开销）
	if l.caller.Load() {
		if _, file, line, ok := runtime.Caller(skip); ok {
			buff.WriteByte(' ')
			buff.WriteString(filepath.Base(file))
			buff.WriteByte(':')
			buff.Write(strconv.AppendInt(nil, int64(line), 10))
			// buff.WriteByte(' ')
			// buff.WriteString(filepath.Base(runtime.FuncForPC(pc).Name()))
		}
	}
	return
}

// Output 输出可变参数日志（JSON 序列化）
func (l *Logger) output(skip int, level int32, args ...any) {
	buff, now := l.writeHeader(skip+1, level)
	defer l.put(buff)

	buff.WriteByte(' ')
	for i, arg := range args {
		if i > 0 {
			buff.WriteByte('|')
		}
		switch varg := arg.(type) {
		case string:
			buff.WriteString(varg)
		case []byte:
			buff.WriteString(hex.EncodeToString(varg))
		default:
			body, _ := sonic.ConfigFastest.Marshal(arg)
			buff.Write(body)
		}
	}
	buff.WriteByte('\n')
	if ww := l.writer.Load(); ww != nil {
		ww.Write(now, buff.Bytes())
	}
}

// Outputf 输出格式化日志
func (l *Logger) outputf(skip int, level int32, format string, args ...any) {
	buff, now := l.writeHeader(skip+1, level)
	defer l.put(buff)

	buff.WriteByte(' ')
	fmt.Fprintf(buff, format, args...)
	buff.WriteByte('\n')

	if ww := l.writer.Load(); ww != nil {
		ww.Write(now, buff.Bytes())
	}
}

// jsonLogEntry JSON 日志条目结构（sonic 零分配序列化）
type jsonLogEntry struct {
	TS    string `json:"ts"`
	Level string `json:"level"`
	File  string `json:"file,omitempty"`
	Tag   string `json:"tag"`
	Msg   any    `json:"msg"`
}

// outputJSON 输出 JSON 格式日志。msgOrArgs 可以是格式化字符串或多个参数。
func (l *Logger) outputj(skip int, level int32, msgOrArgs ...any) {
	buff := l.get()
	defer l.put(buff)
	now := datetime.Now()
	entry := jsonLogEntry{
		TS:    now.Format("2006-01-02T15:04:05.000Z07:00"),
		Level: Level2Name(level),
		Msg:   msgOrArgs,
	}
	if l.caller.Load() {
		if _, file, line, ok := runtime.Caller(skip); ok {
			entry.File = filepath.Base(file) + ":" + strconv.Itoa(line)
		}
	}
	sonic.ConfigDefault.NewEncoder(buff).Encode(entry)
	if ww := l.writer.Load(); ww != nil {
		ww.Write(now, buff.Bytes())
	}
}

// outputJSON 输出 JSON 格式日志。msgOrArgs 可以是格式化字符串或多个参数。
func (l *Logger) outputjf(skip int, level int32, format string, msgOrArgs ...any) {
	buff := l.get()
	defer l.put(buff)
	fmt.Fprintf(buff, format, msgOrArgs...)
	now := datetime.Now()
	entry := jsonLogEntry{
		TS:    now.Format("2006-01-02T15:04:05.000Z07:00"),
		Level: Level2Name(level),
		Msg:   buff.String(),
	}
	if l.caller.Load() {
		if _, file, line, ok := runtime.Caller(skip); ok {
			entry.File = filepath.Base(file) + ":" + strconv.Itoa(line)
		}
	}
	buff.Reset()
	sonic.ConfigDefault.NewEncoder(buff).Encode(entry)
	if ww := l.writer.Load(); ww != nil {
		ww.Write(now, buff.Bytes())
	}
}
