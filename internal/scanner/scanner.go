package scanner

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Path      string
	Size      int64
	Extension string
	MIMEType  string
	ZoneName  string // Neu: Biographische Herkunft
}

type Scanner struct {
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(dir string, zoneName string) ([]FileInfo, error) {
	var found []FileInfo
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue 
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".db") {
			continue
		}

		path := filepath.Join(dir, name)
		info, err := entry.Info()
		if err != nil {
			continue
		}

		mime, _ := detectMIME(path)

		found = append(found, FileInfo{
			Path:      path,
			Size:      info.Size(),
			Extension: filepath.Ext(name),
			MIMEType:  mime,
			ZoneName:  zoneName,
		})
	}

	return found, nil
}

func detectMIME(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "application/octet-stream", err
	}
	defer f.Close()

	// Lies die ersten 512 Bytes für die topologische Signatur
	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && n == 0 {
		return "application/octet-stream", err
	}

	return http.DetectContentType(buffer[:n]), nil
}
