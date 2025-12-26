package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	DB *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &SQLiteStore{DB: db}, nil
}

func (s *SQLiteStore) Close() error {
	if s.DB == nil {
		return nil
	}
	return s.DB.Close()
}

func (s *SQLiteStore) Migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS trades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    side TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL NOT NULL,
    time TEXT NOT NULL,
    strategy TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS positions (
    code TEXT PRIMARY KEY,
    quantity INTEGER NOT NULL,
    avg_price REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS daily_pnl (
    date TEXT PRIMARY KEY,
    equity REAL NOT NULL,
    profit REAL NOT NULL,
    drawdown REAL NOT NULL
);
`
	_, err := s.DB.Exec(schema)
	return err
}
