ALTER TABLE order_items
    DROP COLUMN quantity;

ALTER TABLE orders
    DROP COLUMN first_name,
    DROP COLUMN last_name,
    DROP COLUMN additional_name,
    DROP COLUMN country,
    DROP COLUMN email,
    DROP COLUMN postal_code;

DROP TABLE transactions;
DROP TABLE transactions_history;