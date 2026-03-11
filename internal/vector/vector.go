package vector

import (
	"path/filepath"
	"strings"

	"github.com/topokrat/topoclean/internal/scanner"
)

// Strategy definiert das Interface für Klassifizierungsketten
type Strategy interface {
	Classify(file scanner.FileInfo) (string, bool)
}

// MIMEStrategy klassifiziert nach Magic-Byte Inhalten (Priorität)
type MIMEStrategy struct{}

func (s *MIMEStrategy) Classify(file scanner.FileInfo) (string, bool) {
	mime := strings.ToLower(file.MIMEType)
	if strings.HasPrefix(mime, "video/") || strings.HasPrefix(mime, "audio/") || strings.HasPrefix(mime, "image/") {
		return "05-Media", true
	}
	if strings.HasPrefix(mime, "application/pdf") || strings.HasPrefix(mime, "application/msword") || strings.Contains(mime, "document") {
		return "02-Identity", true
	}
	if strings.HasPrefix(mime, "text/x-") || strings.Contains(mime, "code") {
		return "03-Creation", true
	}
	return "", false
}

// ExtensionStrategy klassifiziert nach bekannten Dateiendungen (Fallback)
type ExtensionStrategy struct{}

func (s *ExtensionStrategy) Classify(file scanner.FileInfo) (string, bool) {
	ext := strings.ToLower(file.Extension)
	
	// Creation (Werkstatt)
	codeExts := map[string]bool{
		".go": true, ".py": true, ".rs": true, ".sh": true, ".js": true, ".ts": true,
		".html": true, ".css": true, ".sql": true, ".c": true, ".cpp": true,
	}
	if codeExts[ext] {
		return "03-Creation", true
	}

	// Media (Konsum)
	mediaExts := map[string]bool{
		".mp4": true, ".mkv": true, ".mov": true, ".avi": true, ".mp3": true, ".flac": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".webm": true,
	}
	if mediaExts[ext] {
		return "05-Media", true
	}

	// Identity (Person)
	docExts := map[string]bool{
		".pdf": true, ".tex": true, ".doc": true, ".docx": true, ".odt": true, ".rtf": true,
	}
	if docExts[ext] {
		return "02-Identity", true
	}

	// Archive (Gedächtnis)
	archExts := map[string]bool{
		".zip": true, ".tar": true, ".gz": true, ".7z": true, ".rar": true, ".apk": true, ".iso": true,
	}
	if archExts[ext] {
		return "04-Archive", true
	}

	return "", false
}

// SubstringStrategy klassifiziert nach semantischen Mustern im Dateinamen (Kontext)
type SubstringStrategy struct{}

func (s *SubstringStrategy) Classify(file scanner.FileInfo) (string, bool) {
	name := strings.ToLower(filepath.Base(file.Path))
	
	// Core / Security
	if strings.Contains(name, "vault") || strings.Contains(name, "key") {
		return "01-Core", true
	}
	
	// Identity / Documents
	if strings.Contains(name, "anschreiben") || strings.Contains(name, "inkasso") || strings.Contains(name, "lebenslauf") {
		return "02-Identity", true
	}
	
	return "", false
}

type Vector struct {
	strategies []Strategy
}

func New() *Vector {
	return &Vector{
		strategies: []Strategy{
			&MIMEStrategy{},
			&ExtensionStrategy{},
			&SubstringStrategy{},
		},
	}
}

func (v *Vector) Classify(file scanner.FileInfo) string {
	for _, strategy := range v.strategies {
		if sphere, ok := strategy.Classify(file); ok {
			return sphere
		}
	}
	return "07-Inbox"
}
