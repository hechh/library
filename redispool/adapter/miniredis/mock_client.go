package miniredis

import (
	"strconv"

	"github.com/hechh/library/redispool"
	"github.com/hechh/library/redispool/adapter/goredis"

	"github.com/alicebob/miniredis/v2"
)

type Client struct {
	*goredis.Client
	miniredis *miniredis.Miniredis
}

func (m *Client) Init(cfg *redispool.DbConfig) error {
	m.Client = &goredis.Client{}
	s, err := miniredis.Run()
	if err != nil {
		return err
	}
	m.miniredis = s

	port, _ := strconv.Atoi(s.Port())
	return m.Client.Init(&redispool.DbConfig{
		Ip:     s.Host(),
		Port:   uint32(port),
		Db:     0,
		Prefix: "mock",
	})
}

func (m *Client) Close() error {
	if m.Client != nil {
		_ = m.Client.Close()
	}
	if m.miniredis != nil {
		m.miniredis.Close()
	}
	return nil
}
