CREATE TABLE "users"
(
    "id"            serial       NOT NULL,
    "email"         varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL,
    "first_name"    varchar(255) NOT NULL,
    "last_name"     varchar(255) NOT NULL,
    "registered_at" TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- Product-Related Tables

CREATE TABLE "categories"
(
    "id"   serial       NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE
);

CREATE TABLE "images"
(
    "id"       serial       NOT NULL,
    "url"      varchar(255) NOT NULL UNIQUE,
    "alt_text" varchar(255)
);

CREATE TYPE language AS ENUM ('english', 'russian', 'ukrainian');

CREATE TABLE "product_titles"
(
    "id"       serial       NOT NULL,
    "title"    varchar(255) NOT NULL,
    "language" language     NOT NULL
);

CREATE TABLE "product_descriptions"
(
    "id"          serial       NOT NULL,
    "description" varchar(255) NOT NULL,
    "language"    language     NOT NULL
);

CREATE TABLE "products"
(
    "id"             serial                                   NOT NULL,
    "title_id"       int REFERENCES product_titles (id)       NOT NULL,
    "current_price"  float                                    NOT NULL,
    "previous_price" float,
    "code"           varchar(255),
    "description_id" int REFERENCES product_descriptions (id) NOT NULL,
    "image_id"       int REFERENCES images (id)               NOT NULL
);

-- Orders-Related Tables

CREATE TABLE "orders"
(
    "id"         serial                    NOT NULL,
    "user_id"    int REFERENCES users (id) NOT NULL,
    "ordered_at" TIMESTAMP                 NOT NULL DEFAULT NOW(),
    "status"     int                       NOT NULL,
    "address"    varchar(255)              NOT NULL
);

CREATE TABLE "order_items"
(
    "id"         serial                       NOT NULL,
    "order_id"   int REFERENCES orders (id)   NOT NULL,
    "product_id" int REFERENCES products (id) NOT NULL
);

-- Admin Users

CREATE TABLE "admin_users"
(
    "id"            serial       NOT NULL,
    "login"         varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL
);
