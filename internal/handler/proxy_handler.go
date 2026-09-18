package handler

import (
	"net/http"

	"github.com/MrAndrey24/cache-proxy/internal/domain"
	"github.com/MrAndrey24/cache-proxy/internal/service"
)

type proxyHandler struct {
	service *service.ProxyService
}

func NewProxyHandler(service *service.ProxyService) *proxyHandler {
	return &proxyHandler{
		service: service,
	}
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, cacheStatus, err := h.service.ProcessRequest(r)
	if err != nil {
		http.Error(w, "failed to reach origin server"+err.Error(), http.StatusBadGateway)
		return
	}

	h.writeResponse(w, resp, cacheStatus)
}

func (h *proxyHandler) writeResponse(w http.ResponseWriter, resp *domain.CacheResponse, cacheStatus string) {
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	w.Header().Set("X-Cache", cacheStatus)
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}
