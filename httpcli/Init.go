package httpcli

import (
	"context"
	"fmt"
	"net/http"
)

var object *HttpClient

func SetObject(oj *HttpClient) {
	object = oj
}

// Get 便捷 GET 请求（全局单例）
func Get[T any](ctx context.Context, url string, opts ...RequestOption) (*T, error) {
	if object != nil {
		body, err := object.Do(ctx, http.MethodGet, url, nil, opts...)
		if err != nil {
			return nil, err
		}
		return unmarshalResponse[T](body)
	}
	return nil, fmt.Errorf("http client not initialized")
}

// Post 便捷 POST 请求（全局单例）
func Post[T any](ctx context.Context, url string, reqBody interface{}, opts ...RequestOption) (*T, error) {
	if object != nil {
		body, err := object.Do(ctx, http.MethodPost, url, reqBody, opts...)
		if err != nil {
			return nil, err
		}
		return unmarshalResponse[T](body)
	}
	return nil, fmt.Errorf("http client not initialized")
}

// Put 便捷 PUT 请求（全局单例）
func Put[T any](ctx context.Context, url string, reqBody interface{}, opts ...RequestOption) (*T, error) {
	if object != nil {
		body, err := object.Do(ctx, http.MethodPut, url, reqBody, opts...)
		if err != nil {
			return nil, err
		}
		return unmarshalResponse[T](body)
	}
	return nil, fmt.Errorf("http client not initialized")
}

// Delete 便捷 DELETE 请求（全局单例）
func Delete[T any](ctx context.Context, url string, opts ...RequestOption) (*T, error) {
	if object != nil {
		body, err := object.Do(ctx, http.MethodDelete, url, nil, opts...)
		if err != nil {
			return nil, err
		}
		return unmarshalResponse[T](body)
	}
	return nil, fmt.Errorf("http client not initialized")
}

// Do 便捷通用请求（全局单例）
func Do(ctx context.Context, method, url string, reqBody interface{}, opts ...RequestOption) ([]byte, error) {
	if object != nil {
		return object.Do(ctx, method, url, reqBody, opts...)
	}
	return nil, fmt.Errorf("http client not initialized")
}
