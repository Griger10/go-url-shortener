package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"url-shortener/cmd/internal/storage"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS url (
			id    INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url   TEXT NOT NULL
		);
	`); err != nil {
		return nil, fmt.Errorf("%s: create table: %w", op, err)
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`); err != nil {
		return nil, fmt.Errorf("%s: create index: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	res, err := s.db.Exec(`INSERT INTO url(url, alias) VALUES(?, ?)`, urlToSave, alias)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("storage.sqlite.SaveURL: %w", storage.ErrURLExists)
		}
		return 0, fmt.Errorf("storage.sqlite.SaveURL: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	var resURL string
	if err := s.db.QueryRow(`SELECT url FROM url WHERE alias = ?`, alias).Scan(&resURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: query: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	res, err := s.db.Exec(`DELETE FROM url WHERE alias = ?`, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return storage.ErrURLNotFound
	}

	return nil
}
