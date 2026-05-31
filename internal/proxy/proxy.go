package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

// Config holds the proxy server configuration.
type Config struct {
	ListenAddr          string
	BackendURL          string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		ListenAddr:          ":8080",
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
	}
}

// NewProxy creates an http.Handler that reverse proxies requests to the given backend URL.
func NewProxy(cfg Config) (http.Handler, error) {
	backendURL, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return proxy, nil
}

// ParseConfig parses configuration from environment variable strings.
func ParseConfig(listenAddr, backendURL, maxIdleConnsStr, maxIdleConnsPerHostStr string) Config {
	cfg := DefaultConfig()

	if listenAddr != "" {
		cfg.ListenAddr = listenAddr
	}

	if backendURL != "" {
		cfg.BackendURL = backendURL
	}

	if maxIdleConnsStr != "" {
		if v, err := strconv.Atoi(maxIdleConnsStr); err == nil && v > 0 {
			cfg.MaxIdleConns = v
		}
	}

	if maxIdleConnsPerHostStr != "" {
		if v, err := strconv.Atoi(maxIdleConnsPerHostStr); err == nil && v > 0 {
			cfg.MaxIdleConnsPerHost = v
		}
	}

	return cfg
}
