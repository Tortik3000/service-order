-- +goose Up
ALTER TABLE menu_item ADD COLUMN image_url TEXT;

-- +goose Down
ALTER TABLE menu_item DROP COLUMN image_url;
