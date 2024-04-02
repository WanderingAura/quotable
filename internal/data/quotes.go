package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/WanderingAura/quotable/internal/validator"
	"github.com/lib/pq"
)

type Quote struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	Author       string    `json:"author"`
	Source       Source    `json:"source,omitempty"`
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

func (s *Source) isPartial() bool {
	return (s.Title == "" || s.Type == "") && !(s.Title == "" && s.Type == "")
}

func ValidateQuote(v *validator.Validator, quote *Quote) {
	v.Check(quote.Author != "", "author", "author must be provided")
	v.Check(len(quote.Author) <= 100, "author", "author must be less than 100 bytes")

	v.Check(!quote.Source.isPartial(), "source", "either provide both source title and type or provide neither")

	v.Check(quote.Tags != nil, "tags", "must be provided")
	v.Check(len(quote.Tags) >= 1, "tags", "must contain at least one tag")
	v.Check(len(quote.Tags) <= 5, "tags", "must not contain more than 5 tags")

	v.Check(validator.Unique(quote.Tags), "tags", "must not contain duplicate values")
}

func (m *QuoteDatabaseModel) Get(id int64) (*Quote, error) {
	query := `
	SELECT id, created_at, last_modified, user_id, content, author, source_title, source_type, tags, version
	FROM quotes
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var quote Quote

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&quote.ID,
		&quote.CreatedAt,
		&quote.LastModified,
		&quote.UserID,
		&quote.Content,
		&quote.Author,
		&quote.Source.Title,
		&quote.Source.Type,
		pq.Array(&quote.Tags),
		&quote.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &quote, nil
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

func (m *QuoteDatabaseModel) GetAll(content string, tags []string, filters Filters) ([]*Quote, Metadata, error) {

	// if title or genre is empty then the WHERE conditions default to true
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), quotes.id, quotes.created_at, quotes.last_modified, users.user_id, users.username, 
		quotes.content, quotes.author, quotes.source_title, quotes.source_type, quotes.tags, quotes.version
		FROM quotes INNER JOIN users ON users.id = quotes.user_id
		WHERE (to_tsvector('english', quotes.content) @@ plainto_tsquery('english', $1) OR $1 = '')
		AND (quotes.tags @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{content, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	quotes := []*Quote{}

	var totalRecords int

	for rows.Next() {
		var quote Quote
		err := rows.Scan(
			&totalRecords,
			&quote.ID,
			&quote.CreatedAt,
			&quote.LastModified,
			&quote.UserID,
			&quote.Content,
			&quote.Author,
			&quote.Source.Title,
			&quote.Source.Type,
			pq.Array(&quote.Tags),
			&quote.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		quotes = append(quotes, &quote)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return quotes, metadata, nil
}

func (m *QuoteDatabaseModel) GetAllForUser(userID int64, content string, tags []string, filters Filters) ([]*Quote, Metadata, error) {

	// if title or genre is empty then the WHERE conditions default to true
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), quotes.id, quotes.created_at, quotes.last_modified, users.user_id, users.username, 
		quotes.content, quotes.author, quotes.source_title, quotes.source_type, quotes.tags, quotes.version
		FROM quotes INNER JOIN users ON users.id = quotes.user_id
		WHERE users.id = $1
		AND(to_tsvector('english', quotes.content) @@ plainto_tsquery('english', $2) OR $2 = '')
		AND (quotes.tags @> $3 OR $3 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{content, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	quotes := []*Quote{}

	var totalRecords int

	for rows.Next() {
		var quote Quote
		err := rows.Scan(
			&totalRecords,
			&quote.ID,
			&quote.CreatedAt,
			&quote.LastModified,
			&quote.UserID,
			&quote.Content,
			&quote.Author,
			&quote.Source.Title,
			&quote.Source.Type,
			pq.Array(&quote.Tags),
			&quote.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		quotes = append(quotes, &quote)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return quotes, metadata, nil
}
