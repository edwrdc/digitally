package store

import (
	"context"
	"database/sql"
	"time"
)

type Review struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
}

type ReviewStore struct {
	DB *sql.DB
}

func (s *ReviewStore) GetByProductID(ctx context.Context, productID int64) ([]Review, error) {

	query := `
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.created_at, users.username, users.id FROM reviews r
		JOIN users on users.id = r.user_id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]Review, 0)
	for rows.Next() {
		var r Review
		r.User = User{}
		err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.ProductID,
			&r.Rating,
			&r.Comment,
			&r.CreatedAt,
			&r.User.Username,
			&r.User.ID,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func (s *ReviewStore) Create(ctx context.Context, review *Review) error {
	query := `
		INSERT INTO reviews (user_id, product_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return s.DB.QueryRowContext(
		ctx,
		query,
		review.UserID,
		review.ProductID,
		review.Rating,
		review.Comment,
	).Scan(
		&review.ID,
		&review.CreatedAt,
	)
}
