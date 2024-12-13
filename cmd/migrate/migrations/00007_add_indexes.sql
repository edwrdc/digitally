-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_product_name ON products USING gin (name gin_trgm_ops);
CREATE INDEX idx_product_description ON products USING gin (description gin_trgm_ops);
CREATE INDEX idx_product_categories ON products USING gin (categories);

CREATE INDEX idx_review_comment ON reviews using gin (comment gin_trgm_ops);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_product_user_id ON products (user_id);
CREATE INDEX idx_wishlist_user_id ON user_wishlist (user_id);
CREATE INDEX idx_reviews_product_id ON reviews (product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_name;
DROP INDEX IF EXISTS idx_product_description;
DROP INDEX IF EXISTS idx_product_categories;
DROP INDEX IF EXISTS idx_review_comment;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_product_user_id;
DROP INDEX IF EXISTS idx_wishlist_user_id;
DROP INDEX IF EXISTS idx_reviews_product_id;
-- +goose StatementEnd
