CREATE TABLE users
(
    "id"            serial       NOT NULL UNIQUE,
    "email"         varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL,
    "first_name"    varchar(255) NOT NULL,
    "last_name"     varchar(255) NOT NULL,
    "registered_at" TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- Product-Related Tables

CREATE TABLE categories
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE titles
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE descriptions
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE images
(
    "id"       serial       NOT NULL UNIQUE,
    "url"      varchar(255) NOT NULL UNIQUE,
    "alt_text" varchar(255)
);

CREATE TABLE materials
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE prices
(
    "id"  serial         NOT NULL UNIQUE,
    "usd" DECIMAL(10, 2) NOT NULL,
    "eur" DECIMAL(10, 2) NOT NULL,
    "uah" DECIMAL(10, 2) NOT NULL
);

CREATE TABLE products
(
    "id"             serial                                             NOT NULL UNIQUE,
    "title_id"       int REFERENCES titles (id) ON DELETE CASCADE       NOT NULL,
    "code"           varchar(255),
    "description_id" int REFERENCES descriptions (id) ON DELETE CASCADE NOT NULL,
    "category_id"    int REFERENCES categories (id) ON DELETE CASCADE   NOT NULL,
    "material_id"    int REFERENCES materials (id) ON DELETE CASCADE    NOT NULL,
    "price_id"       int REFERENCES prices (id) ON DELETE CASCADE       NOT NULL,
    "in_stock"       bool DEFAULT true
);

CREATE TABLE orders
(
    "id"              serial         NOT NULL UNIQUE,
    "user_id"         int REFERENCES users (id),
    "ordered_at"      TIMESTAMP      NOT NULL DEFAULT NOW(),
    "first_name"      varchar(255)   NOT NULL,
    "last_name"       varchar(255)   NOT NULL,
    "additional_name" varchar(255),
    "email"           varchar(255)   NOT NULL,
    "country"         varchar(255)   NOT NULL,
    "address"         varchar(255)   NOT NULL,
    "postal_code"     varchar(10)    NOT NULL,
    "total_cost"      DECIMAL(10, 2) NOT NULL
);

-- Orders-Related Tables

CREATE TABLE product_images
(
    "id"         serial                       NOT NULL UNIQUE,
    "product_id" int REFERENCES products (id) NOT NULL,
    "image_id"   int REFERENCES images (id)   NOT NULL
);

CREATE TABLE order_items
(
    "id"         serial                                         NOT NULL UNIQUE,
    "order_id"   int REFERENCES orders (id) ON DELETE CASCADE   NOT NULL,
    "product_id" int REFERENCES products (id) ON DELETE CASCADE NOT NULL,
    "quantity"   int                                            NOT NULL DEFAULT 1
);


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

-- Admin Panel

CREATE TABLE admin_users
(
    "id"            serial       NOT NULL UNIQUE,
    "login"         varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL
);

-- Data

INSERT INTO categories (english, russian, ukrainian)
VALUES ('Rings', 'Кольца', 'Кільця'),
       ('Bracelets', 'Браслеты', 'Браслети'),
       ('Pendants', 'Подвески', 'Кулони'),
       ('Earrings', 'Серьги', 'Сережки'),
       ('Necklaces', 'Колье', 'Кольє'),
       ('Sets', 'Наборы', 'Набори');