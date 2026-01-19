package mlog

import (
	"strings"
	"sync"
	"time"
)

const (
	LOG_TRACE = 1
	LOG_DEBUG = 2
	LOG_WARN  = 3
	LOG_INFO  = 4
	LOG_ERROR = 5
	LOG_FATAL = 6
)

var (
	logObj   = NewLogger(LOG_DEBUG, &StdWriter{})
	dataPool = sync.Pool{New: func() any { return NewData() }}
)

type Meta struct {
	FileName string
	Line     int
	FuncName string
	Level    int32
	Msg      string
}

type IData interface {
	Now() time.Time // 获取时间
	Add(int32)      // 并发次数
	Done() int32    // 完成次数
	Write(Meta)     // 写入数据
	Read() []byte   // 读取数据
}

type IWriter interface {
	Push(IData) // 推送日志
	Close()     // 关闭
}

func get(times int) IData {
	obj := dataPool.Get().(IData)
	obj.Add(int32(times))
	return obj
}

func put(obj IData) {
	if obj.Done() == 0 {
		dataPool.Put(obj)
	}
}

func Init(level string, lpath string, lname string) {
	logObj.Close()
	logObj = NewLogger(level, NewLogWriter(lpath, lname), &StdWriter{})
}

func Close() {
	logObj.Close()
}

func Tracef(format string, args ...any) {
	logObj.Trace(1, format, args...)
}

func Debugf(format string, args ...any) {
	logObj.Debug(1, format, args...)
}

func Warnf(format string, args ...any) {
	logObj.Warn(1, format, args...)
}

func Infof(format string, args ...any) {
	logObj.Info(1, format, args...)
}

func Errorf(format string, args ...any) {
	logObj.Error(1, format, args...)
}

func Fatalf(format string, args ...any) {
	logObj.Fatal(1, format, args...)
}

func Trace(skip int, format string, args ...any) {
	logObj.Trace(skip+1, format, args...)
}

func Debug(skip int, format string, args ...any) {
	logObj.Debug(skip+1, format, args...)
}

func Warn(skip int, format string, args ...any) {
	logObj.Warn(skip+1, format, args...)
}

func Info(skip int, format string, args ...any) {
	logObj.Info(skip+1, format, args...)
}

func Error(skip int, format string, args ...any) {
	logObj.Error(skip+1, format, args...)
}

func Fatal(skip int, format string, args ...any) {
	logObj.Fatal(skip+1, format, args...)
}

func LevelToString(level int32) string {
	switch level {
	case LOG_TRACE:
		return "TRACE"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_WARN:
		return "WARN"
	case LOG_INFO:
		return "INFO"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	}
	return ""
}

func StringToLevel(str string) int32 {
	switch strings.ToUpper(str) {
	case "TRACE":
		return LOG_TRACE
	case "DEBUG":
		return LOG_DEBUG
	case "WARN":
		return LOG_WARN
	case "INFO":
		return LOG_INFO
	case "ERROR":
		return LOG_ERROR
	case "FATAL":
		return LOG_FATAL
	}
	return LOG_WARN
}
