-- +goose Up
CREATE OR REPLACE VIEW order_summary_view AS
SELECT
    o.id AS order_id,
    c.phone AS customer_phone,
    p.address AS pickup_address,
    o.status,
    o.total_amount,
    o.pickup_time,
    COUNT(oi.id) AS items_count,
    COALESCE(SUM(oi.quantity), 0) AS total_quantity
FROM orders o
JOIN customer c ON c.id = o.customer_id
LEFT JOIN place p ON p.id = o.place_id
LEFT JOIN order_item oi ON oi.order_id = o.id
GROUP BY o.id, c.phone, p.address, o.status, o.total_amount, o.pickup_time;

CREATE OR REPLACE FUNCTION recalculate_order_total(p_order_id UUID)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    v_total BIGINT;
BEGIN
    SELECT COALESCE(SUM(quantity * unit_price), 0)
    INTO v_total
    FROM order_item
    WHERE order_id = p_order_id;

    UPDATE orders
    SET total_amount = v_total
    WHERE id = p_order_id;

    RETURN v_total;
END;
$$;

-- +goose Down
DROP FUNCTION IF EXISTS recalculate_order_total(UUID);
DROP VIEW IF EXISTS order_summary_view;
