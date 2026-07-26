package dbpool

import (
	"github.com/hechh/library/dbpool/domain"
	"github.com/hechh/library/dbpool/internal/registry"
)

var (
	object *DbPool
)

func SetObject(oj *DbPool) {
	object = oj
}

func RegisterGlobal(name string, datas ...any) {
	registry.RegisterGlobal(name, datas...)
}

func RegisterShards(datas ...any) {
	registry.RegisterShards(datas...)
}

func GetByName(name string) domain.IClient {
	if object != nil {
		return object.GetByName(name)
	}
	return nil
}

func GetById(id uint32) domain.IClient {
	if object != nil {
		return object.GetById(id)
	}
	return nil
}

func GetByUid(uid uint64) domain.IClient {
	if object != nil {
		return object.GetByUid(uid)
	}
	return nil
}
