package vector_test

import (
	"testing"
	"github.com/topokrat/topoclean/internal/vector"
	"github.com/topokrat/topoclean/internal/scanner"
)

func TestVectorStrategyPipeline(t *testing.T) {
	v := vector.New()

	tests := []struct {
		name     string
		file     scanner.FileInfo
		expected string
	}{
		{
			name: "MIME-Strategy (Priority)",
			file: scanner.FileInfo{Path: "video.txt", Extension: ".txt", MIMEType: "video/mp4"},
			expected: "05-Media",
		},
		{
			name: "Extension-Strategy (Fallback)",
			file: scanner.FileInfo{Path: "main.go", Extension: ".go", MIMEType: "text/plain"},
			expected: "03-Creation",
		},
		{
			name: "Substring-Strategy (Context)",
			file: scanner.FileInfo{Path: "my-secret-vault.bin", Extension: ".bin", MIMEType: "application/octet-stream"},
			expected: "01-Core",
		},
		{
			name: "Substring-Strategy (Context 2)",
			file: scanner.FileInfo{Path: "Arbeitsagentur-Inkasso.tex", Extension: ".tex", MIMEType: "text/plain"},
			expected: "02-Identity", // Identitätsrelevantes Dokument
		},
		{
			name: "Default-Strategy (Inbox)",
			file: scanner.FileInfo{Path: "unknown.data", Extension: ".data", MIMEType: "application/octet-stream"},
			expected: "07-Inbox",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sphere := v.Classify(tt.file)
			if sphere != tt.expected {
				t.Errorf("%s: erwartet %s, erhalten %s", tt.name, tt.expected, sphere)
			}
		})
	}
}
