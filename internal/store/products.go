package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	// TODO: Type - {Service, Item, File}
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) Create(ctx context.Context, product *Product) error {

	query := `
		INSERT INTO products (user_id, name, price, description, categories)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		product.UserID,
		product.Name,
		product.Price,
		product.Description,
		pq.Array(product.Categories),
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProductStore) GetByID(ctx context.Context, productID int64) (*Product, error) {
	query := `
		SELECT id, user_id, name, price, description, categories, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	var product Product

	err := s.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Price,
		&product.Description,
		pq.Array(&product.Categories),
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}
