-- +goose Up
CREATE TABLE order_item
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id     UUID REFERENCES orders (id) NOT NULL,
    menu_item_id UUID REFERENCES menu_item(id) NOT NULL,
    quantity     INT NOT NULL,
    unit_price   BIGINT NOT NULL
);

-- +goose Down
DROP TABLE order_item;