package service

import (
	"io"
	"net/http"
	"net/url"

	"github.com/MrAndrey24/cache-proxy/internal/domain"
)

type ProxyService struct {
	repo   domain.CacheRepository
	origin *url.URL
}

func NewProxyService(repo domain.CacheRepository, origin *url.URL) *ProxyService {
	return &ProxyService{
		repo:   repo,
		origin: origin,
	}
}

func (s *ProxyService) ProcessRequest(r *http.Request) (*domain.CacheResponse, string, error) {
	key := r.Method + " " + r.URL.RequestURI()

	if resp, hit := s.repo.Get(key); hit {
		return resp, "HIT", nil
	}

	resp, err := s.forwardRequest(r)
	if err != nil {
		return &domain.CacheResponse{}, "", err
	}

	if r.Method == http.MethodGet && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.repo.Set(key, resp)
	}

	return resp, "MISS", nil
}

func (s *ProxyService) forwardRequest(r *http.Request) (*domain.CacheResponse, error) {
	target := *s.origin
	target.Path = s.origin.Path + r.URL.Path
	target.RawQuery = r.URL.RawQuery

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &domain.CacheResponse{}, err
	}

	req, err := http.NewRequest(r.Method, target.String(), newBodyReader(body))
	if err != nil {
		return &domain.CacheResponse{}, err
	}
	copyHeader(req.Header, r.Header)
	req.Host = s.origin.Host

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &domain.CacheResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &domain.CacheResponse{}, err
	}

	return &domain.CacheResponse{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}, nil
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, v := range values {
			dst.Add(key, v)
		}
	}
}

func newBodyReader(body []byte) io.Reader {
	if len(body) == 0 {
		return nil
	}
	return &bytesReader{data: body}
}

type bytesReader struct {
	data []byte
	pos  int
}

func (b *bytesReader) Read(p []byte) (int, error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}
