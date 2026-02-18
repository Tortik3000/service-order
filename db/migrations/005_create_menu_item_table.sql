-- +goose Up
CREATE TABLE menu_item
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID REFERENCES menu_category(id) NOT NULL,
    name        TEXT NOT NULL unique,
    description TEXT,
    price       BIGINT NOT NULL,
    active      BOOLEAN DEFAULT TRUE
);

-- +goose Down
DROP TABLE menu_item;