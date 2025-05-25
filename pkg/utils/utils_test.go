package utils

import (
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid filename",
			input:    "test-image.jpg",
			expected: "test-image.jpg",
		},
		{
			name:     "filename with invalid characters",
			input:    "test/image:with*invalid?chars",
			expected: "test_image_with_invalid_chars",
		},
		{
			name:     "filename with spaces",
			input:    "test image with spaces.jpg",
			expected: "test image with spaces.jpg",
		},
		{
			name:     "filename too long",
			input:    strings.Repeat("a", 300),
			expected: strings.Repeat("a", 255),
		},
		{
			name:     "filename with control characters",
			input:    "test\x00\x1Fimage.jpg",
			expected: "test__image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilename(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}
