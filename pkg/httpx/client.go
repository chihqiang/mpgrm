package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

var (
	defaultClient = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
		Timeout: 10 * time.Second,
	}
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

func Request(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	var lastErr error
	baseDelay := retryDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}
		for hk, hv := range headers {
			req.Header.Set(hk, hv)
		}
		resp, err := defaultClient.Do(req)
		if err != nil {
			lastErr = err
		} else {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, nil
			}
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			defer func() {
				_ = resp.Body.Close()
			}()
			lastErr = fmt.Errorf("%s %s failed: status %d, body: %q", method, url, resp.StatusCode, bodyBytes)
		}

		if attempt < maxRetries {
			backoff := time.Duration(1<<attempt) * baseDelay
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, lastErr
}
