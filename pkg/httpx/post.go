package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// PostD 发送 JSON POST 请求，将响应 JSON 解析到 d，并返回 http.Response
func PostD(ctx context.Context, url string, body io.Reader, d any, headers map[string]string) (*http.Response, error) {
	// 调用 Post 请求（支持重试）
	resp, err := Post(ctx, url, body, headers)

	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	// 解析 JSON 响应
	err = json.NewDecoder(resp.Body).Decode(d)
	return resp, err
}

// Post 发送 JSON POST 请求，支持重试
func Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return Request(ctx, http.MethodPost, url, body, headers)
}

func Upload(ctx context.Context, url string, filePath, fieldName string, fields map[string]string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create form file failed: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copy file failed: %w", err)
	}
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer failed: %w", err)
	}
	return Request(ctx, http.MethodPost, url, body, map[string]string{"Content-Type": writer.FormDataContentType()})
}
