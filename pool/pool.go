package pool

import (
	"bytes"
	"crypto/md5"
	"hash"
	"hash/fnv"
	"sync"
)

var (
	bytesPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(nil)
		},
	}
	poolHash = sync.Pool{
		New: func() any {
			return fnv.New64a()
		},
	}
	md5Pool = sync.Pool{
		New: func() any {
			return md5.New()
		},
	}
)

func GetHash64() hash.Hash64 {
	return poolHash.Get().(hash.Hash64)
}

func PutHash64(h hash.Hash64) {
	poolHash.Put(h)
}

func GetMD5() hash.Hash {
	return md5Pool.Get().(hash.Hash)
}

func PutMD5(h hash.Hash) {
	md5Pool.Put(h)
}

func GetBytes() *bytes.Buffer {
	obj := bytesPool.Get().(*bytes.Buffer)
	obj.Reset()
	return obj
}

func PutBytes(obj *bytes.Buffer) {
	bytesPool.Put(obj)
}
