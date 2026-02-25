-- +goose Up
CREATE TABLE order_status
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO order_status (id, name)
VALUES (0, 'unspecified'),
       (1, 'draft'),
       (2, 'awaiting_payment'),
       (3, 'paid'),
       (4, 'in_progress'),
       (5, 'ready'),
       (6, 'completed'),
       (7, 'cancelled'),
       (8, 'failed');

-- Изменяем тип колонки status в таблице orders на INTEGER
-- Сначала удалим существующие значения и преобразуем их на лету.
-- Предполагаем, что старые текстовые значения соответствуют новым именам.

ALTER TABLE orders
    ALTER COLUMN status TYPE INTEGER USING (CASE
                                                 WHEN status = 'unspecified' THEN 0
                                                 WHEN status = 'draft' THEN 1
                                                 WHEN status = 'awaiting_payment' THEN 2
                                                 WHEN status = 'paid' THEN 3
                                                 WHEN status = 'in_progress' THEN 4
                                                 WHEN status = 'ready' THEN 5
                                                 WHEN status = 'completed' THEN 6
                                                 WHEN status = 'cancelled' THEN 7
                                                 WHEN status = 'failed' THEN 8
                                                 ELSE 0
        END);

-- +goose Down
ALTER TABLE orders
    ALTER COLUMN status TYPE TEXT USING (CASE
                                             WHEN status = 0 THEN 'unspecified'
                                             WHEN status = 1 THEN 'draft'
                                             WHEN status = 2 THEN 'awaiting_payment'
                                             WHEN status = 3 THEN 'paid'
                                             WHEN status = 4 THEN 'in_progress'
                                             WHEN status = 5 THEN 'ready'
                                             WHEN status = 6 THEN 'completed'
                                             WHEN status = 7 THEN 'cancelled'
                                             WHEN status = 8 THEN 'failed'
                                             ELSE 'unspecified'
        END);

DROP TABLE order_status;
