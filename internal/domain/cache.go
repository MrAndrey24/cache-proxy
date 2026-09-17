package domain

import "net/http"

type CacheResponse struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

type CacheRepository interface {
	Get(key string) (*CacheResponse, bool)
	Set(key string, response *CacheResponse) error
	Clear() error
}
