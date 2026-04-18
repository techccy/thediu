package memory

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var (
	dbPath = filepath.Join(os.Getenv("HOME"), ".ccy", "history.db")
)

type MemoryEntry struct {
	ID         int
	Hash       string
	Command    string
	Error      string
	FixCommand string
	Timestamp  time.Time
	Success    bool
}

type Memory struct {
	db *sql.DB
}

func NewMemory() (*Memory, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create memory directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Memory{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hash TEXT UNIQUE NOT NULL,
			command TEXT NOT NULL,
			error TEXT NOT NULL,
			fix_command TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			success BOOLEAN DEFAULT 0
		);

		CREATE INDEX IF NOT EXISTS idx_hash ON history(hash);
		CREATE INDEX IF NOT EXISTS idx_timestamp ON history(timestamp);
	`)
	return err
}

func (m *Memory) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *Memory) GenerateHash(command, error string) string {
	data := fmt.Sprintf("%s|%s", command, error)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (m *Memory) Save(command, errorMessage, fixCommand string, success bool) error {
	hash := m.GenerateHash(command, errorMessage)

	_, err := m.db.Exec(`
		INSERT OR REPLACE INTO history (hash, command, error, fix_command, success)
		VALUES (?, ?, ?, ?, ?)
	`, hash, command, errorMessage, fixCommand, success)

	return err
}

func (m *Memory) Find(command, errorMessage string) (*MemoryEntry, error) {
	hash := m.GenerateHash(command, errorMessage)

	var entry MemoryEntry
	var timestampStr string

	err := m.db.QueryRow(`
		SELECT id, hash, command, error, fix_command, timestamp, success
		FROM history
		WHERE hash = ? AND success = 1
		ORDER BY timestamp DESC
		LIMIT 1
	`, hash).Scan(&entry.ID, &entry.Hash, &entry.Command, &entry.Error, &entry.FixCommand, &timestampStr, &entry.Success)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	entry.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestampStr)
	return &entry, nil
}

func (m *Memory) GetAll() ([]MemoryEntry, error) {
	rows, err := m.db.Query(`
		SELECT id, hash, command, error, fix_command, timestamp, success
		FROM history
		ORDER BY timestamp DESC
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close rows: %v\n", closeErr)
		}
	}()

	var entries []MemoryEntry
	for rows.Next() {
		var entry MemoryEntry
		var timestampStr string

		if err := rows.Scan(&entry.ID, &entry.Hash, &entry.Command, &entry.Error, &entry.FixCommand, &timestampStr, &entry.Success); err != nil {
			return nil, err
		}

		entry.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestampStr)
		entries = append(entries, entry)
	}

	return entries, nil
}

func (m *Memory) Clear() error {
	_, err := m.db.Exec("DELETE FROM history")
	return err
}
