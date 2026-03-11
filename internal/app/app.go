package app

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/topokrat/topoclean/internal/config"
	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

type ProposedMove struct {
	SourcePath   string
	TargetSphere string
	TargetDir    string 
	MIMEType     string
}

type App struct {
	ledger  *ledger.Ledger
	scanner *scanner.Scanner
	vector  *vector.Vector
	config  *config.Config
}

func New(l *ledger.Ledger, s *scanner.Scanner, v *vector.Vector, cfg *config.Config) *App {
	return &App{
		ledger:  l,
		scanner: s,
		vector:  v,
		config:  cfg,
	}
}

func (a *App) Plan() ([]ProposedMove, error) {
	var allFiles []scanner.FileInfo

	if len(a.config.Zones) > 0 {
		for _, zone := range a.config.Zones {
			files, err := a.scanner.Scan(zone.Path, zone.Name)
			if err != nil {
				fmt.Printf("Warnung: Konnte Zone %s (%s) nicht scannen: %v\n", zone.Name, zone.Path, err)
				continue
			}
			allFiles = append(allFiles, files...)
		}
	} else {
		files, err := a.scanner.Scan(a.config.HeptagonRoot, "")
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, files...)
	}

	var plan []ProposedMove
	for _, f := range allFiles {
		sphere := a.vector.Classify(f)
		targetDir := filepath.Join(a.config.HeptagonRoot, sphere)
		if a.config.Mapping.PreserveOrigin && f.ZoneName != "" {
			dateFolder := time.Now().Format(a.config.Mapping.DateFormat)
			targetDir = filepath.Join(targetDir, "From-"+f.ZoneName, dateFolder)
		}

		plan = append(plan, ProposedMove{
			SourcePath:   f.Path,
			TargetSphere: sphere,
			TargetDir:    targetDir,
			MIMEType:     f.MIMEType,
		})
	}

	return plan, nil
}

func (a *App) Execute() error {
	plan, err := a.Plan()
	if err != nil {
		return err
	}

	if len(plan) == 0 {
		return nil
	}

	tx, err := a.ledger.Begin()
	if err != nil {
		return fmt.Errorf("konnte Transaktion nicht starten: %v", err)
	}
	defer a.ledger.Save(tx)

	processedHashes := make(map[string]string)

	for _, move := range plan {
		err := a.executeMove(tx.UUID, move, processedHashes)
		if err != nil {
			fmt.Printf("Fehler bei %s: %v\n", move.SourcePath, err)
		}
	}

	tx.State = "Committed"
	return nil
}

// Forget verschiebt eine Datei soterisch in den Limbo (04-Archive/.trash)
func (a *App) Forget(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// 1. Ziel im Archiv bestimmen
	trashDir := filepath.Join(a.config.HeptagonRoot, "04-Archive", ".trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}

	// 2. Transaktion starten
	tx, err := a.ledger.Begin()
	if err != nil {
		return err
	}
	tx.State = "PurgeStaged"
	defer a.ledger.Save(tx)

	// 3. Move ausführen
	move := ProposedMove{
		SourcePath:   absPath,
		TargetSphere: "04-Archive",
		TargetDir:    trashDir,
	}
	
	// Wir nutzen executeMove ohne Session-Cache für Single-File Operationen
	return a.executeMove(tx.UUID, move, make(map[string]string))
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
	currentHash, err := calculateHash(op.DestPath)
	if err != nil {
		return fmt.Errorf("konnte Hash von %s nicht prüfen: %v", op.DestPath, err)
	}

	if currentHash != op.FileHash {
		return fmt.Errorf("Integrität von %s verletzt! Rollback abgebrochen", op.DestPath)
	}

	if _, err := os.Stat(op.SourcePath); os.IsNotExist(err) {
		if err := copyFile(op.DestPath, op.SourcePath); err != nil {
			return fmt.Errorf("konnte Datei nicht zurückrollen: %v", err)
		}
	}

	return os.Remove(op.DestPath)
}

func (a *App) executeMove(txUUID string, move ProposedMove, processedHashes map[string]string) error {
	if err := os.MkdirAll(move.TargetDir, 0755); err != nil {
		return fmt.Errorf("konnte Ziel-Verzeichnis %s nicht erstellen: %v", move.TargetDir, err)
	}

	sourceHash, err := calculateHash(move.SourcePath)
	if err != nil {
		return fmt.Errorf("konnte Quell-Hash nicht berechnen: %v", err)
	}

	targetPath := filepath.Join(move.TargetDir, filepath.Base(move.SourcePath))
	
	var existingPath string
	if path, ok := processedHashes[sourceHash]; ok {
		existingPath = path
	} else {
		if path, err := a.ledger.GetPathByHash(sourceHash); err == nil {
			if _, err := os.Stat(path); err == nil {
				existingPath = path
			}
		}
	}

	if existingPath != "" && existingPath != targetPath {
		if err := os.Link(existingPath, targetPath); err == nil {
			processedHashes[sourceHash] = targetPath
			return a.finalizeMove(txUUID, move.SourcePath, targetPath, sourceHash)
		}
	}

	if err := copyFile(move.SourcePath, targetPath); err != nil {
		return fmt.Errorf("konnte Datei nicht kopieren: %v", err)
	}

	targetHash, err := calculateHash(targetPath)
	if err != nil {
		return fmt.Errorf("konnte Ziel-Hash nicht berechnen: %v", err)
	}

	if sourceHash != targetHash {
		os.Remove(targetPath)
		return fmt.Errorf("Hash-Mismatch! Integrität gefährdet")
	}

	processedHashes[sourceHash] = targetPath
	return a.finalizeMove(txUUID, move.SourcePath, targetPath, sourceHash)
}

func (a *App) finalizeMove(txUUID, sourcePath, targetPath, hash string) error {
	info, _ := os.Stat(targetPath)
	op := ledger.Operation{
		SourcePath: sourcePath,
		DestPath:   targetPath,
		FileHash:   hash,
		FileSize:   info.Size(),
	}
	if err := a.ledger.AddOperation(txUUID, op); err != nil {
		return err
	}
	return os.Remove(sourcePath)
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
