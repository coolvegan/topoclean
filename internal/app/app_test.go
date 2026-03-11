package app_test

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/topokrat/topoclean/internal/app"
	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

func TestDryRun(t *testing.T) {
	// Setup
	tempDir, _ := os.MkdirTemp("", "topoclean_app_test")
	defer os.RemoveAll(tempDir)
	
	// Erstelle eine Testdatei (Video)
	jpegMagic := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	os.WriteFile(filepath.Join(tempDir, "image.jpg"), jpegMagic, 0644)

	l, _ := ledger.New(filepath.Join(tempDir, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	
	core := app.New(l, s, v)
	report, err := core.Plan(tempDir)
	if err != nil {
		t.Fatalf("Planung fehlgeschlagen: %v", err)
	}

	if len(report) != 1 {
		t.Fatalf("erwartete 1 Datei im Plan, erhalten: %d", len(report))
	}

	if report[0].TargetSphere != "05-Media" {
		t.Errorf("erwartete Ziel-Sphäre 05-Media, erhalten: %s", report[0].TargetSphere)
	}
}

func TestExecute(t *testing.T) {
	// Setup: Temporäres Home und Testdatei
	tempHome, _ := os.MkdirTemp("", "topoclean_home_execute")
	defer os.RemoveAll(tempHome)
	
	fileName := "document.pdf"
	sourcePath := filepath.Join(tempHome, fileName)
	os.WriteFile(sourcePath, []byte("%PDF-1.4\nsoterisches dokument"), 0644)

	// Initialisiere App
	l, _ := ledger.New(filepath.Join(tempHome, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v)

	// Führe Execute aus
	err := core.Execute(tempHome)
	if err != nil {
		t.Fatalf("Execute fehlgeschlagen: %v", err)
	}

	// Validierung 1: Datei sollte nun in 02-Identity/ liegen (da .pdf)
	targetPath := filepath.Join(tempHome, "02-Identity", fileName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("Datei wurde nicht nach %s verschoben", targetPath)
	}

	// Validierung 2: Quelldatei sollte weg sein
	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		t.Error("Quelldatei existiert noch am Ursprungsort")
	}

	// Validierung 3: Ledger sollte die Operation enthalten
	txs, _ := l.GetRecentTransactions(1)
	if len(txs) == 0 || len(txs[0].Ops) == 0 {
		t.Error("Keine Operation im Ledger protokolliert")
	}
}

func TestRollback(t *testing.T) {
	// Setup
	tempHome, _ := os.MkdirTemp("", "topoclean_rollback")
	defer os.RemoveAll(tempHome)
	
	fileName := "important.go"
	sourcePath := filepath.Join(tempHome, fileName)
	os.WriteFile(sourcePath, []byte("package main\nfunc main() {}"), 0644)

	l, _ := ledger.New(filepath.Join(tempHome, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v)

	// 1. Ausführen
	core.Execute(tempHome)
	
	txs, _ := l.GetRecentTransactions(1)
	txUUID := txs[0].UUID

	// 2. Rollback
	err := core.Rollback(txUUID)
	if err != nil {
		t.Fatalf("Rollback fehlgeschlagen: %v", err)
	}

	// 3. Validierung: Datei muss wieder am Ursprungsort sein
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		t.Error("Datei wurde durch Rollback nicht wiederhergestellt")
	}

	// 4. Validierung: Ziel-Ordner (Sphäre) sollte leer sein (optional)
	targetPath := filepath.Join(tempHome, "03-Creation", fileName)
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("Zieldatei existiert nach Rollback immer noch")
	}
}
