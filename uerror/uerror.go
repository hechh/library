package uerror

import (
	"fmt"
	"path"
	"runtime"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ICode interface {
	Number() protoreflect.EnumNumber
}

type UError struct {
	file  string
	fname string
	line  int
	code  int32
	msg   string
}

func toInt32(code any) int32 {
	switch vv := code.(type) {
	case int32:
		return vv
	case ICode:
		return int32(vv.Number())
	default:
		return -1
	}
}

func New(code any, format string, args ...any) *UError {
	pc, file, line, _ := runtime.Caller(1)
	fname := runtime.FuncForPC(pc).Name()
	return &UError{
		file:  path.Base(file),
		line:  line,
		fname: path.Base(fname),
		code:  toInt32(code),
		msg:   fmt.Sprintf(format, args...),
	}
}

func Err(code any, format string, args ...any) *UError {
	return &UError{
		code: toInt32(code),
		msg:  fmt.Sprintf(format, args...),
	}
}

func Wrap(code any, err error) *UError {
	return &UError{code: toInt32(code), msg: err.Error()}
}

func Turn(code any, err error) *UError {
	if vv, ok := err.(*UError); ok {
		return vv
	}
	return &UError{
		code: toInt32(code),
		msg:  err.Error(),
	}
}

func (ue *UError) Error() string {
	if len(ue.file) <= 0 {
		return ue.msg
		//return fmt.Sprintf("[%d] %s", ue.code, ue.msg)
	}
	return fmt.Sprintf("%s:%d %s %s", ue.file, ue.line, ue.fname, ue.msg)
	//return fmt.Sprintf("%s:%d %s [%d] %s", ue.file, ue.line, ue.fname, ue.code, ue.msg)
}

func (ue *UError) GetFile() string {
	return ue.file
}

func (ue *UError) GetFunc() string {
	return ue.fname
}

func (ue *UError) GetLine() int {
	return ue.line
}

func (ue *UError) GetCode() int32 {
	return ue.code
}

func (ue *UError) GetMsg() string {
	return ue.msg
}
