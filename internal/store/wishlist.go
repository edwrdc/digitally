package store

import (
	"context"
	"database/sql"
	"time"
)

type UserWishlist struct {
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WishlistStore struct {
	db *sql.DB
}

func (s *WishlistStore) Add(ctx context.Context, userID, productID int64) error {
	query := `INSERT INTO user_wishlist (user_id, product_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_wishlist_pkey"`:
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *WishlistStore) Remove(ctx context.Context, userID, productID int64) error {
	query := `DELETE FROM user_wishlist WHERE user_id = $1 AND product_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return err
	}
	return nil
}
