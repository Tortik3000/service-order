-- +goose Up
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE place
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address         TEXT NOT NULL
);


-- +goose Down
DROP TABLE place;