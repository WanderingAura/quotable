package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Quotes QuoteDatabaseModel // change to corresponding interfaces when ready.
	Users  UserDatabaseModel
}

func New(db *sql.DB) Models {
	return Models{
		Quotes: QuoteDatabaseModel{DB: db},
		Users:  UserDatabaseModel{DB: db},
	}
}
