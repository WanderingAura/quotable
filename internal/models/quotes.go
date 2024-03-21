package models

import "database/sql"

type Quote struct {
	ID         int64
	Content    string
	Source     string
	SourceType string
	Tags       []string
}

type QuoteModel interface {
	Insert(quote *Quote) error
	Get(id int64) (*Quote, error)
	Latest() ([]*Quote, error)
}

type QuoteDatabaseModel struct {
	DB *sql.DB
}
