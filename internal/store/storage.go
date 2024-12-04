package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Products interface {
		Create(context.Context, *Product) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
	}
}
