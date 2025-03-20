package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
	pq squirrel.StatementBuilderType
}

func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &Store{db: db, pq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	if err := store.initTables(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		email VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tokens(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id),
		expiry TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS requests(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id),
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := s.db.Exec(query)
	return err
}
