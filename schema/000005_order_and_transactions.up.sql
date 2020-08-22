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
    "id"       serial                                       NOT NULL UNIQUE,
    "order_id" int REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    "uuid"     uuid                                         NOT NULL
);

CREATE TABLE transactions_history
(
    "id"         serial       NOT NULL UNIQUE,
    "uuid"       uuid         NOT NULL,
    "card_mask"  varchar(255),
    "status"     varchar(255) NOT NULL DEFAULT 'created',
    "created_at" timestamp    NOT NULL DEFAULT NOW()
);