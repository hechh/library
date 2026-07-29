package enum

import (
	"github.com/spf13/cast"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IEnum interface {
	Number() protoreflect.EnumNumber
}

func ToInt32(val any) int32 {
	switch vv := val.(type) {
	case int32:
		return vv
	case IEnum:
		return int32(vv.Number())
	default:
		return cast.ToInt32(val)
	}
}

func ToUint32(val any) uint32 {
	switch vv := val.(type) {
	case IEnum:
		return uint32(vv.Number())
	default:
		return cast.ToUint32(val)
	}
}
