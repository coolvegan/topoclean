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
			fmt.Printf("Fehler bei %s: %v\n", move.SourcePath, err)
		}
	}

	tx.State = "Committed"
	return nil
}

func (a *App) Rollback(txUUID string) error {
	tx, err := a.ledger.Get(txUUID)
	if err != nil {
		return fmt.Errorf("Transaktion %s nicht gefunden: %v", txUUID, err)
	}

	if tx.State == "RolledBack" {
		return fmt.Errorf("Transaktion %s wurde bereits zurückgerollt", txUUID)
	}

	for _, op := range tx.Ops {
		err := a.executeRollbackOp(op)
		if err != nil {
			fmt.Printf("Rollback-Fehler bei %s: %v\n", op.DestPath, err)
		}
	}

	return a.ledger.UpdateTransactionState(txUUID, "RolledBack")
}

func (a *App) executeRollbackOp(op ledger.Operation) error {
	// 1. Integrität prüfen (Inversions-Validierung)
	currentHash, err := calculateHash(op.DestPath)
	if err != nil {
		return fmt.Errorf("konnte Hash von %s nicht prüfen: %v", op.DestPath, err)
	}

	if currentHash != op.FileHash {
		return fmt.Errorf("Integrität von %s verletzt! Rollback abgebrochen", op.DestPath)
	}

	// 2. Zurückkopieren
	if err := copyFile(op.DestPath, op.SourcePath); err != nil {
		return fmt.Errorf("konnte Datei nicht zurückrollen: %v", err)
	}

	// 3. Ziel löschen
	return os.Remove(op.DestPath)
}

func (a *App) executeMove(txUUID string, move ProposedMove, baseDir string) error {
	targetDir := filepath.Join(baseDir, move.TargetSphere)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("konnte Sphäre %s nicht erstellen: %v", move.TargetSphere, err)
	}

	sourceHash, err := calculateHash(move.SourcePath)
	if err != nil {
		return fmt.Errorf("konnte Quell-Hash nicht berechnen: %v", err)
	}

	targetPath := filepath.Join(targetDir, filepath.Base(move.SourcePath))
	
	if err := copyFile(move.SourcePath, targetPath); err != nil {
		return fmt.Errorf("konnte Datei nicht kopieren: %v", err)
	}

	targetHash, err := calculateHash(targetPath)
	if err != nil {
		return fmt.Errorf("konnte Ziel-Hash nicht berechnen: %v", err)
	}

	if sourceHash != targetHash {
		os.Remove(targetPath)
		return fmt.Errorf("Hash-Mismatch! Integrität gefährdet: %s != %s", sourceHash, targetHash)
	}

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
