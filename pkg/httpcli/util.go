package httpcli

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	Timeout             int64 `yaml:"timeout,omitempty"`                 // 总超时时间（秒，推荐 5-10s）
	MaxIdleConns        int   `yaml:"max_idle_conns,omitempty"`          // 最大空闲连接数
	MaxIdleConnsPerHost int   `yaml:"max_idle_conns_per_host,omitempty"` // 每个Host最大空闲连接
	IdleConnTimeout     int64 `yaml:"idle_conn_timeout,omitempty"`       // 空闲连接超时（秒）
}

func defaultConfig() *Config {
	return &Config{
		Timeout:             5, // * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30, // * time.Second,
	}
}

type RequestOption func(r *http.Request)

func WithHeader(key, value string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}

func WithContext(ctx context.Context) RequestOption {
	return func(r *http.Request) {
		*r = *r.WithContext(ctx)
	}
}

// getRealIP 从 HTTP 请求中提取客户端真实 IP 地址
// 优先级：X-Forwarded-For > X-Real-IP > RemoteAddr
func GetRealIP(r *http.Request) string {
	// 1. 检查 X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ","); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	// 2. 检查 X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	// 3. 回退到 RemoteAddr
	if host := r.RemoteAddr; host != "" {
		if i := strings.LastIndex(host, ":"); i > 0 {
			return host[:i]
		}
		return host
	}
	return ""
}
