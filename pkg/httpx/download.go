package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Download 下载远程文件到本地，支持断点续传
func Download(ctx context.Context, url, filePath string) error {
	// 创建目录（自动支持多级目录）
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	// 检查远程是否支持 Range 请求
	supportRange, remoteSize, err := checkRemoteSupportRange(ctx, url)
	if err != nil {
		return fmt.Errorf("checking remote file failed: %w", err)
	}
	// 获取本地已下载文件大小
	var localSize int64
	if fi, err := os.Stat(filePath); err == nil {
		localSize = fi.Size()
	}
	// 如果远程不支持断点，或者本地文件不存在/为 0，就整文件下载
	if !supportRange || localSize == 0 {
		return fullDownload(ctx, url, filePath)
	}
	// 如果本地已经完整下载了，直接跳过
	if localSize >= remoteSize {
		return nil
	}
	// 否则进行断点下载
	return resumeDownload(ctx, url, filePath, localSize)
}

// checkRemoteSupportRange 检查服务器是否支持断点续传
func checkRemoteSupportRange(ctx context.Context, url string) (bool, int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// 判断是否支持 Range
	supportRange := resp.Header.Get("Accept-Ranges") == "bytes"
	// 获取远程文件大小
	contentLengthStr := resp.Header.Get("Content-Length")
	contentLength, _ := strconv.ParseInt(contentLengthStr, 10, 64)
	return supportRange, contentLength, nil
}

// fullDownload 整文件下载
func fullDownload(ctx context.Context, url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed, status code: %d", resp.StatusCode)
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, resp.Body)
	return err
}

// resumeDownload 断点续传下载
func resumeDownload(ctx context.Context, url, filePath string, start int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", start))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("the server does not support breakpoint resuming or the request failed. Status code: %d", resp.StatusCode)
	}
	// 追加写入
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, resp.Body)
	return err
}
