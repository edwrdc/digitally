package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrEditConflict      = errors.New("edit conflict ")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Products interface {
		GetByID(context.Context, int64) (*Product, error)
		Create(context.Context, *Product) error
		Delete(context.Context, int64) error
		Update(context.Context, *Product) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Reviews interface {
		GetByProductID(context.Context, int64) ([]Review, error)
		Create(context.Context, *Review) error
	}
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
		Reviews:  &ReviewStore{db},
	}
}
