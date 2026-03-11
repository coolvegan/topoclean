package ledger_test

import (
	"os"
	"testing"
	"github.com/topokrat/topoclean/internal/ledger"
)

func TestNewTransaction(t *testing.T) {
	l, err := ledger.New("test_ledger.db")
	if err != nil {
		t.Fatalf("konnte Ledger nicht initialisieren: %v", err)
	}
	tx, err := l.Begin()
	if err != nil {
		t.Fatalf("konnte Transaktion nicht starten: %v", err)
	}
	if tx.UUID == "" {
		t.Error("Transaktion sollte eine UUID haben")
	}
	if tx.State != "Pending" {
		t.Errorf("erwarteter Status 'Pending', erhalten '%s'", tx.State)
	}
}

func TestOperations(t *testing.T) {
	dbPath := "test_operations.db"
	l, err := ledger.New(dbPath)
	if err != nil {
		t.Fatalf("konnte Ledger nicht initialisieren: %v", err)
	}
	
	tx, _ := l.Begin()
	op := ledger.Operation{
		SourcePath: "/home/marco/test.mp4",
		DestPath:   "/home/marco/Videos/test.mp4",
		FileHash:   "sha256:12345",
		FileSize:   1024,
	}
	
	err = l.AddOperation(tx.UUID, op)
	if err != nil {
		t.Fatalf("konnte Operation nicht hinzufügen: %v", err)
	}
	
	err = l.Save(tx)
	if err != nil {
		t.Fatalf("konnte Transaktion nicht speichern: %v", err)
	}

	loadedTx, err := l.Get(tx.UUID)
	if err != nil {
		t.Fatalf("konnte Transaktion nicht laden: %v", err)
	}
	
	if len(loadedTx.Ops) != 1 {
		t.Errorf("erwartete 1 Operation, erhalten %d", len(loadedTx.Ops))
	}
	
	if loadedTx.Ops[0].SourcePath != op.SourcePath {
		t.Errorf("SourcePath mismatch: %s != %s", loadedTx.Ops[0].SourcePath, op.SourcePath)
	}
}

func TestLocate(t *testing.T) {
	dbPath := "test_locate.db"
	l, _ := ledger.New(dbPath)
	defer os.Remove(dbPath)

	tx, _ := l.Begin()
	op := ledger.Operation{
		SourcePath: "/home/marco/my_document.pdf",
		DestPath:   "/home/marco/02-Identity/my_document.pdf",
		FileHash:   "abc",
	}
	l.AddOperation(tx.UUID, op)
	l.Save(tx)

	// Suche nach 'document'
	results, err := l.Locate("document")
	if err != nil {
		t.Fatalf("Locate fehlgeschlagen: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("erwartete 1 Ergebnis, erhalten %d", len(results))
	}

	if results[0].DestPath != op.DestPath {
		t.Errorf("falscher Pfad gefunden: %s", results[0].DestPath)
	}
}

