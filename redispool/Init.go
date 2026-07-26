package redispool

var object *RedisPool

func SetObject(oj *RedisPool) {
	object = oj
}

// GetByName 根据名称获取Redis客户端
func GetByName(name string) IClient {
	if object != nil {
		return object.GetByName(name)
	}
	return nil
}

// GetById 根据shardID获取Redis客户端
func GetById(shardID uint32) IClient {
	if object != nil {
		return object.GetById(shardID)
	}
	return nil
}

// GetByUid 根据用户获取Redis客户端
func GetByUid(uid uint64) IClient {
	if object != nil {
		return object.GetByUid(uid)
	}
	return nil
}

// GetShardsId 根据用户获取shardID
func GetShardsId(uid uint64) uint32 {
	if object != nil {
		return object.GetShardsId(uid)
	}
	return 0
}
