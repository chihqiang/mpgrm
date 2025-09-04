package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetD(ctx context.Context, url string, d any) (*http.Response, error) {
	resp, err := Get(ctx, url)
	if err != nil {
		return nil, err
	}
	// 如果不是成功状态码，提前返回错误
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		defer func() {
			_ = resp.Body.Close()
		}()
		return resp, fmt.Errorf("GET %s failed: status %d, body: %q", url, resp.StatusCode, bodyBytes)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// 把 resp.Body 内容解码到 d
	err = json.NewDecoder(resp.Body).Decode(d)
	return resp, err
}

func Get(ctx context.Context, url string) (*http.Response, error) {
	return Request(ctx, http.MethodGet, url, nil, map[string]string{})
}
