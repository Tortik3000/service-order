-- +goose Up
ALTER TABLE customer ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE customer DROP COLUMN name;
