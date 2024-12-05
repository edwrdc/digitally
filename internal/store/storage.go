package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Products interface {
		GetByID(context.Context, int64) (*Product, error)
		Create(context.Context, *Product) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Reviews interface {
		GetByProductID(context.Context, int64) ([]Review, error)
	}
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
		Reviews:  &ReviewStore{db},
	}
}
