package downloader

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fetchimg/pkg/utils"
)

const (
	maxFileSize = 10 * 1024 * 1024
	timeout     = 30 * time.Second
)

type Downloader struct {
	outputDir string
	client    *http.Client
}

func NewDownloader(outputDir string) *Downloader {
	return &Downloader{
		outputDir: outputDir,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (d *Downloader) DownloadImage(name, imageURL string) (string, error) {
	_, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ImageDownloader/1.0)")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if err := d.validateResponse(resp); err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	format, err := d.detectImageFormat(data)
	if err != nil {
		return "", err
	}

	sanitizedName := utils.SanitizeFilename(name)
	filename := fmt.Sprintf("%s.%s", sanitizedName, format)
	filepath := filepath.Join(d.outputDir, filename)

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}

	return filepath, nil
}

func (d *Downloader) validateResponse(resp *http.Response) error {
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("not an image file. Content-Type: %s", contentType)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	if resp.ContentLength > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", resp.ContentLength, maxFileSize)
	}

	return nil
}

func (d *Downloader) detectImageFormat(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return "", fmt.Errorf("invalid image format: %w", err)
	}
	return format, nil
}
