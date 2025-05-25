package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCSV(t *testing.T) {
	// Create a temporary test CSV file
	content := `name,link
image1,https://example.com/image1.jpg
image2,https://example.com/image2.png
image3,https://example.com/image3.gif`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filepath string
		nameCol  string
		linkCol  string
		wantErr  bool
		wantLen  int
	}{
		{
			name:     "valid file with default columns",
			filepath: tmpFile,
			nameCol:  "name",
			linkCol:  "link",
			wantErr:  false,
			wantLen:  3,
		},
		{
			name:     "file does not exist",
			filepath: "nonexistent.csv",
			nameCol:  "name",
			linkCol:  "link",
			wantErr:  true,
			wantLen:  0,
		},
		{
			name:     "invalid column names",
			filepath: tmpFile,
			nameCol:  "invalid",
			linkCol:  "link",
			wantErr:  true,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images, err := ParseCSV(tt.filepath, tt.nameCol, tt.linkCol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(images) != tt.wantLen {
				t.Errorf("ParseCSV() got %d images, want %d", len(images), tt.wantLen)
			}
		})
	}
}

func TestWriteReport(t *testing.T) {
	images := []Image{
		{Name: "test1", Link: "https://example.com/1.jpg", Path: "downloads/test1.jpg", Status: "success"},
		{Name: "test2", Link: "https://example.com/2.jpg", Path: "downloads/test2.jpg", Status: "error: timeout"},
	}

	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report.csv")

	err := WriteReport(images, reportPath)
	if err != nil {
		t.Errorf("WriteReport() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Errorf("Report file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Errorf("Failed to read report file: %v", err)
	}

	expected := "name,link,path,status\n" +
		"test1,https://example.com/1.jpg,downloads/test1.jpg,success\n" +
		"test2,https://example.com/2.jpg,downloads/test2.jpg,error: timeout\n"

	if string(content) != expected {
		t.Errorf("WriteReport() content = %v, want %v", string(content), expected)
	}
}
