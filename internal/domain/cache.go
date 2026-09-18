package domain

type CacheResponse struct {
	StatusCode int                 `json:"status_code"`
	Header     map[string][]string `json:"header"`
	Body       []byte              `json:"body"`
}

type CacheRepository interface {
	Get(key string) (*CacheResponse, bool)
	Set(key string, response *CacheResponse)
	Clear() error
	Load() error
}
