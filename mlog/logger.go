package mlog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/spf13/cast"
)

type Logger struct {
	level   int32
	list    []IWriter
	filters map[string]struct{}
}

func NewLogger(level any, ws ...IWriter) *Logger {
	var logLevel int32
	switch vv := level.(type) {
	case string:
		logLevel = StringToLevel(vv)
	default:
		logLevel = cast.ToInt32(vv)
	}
	return &Logger{
		level:   logLevel,
		list:    ws,
		filters: make(map[string]struct{}),
	}
}

func (d *Logger) Close() {
	for _, w := range d.list {
		w.Close()
	}
}

func (d *Logger) SetLevel(level int32) {
	atomic.StoreInt32(&d.level, level)
}

func (d *Logger) Filter(names ...string) {
	for _, name := range names {
		d.filters[name] = struct{}{}
	}
}

func (d *Logger) Trace(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_TRACE {
		d.output(skip+1, LOG_TRACE, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Debug(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_DEBUG {
		d.output(skip+1, LOG_DEBUG, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Warn(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_WARN {
		d.output(skip+1, LOG_WARN, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Info(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_INFO {
		d.output(skip+1, LOG_INFO, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Error(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_ERROR {
		d.output(skip+1, LOG_ERROR, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) Fatal(skip int, attr, format string, args ...any) {
	if atomic.LoadInt32(&d.level) <= LOG_FATAL {
		d.output(skip+1, LOG_FATAL, attr, fmt.Sprintf(format, args...))
	}
}

func (d *Logger) output(depth int, level int32, attr, msg string) {
	pc, file, _, _ := runtime.Caller(depth + 1)
	if len(attr) <= 0 {
		fname := path.Base(runtime.FuncForPC(pc).Name())
		if pos := strings.LastIndex(fname, "."); pos >= 0 {
			attr = fname[pos+1:]
		}
	}
	if _, ok := d.filters[path.Base(file)]; ok {
		return
	}
	if _, ok := d.filters[attr]; ok {
		return
	}
	meta := Meta{Level: level, Msg: msg}
	/*
		if depth >= 2 {
			meta.FileName = file
			meta.Line = line
			meta.FuncName = fname
		}
	*/
	data := get(len(d.list))
	data.Write(meta)
	for _, w := range d.list {
		w.Push(data)
	}
}
