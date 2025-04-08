package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	if listenAddr == "" {
		listenAddr = ":8080" 
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	if err != nil || maxIdleConns <= 0 {
		maxIdleConns = 1000 
	}

	maxIdleConnsPerHost, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS_PER_HOST"))
	if err != nil || maxIdleConnsPerHost <= 0 {
		maxIdleConnsPerHost = 1000 // default value
	}

	backendURL, err := url.Parse(backendURLStr)
	if err != nil {
		log.Fatal("Invalid BACKEND_URL: ", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request: %s %s", r.Method, r.URL.String())
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Starting proxy server at %s forwarding to %s", listenAddr, backendURL.String())
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
