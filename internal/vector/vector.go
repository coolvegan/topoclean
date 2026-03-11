package vector

import (
	"strings"
	"github.com/topokrat/topoclean/internal/scanner"
)

type Vector struct {
}

func New() *Vector {
	return &Vector{}
}

func (v *Vector) Classify(file scanner.FileInfo) string {
	mime := strings.ToLower(file.MIMEType)

	// 1. MIME-Präfix Klassifizierung
	if strings.HasPrefix(mime, "video/") || strings.HasPrefix(mime, "audio/") || strings.HasPrefix(mime, "image/") {
		return "05-Media"
	}

	if strings.HasPrefix(mime, "application/pdf") || strings.HasPrefix(mime, "application/msword") {
		return "02-Identity"
	}

	// 2. Code-spezifische MIME-Types oder Dateiendungen
	if strings.HasPrefix(mime, "text/x-") || strings.Contains(mime, "code") || isCodeExtension(file.Extension) {
		return "03-Creation"
	}

	// 3. Fallback zu Inbox
	return "07-Inbox"
}

func isCodeExtension(ext string) bool {
	ext = strings.ToLower(ext)
	codeExts := map[string]bool{
		".go": true, ".py": true, ".rs": true, ".sh": true, ".js": true, ".ts": true,
	}
	return codeExts[ext]
}
