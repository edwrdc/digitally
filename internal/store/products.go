package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Version   int            `json:"version"`
	Reviews   []Review       `json:"reviews,omitempty"`
	User      User           `json:"user,omitempty"`
	Wishlist  []UserWishlist `json:"wishlist,omitempty"`
}

type UserFeedProduct struct {
	Product
	ReviewCount  int  `json:"review_count"`
	IsWishlisted bool `json:"is_wishlisted"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) GetUserFeed(ctx context.Context, userID int64) ([]UserFeedProduct, error) {
	query := `
		SELECT
			p.id AS product_id,
			p.user_id,
			u.username AS seller_username,
			p.name,
			p.price,
			p.description,
			p.categories,
			p.version,
			p.created_at,
			COALESCE(COUNT(r.id), 0) AS reviews_count,
			CASE WHEN w.product_id IS NOT NULL THEN true ELSE false END AS is_wishlisted
		FROM
			products p
			INNER JOIN users u ON u.id = p.user_id
			LEFT JOIN reviews r ON r.product_id = p.id
			LEFT JOIN user_wishlist w ON w.product_id = p.id AND w.user_id = $1
			GROUP BY p.id, u.username, w.product_id
			ORDER BY p.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []UserFeedProduct

	for rows.Next() {
		var product UserFeedProduct
		if err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.User.Username,
			&product.Name,
			&product.Price,
			&product.Description,
			pq.Array(&product.Categories),
			&product.Version,
			&product.CreatedAt,
			&product.ReviewCount,
			&product.IsWishlisted,
		); err != nil {
			return nil, err
		}

		feed = append(feed, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feed, nil
}

func (s *ProductStore) Create(ctx context.Context, product *Product) error {

	query := `
		INSERT INTO products (user_id, name, price, description, categories)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		SELECT id, user_id, name, price, description, categories, created_at, updated_at, version
		FROM products
		WHERE id = $1
	`
	var product Product

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Price,
		&product.Description,
		pq.Array(&product.Categories),
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Version,
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

func (s *ProductStore) Delete(ctx context.Context, productID int64) error {
	query := `
		DELETE FROM products WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *ProductStore) Update(ctx context.Context, product *Product) error {

	query := `
		UPDATE products 
		SET name = $1, price = $2, description = $3, categories = $4, updated_at = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Description,
		pq.Array(product.Categories),
		time.Now().UTC(),
		product.ID,
		product.Version,
	).Scan(&product.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		// TODO: Refactor this to use PG error codes
		case strings.Contains(err.Error(), "version"):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
