package app_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/topokrat/topoclean/internal/app"
	"github.com/topokrat/topoclean/internal/config"
	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

func TestDryRun(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "topoclean_app_test")
	defer os.RemoveAll(tempDir)
	
	jpegMagic := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	os.WriteFile(filepath.Join(tempDir, "image.jpg"), jpegMagic, 0644)

	cfg, _ := config.Load("")
	cfg.HeptagonRoot = tempDir // Simuliere Root im Temp-Dir

	l, _ := ledger.New(filepath.Join(tempDir, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	
	core := app.New(l, s, v, cfg)
	report, err := core.Plan()
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
	tempHome, _ := os.MkdirTemp("", "topoclean_home_execute")
	defer os.RemoveAll(tempHome)
	
	fileName := "document.pdf"
	sourcePath := filepath.Join(tempHome, fileName)
	os.WriteFile(sourcePath, []byte("%PDF-1.4\nsoterisches dokument"), 0644)

	cfg, _ := config.Load("")
	cfg.HeptagonRoot = tempHome

	l, _ := ledger.New(filepath.Join(tempHome, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	err := core.Execute()
	if err != nil {
		t.Fatalf("Execute fehlgeschlagen: %v", err)
	}

	targetPath := filepath.Join(tempHome, "02-Identity", fileName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("Datei wurde nicht nach %s verschoben", targetPath)
	}

	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		t.Error("Quelldatei existiert noch am Ursprungsort")
	}

	txs, _ := l.GetRecentTransactions(1)
	if len(txs) == 0 || len(txs[0].Ops) == 0 {
		t.Error("Keine Operation im Ledger protokolliert")
	}
}

func TestRollback(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "topoclean_rollback")
	defer os.RemoveAll(tempHome)
	
	fileName := "important.go"
	sourcePath := filepath.Join(tempHome, fileName)
	os.WriteFile(sourcePath, []byte("package main\nfunc main() {}"), 0644)

	cfg, _ := config.Load("")
	cfg.HeptagonRoot = tempHome

	l, _ := ledger.New(filepath.Join(tempHome, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	core.Execute()
	
	txs, _ := l.GetRecentTransactions(1)
	txUUID := txs[0].UUID

	err := core.Rollback(txUUID)
	if err != nil {
		t.Fatalf("Rollback fehlgeschlagen: %v", err)
	}

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		t.Error("Datei wurde durch Rollback nicht wiederhergestellt")
	}

	targetPath := filepath.Join(tempHome, "03-Creation", fileName)
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("Zieldatei existiert nach Rollback immer noch")
	}
}

func TestDeduplication(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "topoclean_dedup")
	defer os.RemoveAll(tempHome)
	
	content := []byte("identischer soterischer inhalt")
	os.WriteFile(filepath.Join(tempHome, "file1.txt"), content, 0644)
	os.WriteFile(filepath.Join(tempHome, "file2.txt"), content, 0644)

	cfg, _ := config.Load("")
	cfg.HeptagonRoot = tempHome

	l, _ := ledger.New(filepath.Join(tempHome, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	err := core.Execute()
	if err != nil {
		t.Fatalf("Execute fehlgeschlagen: %v", err)
	}

	f1, _ := os.Stat(filepath.Join(tempHome, "07-Inbox", "file1.txt"))
	f2, _ := os.Stat(filepath.Join(tempHome, "07-Inbox", "file2.txt"))

	if !os.SameFile(f1, f2) {
		t.Error("Deduplizierung fehlgeschlagen: Dateien zeigen nicht auf dieselbe Inode (kein Hardlink)")
	}
}

func TestMultiZoneMapping(t *testing.T) {
	tempRoot, _ := os.MkdirTemp("", "topoclean_multizone")
	defer os.RemoveAll(tempRoot)
	
	// Erzeuge Zonen
	downloadsDir := filepath.Join(tempRoot, "Downloads")
	desktopDir := filepath.Join(tempRoot, "Desktop")
	os.MkdirAll(downloadsDir, 0755)
	os.MkdirAll(desktopDir, 0755)
	
	os.WriteFile(filepath.Join(downloadsDir, "video.mp4"), []byte("video data"), 0644)
	os.WriteFile(filepath.Join(desktopDir, "script.py"), []byte("python code"), 0644)

	cfg := &config.Config{
		HeptagonRoot: tempRoot,
		Zones: []config.Zone{
			{Path: downloadsDir, Name: "Downloads"},
			{Path: desktopDir, Name: "Desktop"},
		},
		Mapping: config.Mapping{
			PreserveOrigin: true,
			DateFormat:     "2006-01",
		},
	}

	l, _ := ledger.New(filepath.Join(tempRoot, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	err := core.Execute()
	if err != nil {
		t.Fatalf("Execute fehlgeschlagen: %v", err)
	}

	// Verifizierung
	expectedVideo := filepath.Join(tempRoot, "05-Media", "From-Downloads", time.Now().Format("2006-01"), "video.mp4")
	if _, err := os.Stat(expectedVideo); os.IsNotExist(err) {
		t.Errorf("Video nicht an erwartetem Ort: %s", expectedVideo)
	}

	expectedScript := filepath.Join(tempRoot, "03-Creation", "From-Desktop", time.Now().Format("2006-01"), "script.py")
	if _, err := os.Stat(expectedScript); os.IsNotExist(err) {
		t.Errorf("Script nicht an erwartetem Ort: %s", expectedScript)
	}
}

func TestForget(t *testing.T) {
	tempRoot, _ := os.MkdirTemp("", "topoclean_forget")
	defer os.RemoveAll(tempRoot)
	
	fileName := "to_be_forgotten.txt"
	filePath := filepath.Join(tempRoot, fileName)
	os.WriteFile(filePath, []byte("vergänglicher inhalt"), 0644)

	cfg := &config.Config{HeptagonRoot: tempRoot}
	l, _ := ledger.New(filepath.Join(tempRoot, "ledger.db"))
	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	// Aktion: Vergessen
	err := core.Forget(filePath)
	if err != nil {
		t.Fatalf("Forget fehlgeschlagen: %v", err)
	}

	// Verifizierung 1: Datei sollte nun im Trash sein
	trashPath := filepath.Join(tempRoot, "04-Archive", ".trash", fileName)
	if _, err := os.Stat(trashPath); os.IsNotExist(err) {
		t.Errorf("Datei nicht im Trash gefunden: %s", trashPath)
	}

	// Verifizierung 2: Original sollte weg sein
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Originaldatei existiert noch")
	}

	// Verifizierung 3: Ledger-Status prüfen
	txs, _ := l.GetRecentTransactions(1)
	if len(txs) == 0 || txs[0].State != "PurgeStaged" {
		t.Errorf("Falscher Transaktions-Status im Ledger: %s", txs[0].State)
	}
}
