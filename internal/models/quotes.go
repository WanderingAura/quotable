package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Quote struct {
	ID           int64     `json:"id"`
	LastModified time.Time `json:"last_modified"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	Author       string    `json:"author,omitempty"`
	Source       Source    `json:"source"`
	Tags         []string  `json:"tags"`
	Version      int       `json:"version"`
}

// TODO: make the source type marhsal JSON and unmarshal using the format sourceTitle(sourceType)?
type Source struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}

type QuoteModel interface {
	Insert(quote *Quote) error
	Get(id int64) (*Quote, error)
	Update(quote *Quote) error
	Latest() ([]*Quote, error)
}

type QuoteDatabaseModel struct {
	DB *sql.DB
}

func (m *QuoteDatabaseModel) Insert(quote *Quote) error {

	query := `
		INSERT INTO quotes (user_id, content, author, source_title, source_type, tags)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	args := []interface{}{quote.UserID, quote.Content, quote.Author, quote.Source.Title, quote.Source.Type, pq.Array(quote.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&quote.ID,
		&quote.LastModified,
		&quote.Version,
	)
}

func (m *QuoteDatabaseModel) Update(quote *Quote) error {
	query := `
		UPDATE quotes
		SET content=$1, author=$2, source_title=$3, source_type=$4, tags=$5, version=version+1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{quote.Content, quote.Author, quote.Source.Title, quote.Source.Type, quote.Tags, quote.ID, quote.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
