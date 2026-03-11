package app

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

type ProposedMove struct {
	SourcePath   string
	TargetSphere string
	MIMEType     string
}

type App struct {
	ledger  *ledger.Ledger
	scanner *scanner.Scanner
	vector  *vector.Vector
}

func New(l *ledger.Ledger, s *scanner.Scanner, v *vector.Vector) *App {
	return &App{
		ledger:  l,
		scanner: s,
		vector:  v,
	}
}

func (a *App) Plan(dir string) ([]ProposedMove, error) {
	files, err := a.scanner.Scan(dir)
	if err != nil {
		return nil, err
	}

	var plan []ProposedMove
	for _, f := range files {
		sphere := a.vector.Classify(f)
		plan = append(plan, ProposedMove{
			SourcePath:   f.Path,
			TargetSphere: sphere,
			MIMEType:     f.MIMEType,
		})
	}

	return plan, nil
}

func (a *App) Execute(dir string) error {
	plan, err := a.Plan(dir)
	if err != nil {
		return err
	}

	if len(plan) == 0 {
		return nil
	}

	// Start Ledger-Transaktion
	tx, err := a.ledger.Begin()
	if err != nil {
		return fmt.Errorf("konnte Transaktion nicht starten: %v", err)
	}
	defer a.ledger.Save(tx) // Speichere Transaktion am Ende

	for _, move := range plan {
		err := a.executeMove(tx.UUID, move, dir)
		if err != nil {
			// Soterischer Fehler-Report: Wir protokollieren den Fehler, brechen aber nicht ab, 
			// um andere Dateien nicht zu blockieren (Graceful Degradation).
			fmt.Printf("Fehler bei %s: %v\n", move.SourcePath, err)
		}
	}

	tx.State = "Committed"
	return nil
}

func (a *App) executeMove(txUUID string, move ProposedMove, baseDir string) error {
	// 1. Ziel-Ordner (Sphäre) sicherstellen
	targetDir := filepath.Join(baseDir, move.TargetSphere)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("konnte Sphäre %s nicht erstellen: %v", move.TargetSphere, err)
	}

	// 2. Pre-Move Hash berechnen (Integritäts-Mandat)
	sourceHash, err := calculateHash(move.SourcePath)
	if err != nil {
		return fmt.Errorf("konnte Quell-Hash nicht berechnen: %v", err)
	}

	targetPath := filepath.Join(targetDir, filepath.Base(move.SourcePath))
	
	// 3. Atomarer Move (oder Copy + Delete für Cross-Filesystem Sicherheit)
	if err := copyFile(move.SourcePath, targetPath); err != nil {
		return fmt.Errorf("konnte Datei nicht kopieren: %v", err)
	}

	// 4. Post-Move Hash verifizieren (Anti-Corruption Mandat)
	targetHash, err := calculateHash(targetPath)
	if err != nil {
		return fmt.Errorf("konnte Ziel-Hash nicht berechnen: %v", err)
	}

	if sourceHash != targetHash {
		os.Remove(targetPath) // Korruptes Ziel löschen
		return fmt.Errorf("Hash-Mismatch! Integrität gefährdet: %s != %s", sourceHash, targetHash)
	}

	// 5. Operation im Ledger protokollieren
	info, _ := os.Stat(targetPath)
	op := ledger.Operation{
		SourcePath: move.SourcePath,
		DestPath:   targetPath,
		FileHash:   sourceHash,
		FileSize:   info.Size(),
	}
	if err := a.ledger.AddOperation(txUUID, op); err != nil {
		return fmt.Errorf("konnte Operation nicht im Ledger speichern: %v", err)
	}

	// 6. Finalisierung: Quelldatei löschen
	return os.Remove(move.SourcePath)
}

func calculateHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
