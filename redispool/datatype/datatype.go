package datatype

import (
	"reflect"

	"github.com/hechh/library/base/logic"
	"github.com/hechh/library/base/utils"
	"github.com/hechh/library/redispool"
)

type DataType[R any, I any] struct {
	cliFunc func(I) redispool.IClient
	mask    uint32
	id      uint32
}

func NewDataType[R any, I any](cli func(I) redispool.IClient, masks ...uint32) *DataType[R, I] {
	var zero R
	return &DataType[R, I]{
		cliFunc: cli,
		mask:    logic.Or(masks...),
		id:      utils.GetCrc32(utils.ParseName(reflect.TypeOf(zero))),
	}
}

func (d *DataType[R, I]) UniqueId() uint32 {
	return d.id
}

func (d *DataType[R, I]) GetClient(i I) redispool.IClient {
	return d.cliFunc(i)
}

func (d *DataType[R, I]) GetMask() uint32 {
	return d.mask
}

func (d *DataType[R, I]) Marshal(val redispool.Message) ([]byte, error) {
	return val.MarshalVT()
}

func (d *DataType[R, I]) Unmarshal(body []byte) (redispool.Message, error) {
	val := any(new(R)).(redispool.Message)
	if err := val.UnmarshalVT(body); err != nil {
		return nil, err
	}
	return val, nil
}
