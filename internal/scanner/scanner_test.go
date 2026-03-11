package scanner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/topokrat/topoclean/internal/scanner"
)

func TestScanWithMIME(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "topoclean_test_mime")
	if err != nil {
		t.Fatalf("konnte Temp-Dir nicht erstellen: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// JPEG Magic Bytes: FF D8 FF E0
	jpegMagic := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	filePath := filepath.Join(tempDir, "unknown_image")
	err = os.WriteFile(filePath, jpegMagic, 0644)
	if err != nil {
		t.Fatalf("konnte Testdatei nicht erstellen: %v", err)
	}

	s := scanner.New()
	found, err := s.Scan(tempDir, "TestZone")
	if err != nil {
		t.Fatalf("Scan fehlgeschlagen: %v", err)
	}

	if len(found) != 1 {
		t.Fatalf("erwartete 1 Datei, gefunden: %d", len(found))
	}

	if found[0].ZoneName != "TestZone" {
		t.Errorf("erwarteter ZoneName 'TestZone', erhalten: %s", found[0].ZoneName)
	}

	// Wir erwarten, dass der Scanner den MIME-Type erkannt hat
	if !strings.HasPrefix(found[0].MIMEType, "image/") {
		t.Errorf("erwarteter MIME-Type image/*, erhalten: %s", found[0].MIMEType)
	}
}
