package vector_test

import (
	"testing"
	"github.com/topokrat/topoclean/internal/vector"
	"github.com/topokrat/topoclean/internal/scanner"
)

func TestClassify(t *testing.T) {
	v := vector.New()

	tests := []struct {
		name     string
		file     scanner.FileInfo
		expected string
	}{
		{
			name: "MIME video/* zu Media",
			file: scanner.FileInfo{Path: "secret_video.txt", Extension: ".txt", MIMEType: "video/mp4"},
			expected: "05-Media",
		},
		{
			name: "MIME image/* zu Media",
			file: scanner.FileInfo{Path: "photo.jpg", Extension: ".jpg", MIMEType: "image/jpeg"},
			expected: "05-Media",
		},
		{
			name: "MIME application/pdf zu Identity",
			file: scanner.FileInfo{Path: "vertrag.pdf", Extension: ".pdf", MIMEType: "application/pdf"},
			expected: "02-Identity",
		},
		{
			name: "MIME text/x-python zu Creation",
			file: scanner.FileInfo{Path: "script.txt", Extension: ".txt", MIMEType: "text/x-python"},
			expected: "03-Creation",
		},
		{
			name: "Unbekannt zu Inbox",
			file: scanner.FileInfo{Path: "unbekannt.dat", Extension: ".dat", MIMEType: "application/octet-stream"},
			expected: "07-Inbox",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sphere := v.Classify(tt.file)
			if sphere != tt.expected {
				t.Errorf("%s: erwartete Sphäre %s, erhalten %s", tt.name, tt.expected, sphere)
			}
		})
	}
}
