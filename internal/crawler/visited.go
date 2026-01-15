package crawler

import (
	"database/sql"
)

type VisitedSet struct {
	db *sql.DB
}

func NewVisitedSet(db *sql.DB) *VisitedSet {
	return &VisitedSet{db: db}
}

func (s *VisitedSet) Add(url string) bool {
	query := `INSERT INTO crawled_urls (url) VALUES($1) ON CONFLICT (url) DO NOTHING`
	result, err := s.db.Exec(query, url)
	if err != nil {
		return false
	}
	rows, _ := result.RowsAffected()
	return rows > 0
}

func (s *VisitedSet) Len() int {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM crawled_urls").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (s *VisitedSet) UpdateStatus(url string, statusCode int, errMsg string) error {
	query := `
    INSERT INTO crawled_urls (url, status_code, error_message)
    VALUES ($1, $2, $3)
    ON CONFLICT (url) DO UPDATE SET
        status_code = EXCLUDED.status_code,
        error_message = EXCLUDED.error_message,
        crawled_at = CURRENT_TIMESTAMP`
	_, err := s.db.Exec(query, url, statusCode, errMsg)
	return err
}
