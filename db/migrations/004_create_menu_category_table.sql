-- +goose Up
CREATE TABLE menu_category
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL unique,
    sort_order  INT NOT NULL
);

-- +goose Down
DROP TABLE menu_category;