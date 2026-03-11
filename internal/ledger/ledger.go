package ledger

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Transaction struct {
	UUID      string
	Timestamp time.Time
	State     string
	Ops       []Operation
}

type Operation struct {
	SourcePath string
	DestPath   string
	FileHash   string
	FileSize   int64
}

type Ledger struct {
	db *sql.DB
}

func New(dbPath string) (*Ledger, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		uuid TEXT PRIMARY KEY,
		timestamp DATETIME,
		state TEXT
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS operations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tx_uuid TEXT,
		source_path TEXT,
		dest_path TEXT,
		file_hash TEXT,
		file_size INTEGER,
		FOREIGN KEY(tx_uuid) REFERENCES transactions(uuid)
	)`)
	if err != nil {
		return nil, err
	}

	return &Ledger{db: db}, nil
}

func (l *Ledger) Begin() (*Transaction, error) {
	tx := &Transaction{
		UUID:      uuid.New().String(),
		Timestamp: time.Now(),
		State:     "Pending",
	}
	return tx, nil
}

func (l *Ledger) AddOperation(txUUID string, op Operation) error {
	_, err := l.db.Exec(`INSERT INTO operations 
		(tx_uuid, source_path, dest_path, file_hash, file_size) 
		VALUES (?, ?, ?, ?, ?)`,
		txUUID, op.SourcePath, op.DestPath, op.FileHash, op.FileSize)
	return err
}

func (l *Ledger) Save(tx *Transaction) error {
	_, err := l.db.Exec("INSERT OR REPLACE INTO transactions (uuid, timestamp, state) VALUES (?, ?, ?)",
		tx.UUID, tx.Timestamp, tx.State)
	return err
}

func (l *Ledger) Get(txUUID string) (*Transaction, error) {
	var tx Transaction
	err := l.db.QueryRow("SELECT uuid, timestamp, state FROM transactions WHERE uuid = ?", txUUID).
		Scan(&tx.UUID, &tx.Timestamp, &tx.State)
	if err != nil {
		return nil, err
	}

	rows, err := l.db.Query(`SELECT source_path, dest_path, file_hash, file_size 
		FROM operations WHERE tx_uuid = ?`, txUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var op Operation
		err := rows.Scan(&op.SourcePath, &op.DestPath, &op.FileHash, &op.FileSize)
		if err == nil {
			tx.Ops = append(tx.Ops, op)
		}
	}

	return &tx, nil
}

func (l *Ledger) GetRecentTransactions(limit int) ([]*Transaction, error) {
	rows, err := l.db.Query("SELECT uuid FROM transactions ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*Transaction
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err == nil {
			if tx, err := l.Get(uuid); err == nil {
				txs = append(txs, tx)
			}
		}
	}
	return txs, nil
}

func (l *Ledger) UpdateTransactionState(txUUID string, state string) error {
	_, err := l.db.Exec("UPDATE transactions SET state = ? WHERE uuid = ?", state, txUUID)
	return err
}

func (l *Ledger) Locate(pattern string) ([]Operation, error) {
	// Wir suchen nach Operationen, deren Transaktion NICHT 'RolledBack' ist
	query := `
		SELECT o.source_path, o.dest_path, o.file_hash, o.file_size 
		FROM operations o
		JOIN transactions t ON o.tx_uuid = t.uuid
		WHERE o.dest_path LIKE ? AND t.state != 'RolledBack'
		ORDER BY t.timestamp DESC`
	
	rows, err := l.db.Query(query, "%"+pattern+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Operation
	for rows.Next() {
		var op Operation
		if err := rows.Scan(&op.SourcePath, &op.DestPath, &op.FileHash, &op.FileSize); err == nil {
			results = append(results, op)
		}
	}
	return results, nil
}

func (l *Ledger) GetPathByHash(hash string) (string, error) {
	var path string
	// Wir suchen nach dem neuesten Pfad, der diesen Hash hat und noch existiert (stat-check erfolgt in der App)
	err := l.db.QueryRow("SELECT dest_path FROM operations WHERE file_hash = ? ORDER BY id DESC LIMIT 1", hash).
		Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}
