-- +goose Up
ALTER TABLE orders ADD COLUMN pick_up BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE orders DROP COLUMN pick_up;
