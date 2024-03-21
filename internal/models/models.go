package models

import "database/sql"

type Models struct {
	Quotes QuoteDatabaseModel
}

func New(db *sql.DB) Models {
	return Models{
		Quotes: QuoteDatabaseModel{DB: db},
	}
}
