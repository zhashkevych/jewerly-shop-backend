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

CREATE TABLE products
(
    "id"             serial                           NOT NULL UNIQUE,
    "title_id"       int REFERENCES titles (id)       NOT NULL,
    "current_price"  DECIMAL(10, 2)                   NOT NULL,
    "previous_price" DECIMAL(10, 2),
    "code"           varchar(255),
    "description_id" int REFERENCES descriptions (id) NOT NULL,
    "category_id"    int REFERENCES categories (id)   NOT NULL,
    "material_id"    int REFERENCES materials (id)    NOT NULL

);

CREATE TABLE product_images
(
    "id"         serial                       NOT NULL UNIQUE,
    "product_id" int REFERENCES products (id) NOT NULL,
    "image_id"   int REFERENCES images (id)   NOT NULL
);

-- Orders-Related Tables

CREATE TABLE orders
(
    "id"         serial                    NOT NULL UNIQUE,
    "user_id"    int REFERENCES users (id) NOT NULL,
    "ordered_at" TIMESTAMP                 NOT NULL DEFAULT NOW(),
    "status"     int                       NOT NULL,
    "address"    varchar(255)              NOT NULL
);

CREATE TABLE order_items
(
    "id"         serial                       NOT NULL UNIQUE,
    "order_id"   int REFERENCES orders (id)   NOT NULL,
    "product_id" int REFERENCES products (id) NOT NULL
);

-- Admin Users

CREATE TABLE admin_users
(
    "id"            serial       NOT NULL UNIQUE,
    "login"         varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL
);
