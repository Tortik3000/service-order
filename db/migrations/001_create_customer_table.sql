-- +goose Up
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE customer
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone         TEXT NOT NULL
);


-- +goose Down
DROP TABLE customer;