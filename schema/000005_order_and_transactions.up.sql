ALTER TABLE order_items
    ADD COLUMN quantity int NOT NULL DEFAULT 1;

ALTER TABLE orders
    DROP COLUMN user_id,
    DROP COLUMN status,
    ADD COLUMN user_id         int REFERENCES users (id),
    ADD COLUMN first_name      varchar(255)   NOT NULL,
    ADD COLUMN last_name       varchar(255)   NOT NULL,
    ADD COLUMN additional_name varchar(255),
    ADD COLUMN country         varchar(255)   NOT NULL,
    ADD COLUMN email           varchar(255)   NOT NULL,
    ADD COLUMN postal_code     varchar(10)    NOT NULL,
    ADD COLUMN total_cost      DECIMAL(10, 2) NOT NULL;

CREATE TABLE transactions
(
    "id"        serial                                       NOT NULL UNIQUE,
    "order_id"  int REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    "uuid"      uuid                                         NOT NULL,
    "status"    varchar(255)                                 NOT NULL DEFAULT 'created',
    "card_mask" varchar(255),
    "paid_date" timestamp
);