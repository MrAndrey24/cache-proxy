package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/MrAndrey24/cache-proxy/internal/data"
	"github.com/MrAndrey24/cache-proxy/internal/handler"
	"github.com/MrAndrey24/cache-proxy/internal/service"
)

type cacheRepositoryAdapter struct {
	*data.FileCacheRepository
}

func (r cacheRepositoryAdapter) Load() error {
	r.FileCacheRepository.Load()
	return nil
}

func Execute() {
	program()
}

func program() {
	port := flag.Int("port", 0, "Port on which the caching proxy server will run")
	origin := flag.String("origin", "", "URL of the server to which request will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "Clear the cached responses and exit")
	flag.Parse()

	repo := data.NewFileCacheRepository()

	if *clearCache {
		if err := repo.Clear(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *port == 0 || *origin == "" {
		fmt.Println("Usage:")
		fmt.Println(" caching-proxy --port <number> --origin <url>")
		fmt.Println(" caching-proxy --clear-cache")
		os.Exit(1)
	}

	originURL, err := url.Parse(*origin)
	if err != nil || originURL.Scheme == "" || originURL.Host == "" {
		log.Fatal("Invalid --origin URL: %s", *origin)
	}

	cacheRepo := cacheRepositoryAdapter{FileCacheRepository: repo}
	if err := cacheRepo.Load(); err != nil {
		log.Fatal(err)
	}

	proxySvs := service.NewProxyService(cacheRepo, originURL)
	proxyHandler := handler.NewProxyHandler(proxySvs)

	http.HandleFunc("/", proxyHandler.ServeHTTP)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Caching proxy server listening on %s, forwarding to %s", addr, originURL.String())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
