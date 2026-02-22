-- +goose Up
CREATE TABLE orders
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id  UUID REFERENCES customer (id) NOT NULL,
    status       TEXT                          NOT NULL,
    total_amount BIGINT                        NOT NULL,
    pickup_time  TIME,
    place_id     UUID REFERENCES place (id)
);

-- +goose Down
DROP TABLE orders;