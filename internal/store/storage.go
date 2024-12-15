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
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type Storage struct {
	Products interface {
		GetByID(context.Context, int64) (*Product, error)
		Create(context.Context, *Product) error
		Delete(context.Context, int64) error
		Update(context.Context, *Product) error
		GetUserFeed(context.Context, int64, PaginationFeedQuery) ([]UserFeedProduct, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, expiry time.Duration) error
		Activate(context.Context, string) error
	}
	Reviews interface {
		GetByProductID(context.Context, int64) ([]Review, error)
		Create(context.Context, *Review) error
	}
	Wishlist interface {
		Add(ctx context.Context, userID, productID int64) error
		Remove(ctx context.Context, userID, productID int64) error
	}
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
		Reviews:  &ReviewStore{db},
		Wishlist: &WishlistStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
