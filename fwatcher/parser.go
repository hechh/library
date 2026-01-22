package fwatcher

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"sync"
)

type ParseFunc func([]byte) error // 配置解析函数
type ChangeFunc func()            // 配置变更函数

var (
	md5Pool = sync.Pool{
		New: func() any {
			return md5.New()
		},
	}
)

func get() hash.Hash {
	return md5Pool.Get().(hash.Hash)
}

func put(h hash.Hash) {
	md5Pool.Put(h)
}

type Parser struct {
	sheet   string
	value   string
	parse   ParseFunc
	changes []ChangeFunc
}

func NewParser(sheet string) *Parser {
	return &Parser{sheet: sheet}
}

func (d *Parser) Sheet() string {
	return d.sheet
}

func (d *Parser) Register(p ParseFunc) {
	d.parse = p
}

func (d *Parser) Listen(fs ...ChangeFunc) {
	d.changes = append(d.changes, fs...)
}

// 加载配置
func (d *Parser) Parse(isload bool, buf []byte) error {
	hh := get()
	defer put(hh)

	hh.Reset()
	hh.Write(buf)
	value := hex.EncodeToString(hh.Sum(nil))
	if value == d.value {
		return nil
	}

	// 加载配置
	if err := d.parse(buf); err != nil {
		return err
	}
	d.value = value

	// 配置变更通知
	if !isload {
		for _, f := range d.changes {
			f()
		}
	}
	return nil
}
