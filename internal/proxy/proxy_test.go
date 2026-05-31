package proxy_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proxy-server/internal/proxy"
)

// startTestBackend creates an httptest.Server that acts as a configurable backend.
func startTestBackend(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// startProxy creates an httptest.Server that proxies to the given backend URL.
func startProxy(t *testing.T, backendURL string) *httptest.Server {
	t.Helper()

	cfg := proxy.Config{
		BackendURL:          backendURL,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 5,
	}

	handler, err := proxy.NewProxy(cfg)
	require.NoError(t, err)

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

// ---------- Path Preservation ----------

func TestProxy_PathPreservation(t *testing.T) {
	var capturedPath string
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Path))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Get(proxySrv.URL + "/api/v1/resource")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "/api/v1/resource", string(body))
	assert.Equal(t, "/api/v1/resource", capturedPath)
}

// ---------- Query Parameters ----------

func TestProxy_QueryParameters(t *testing.T) {
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "key=value&page=1", r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.RawQuery))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Get(proxySrv.URL + "/search?key=value&page=1")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "key=value&page=1", string(body))
}

// ---------- Header Forwarding ----------

func TestProxy_HeaderForwarding(t *testing.T) {
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	req, err := http.NewRequest(http.MethodGet, proxySrv.URL+"/headers", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "custom-value")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// ---------- Hop-by-hop Headers Stripping ----------

func TestProxy_HopByHopHeadersStripped(t *testing.T) {
	var receivedHeaders http.Header
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	// Hop-by-hop headers that should be stripped by the proxy
	req, err := http.NewRequest(http.MethodGet, proxySrv.URL+"/hop", nil)
	require.NoError(t, err)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "timeout=5")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Trailer", "X-Trailer")
	req.Header.Set("Proxy-Authorization", "basic xyz")
	req.Header.Set("Proxy-Authenticate", "basic")
	req.Header.Set("TE", "trailers")

	// Regular headers that should still be forwarded
	req.Header.Set("X-Custom", "should-be-forwarded")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify hop-by-hop headers are stripped
	assert.Empty(t, receivedHeaders.Get("Connection"), "Connection should be stripped")
	assert.Empty(t, receivedHeaders.Get("Transfer-Encoding"), "Transfer-Encoding should be stripped")
	assert.Empty(t, receivedHeaders.Get("Proxy-Connection"), "Proxy-Connection should be stripped")
	assert.Empty(t, receivedHeaders.Get("Keep-Alive"), "Keep-Alive should be stripped")

	// Verify regular headers are forwarded
	assert.Equal(t, "should-be-forwarded", receivedHeaders.Get("X-Custom"))
}

// ---------- Backend Unavailable ----------

func TestProxy_BackendUnavailable(t *testing.T) {
	cfg := proxy.Config{
		BackendURL:          "http://localhost:1",
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
	}

	handler, err := proxy.NewProxy(cfg)
	require.NoError(t, err)

	proxySrv := httptest.NewServer(handler)
	defer proxySrv.Close()

	resp, err := http.Get(proxySrv.URL + "/test")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}

// ---------- Concurrent Requests ----------

func TestProxy_ConcurrentRequests(t *testing.T) {
	var mu sync.Mutex
	requestCount := 0

	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	const numRequests = 50
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get(proxySrv.URL + "/concurrent")
			if err != nil {
				t.Errorf("request failed: %v", err)
				return
			}
			defer resp.Body.Close()
			io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("unexpected status: %d", resp.StatusCode)
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, numRequests, requestCount)
}

// ---------- Configuration Parsing ----------

func TestParseConfig_Defaults(t *testing.T) {
	cfg := proxy.ParseConfig("", "", "", "")
	assert.Equal(t, ":8080", cfg.ListenAddr)
	assert.Equal(t, "", cfg.BackendURL)
	assert.Equal(t, 1000, cfg.MaxIdleConns)
	assert.Equal(t, 1000, cfg.MaxIdleConnsPerHost)
}

func TestParseConfig_CustomValues(t *testing.T) {
	cfg := proxy.ParseConfig(":9090", "http://example.com", "500", "50")
	assert.Equal(t, ":9090", cfg.ListenAddr)
	assert.Equal(t, "http://example.com", cfg.BackendURL)
	assert.Equal(t, 500, cfg.MaxIdleConns)
	assert.Equal(t, 50, cfg.MaxIdleConnsPerHost)
}

func TestParseConfig_InvalidMaxIdleConns(t *testing.T) {
	cfg := proxy.ParseConfig(":8080", "http://example.com", "invalid", "invalid")
	assert.Equal(t, 1000, cfg.MaxIdleConns)
	assert.Equal(t, 1000, cfg.MaxIdleConnsPerHost)
}

func TestParseConfig_ZeroOrNegativeValues(t *testing.T) {
	cfg := proxy.ParseConfig(":8080", "http://example.com", "0", "-5")
	assert.Equal(t, 1000, cfg.MaxIdleConns)
	assert.Equal(t, 1000, cfg.MaxIdleConnsPerHost)
}

// ---------- Basic Request Forwarding ----------

func TestProxy_GetRequest(t *testing.T) {
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from backend"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Get(proxySrv.URL + "/test")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello from backend", string(body))
}

func TestProxy_PostRequest(t *testing.T) {
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "request body", string(body))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Post(proxySrv.URL+"/data", "text/plain", strings.NewReader("request body"))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "created", string(body))
}

// ---------- Large Payload ----------

func TestProxy_LargePayload(t *testing.T) {
	payloadSize := 10 * 1024 * 1024 // 10 MB
	payload := strings.Repeat("A", payloadSize)

	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, payloadSize, len(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(len(body))))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Post(proxySrv.URL+"/large", "application/octet-stream", strings.NewReader(payload))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(payloadSize), string(body))
}

// ---------- All HTTP Methods ----------

func TestProxy_AllMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, method, r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.Method))
			})
			defer backend.Close()

			proxySrv := startProxy(t, backend.URL)

			req, err := http.NewRequest(method, proxySrv.URL+"/method", bytes.NewReader([]byte("body")))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, method, string(body))
		})
	}
}

// ---------- Response Headers ----------

func TestProxy_ResponseHeaders(t *testing.T) {
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-Header", "from-backend")
		w.Header().Set("X-Custom-ID", "12345")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Get(proxySrv.URL + "/headers")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "from-backend", resp.Header.Get("X-Backend-Header"))
	assert.Equal(t, "12345", resp.Header.Get("X-Custom-ID"))
}

// ---------- Content-Length ----------

func TestProxy_ContentLength(t *testing.T) {
	bodyContent := "hello world"
	backend := startTestBackend(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, bodyContent, string(body))
		w.Header().Set("Content-Length", strconv.Itoa(len(bodyContent)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(bodyContent))
	})
	defer backend.Close()

	proxySrv := startProxy(t, backend.URL)

	resp, err := http.Post(proxySrv.URL+"/data", "text/plain", strings.NewReader(bodyContent))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, bodyContent, string(body))
	assert.Equal(t, strconv.Itoa(len(bodyContent)), resp.Header.Get("Content-Length"))
}
