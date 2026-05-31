// this is a simple HTTP reverse proxy server in Go that forwards requests to a backend server.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"proxy-server/internal/proxy"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	backendURLStr := os.Getenv("BACKEND_URL")
	if backendURLStr == "" {
		log.Fatal("BACKEND_URL not set in .env")
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	maxIdleConnsStr := os.Getenv("MAX_IDLE_CONNS")
	maxIdleConnsPerHostStr := os.Getenv("MAX_IDLE_CONNS_PER_HOST")

	cfg := proxy.ParseConfig(listenAddr, backendURLStr, maxIdleConnsStr, maxIdleConnsPerHostStr)

	handler, err := proxy.NewProxy(cfg)
	if err != nil {
		log.Fatal("Failed to create proxy: ", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request: %s %s", r.Method, r.URL.String())
		handler.ServeHTTP(w, r)
	})

	log.Printf("Starting proxy server at %s forwarding to %s", cfg.ListenAddr, cfg.BackendURL)
	if err := http.ListenAndServe(cfg.ListenAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
