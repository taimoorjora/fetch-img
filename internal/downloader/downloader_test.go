package downloader

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func createTestImage() []byte {
	// Create a simple 1x1 JPEG image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red pixel

	// Encode as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestDownloader_DownloadImage(t *testing.T) {
	// Create a test server that serves a valid JPEG image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(createTestImage())
	}))
	defer server.Close()

	// Create a test server that returns an error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer errorServer.Close()

	tmpDir := t.TempDir()
	downloader := NewDownloader(tmpDir)

	tests := []struct {
		name     string
		imageURL string
		wantErr  bool
	}{
		{
			name:     "valid image URL",
			imageURL: server.URL,
			wantErr:  false,
		},
		{
			name:     "invalid URL",
			imageURL: "invalid-url",
			wantErr:  true,
		},
		{
			name:     "non-existent image",
			imageURL: errorServer.URL,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := downloader.DownloadImage("test", tt.imageURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Downloaded file does not exist at %s", path)
				}

				// Check if file has correct extension
				ext := filepath.Ext(path)
				if ext != ".jpg" && ext != ".jpeg" {
					t.Errorf("Expected .jpg or .jpeg extension, got %s", ext)
				}
			}
		})
	}
}

func TestDownloader_validateResponse(t *testing.T) {
	downloader := NewDownloader("test")

	tests := []struct {
		name          string
		contentType   string
		statusCode    int
		contentLength int64
		wantErr       bool
	}{
		{
			name:          "valid response",
			contentType:   "image/jpeg",
			statusCode:    http.StatusOK,
			contentLength: 1024,
			wantErr:       false,
		},
		{
			name:          "invalid content type",
			contentType:   "text/html",
			statusCode:    http.StatusOK,
			contentLength: 1024,
			wantErr:       true,
		},
		{
			name:          "error status code",
			contentType:   "image/jpeg",
			statusCode:    http.StatusNotFound,
			contentLength: 1024,
			wantErr:       true,
		},
		{
			name:          "file too large",
			contentType:   "image/jpeg",
			statusCode:    http.StatusOK,
			contentLength: maxFileSize + 1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Header:        make(http.Header),
				StatusCode:    tt.statusCode,
				ContentLength: tt.contentLength,
			}
			resp.Header.Set("Content-Type", tt.contentType)

			err := downloader.validateResponse(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
