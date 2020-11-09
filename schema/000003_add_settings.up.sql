CREATE TABLE text
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE homepage_settings
(
    "image_id" int REFERENCES images (id) NOT NULL,
    "button_text_id" int REFERENCES text (id) NOT NULL,
    "text_block_1_id" int REFERENCES text (id) NOT NULL,
    "text_block_2_id" int REFERENCES text (id) NOT NULL
);

CREATE TABLE homepage_product_item
(
    "id"        serial       NOT NULL UNIQUE,
    "product_id" int REFERENCES products (id) NOT NULL
);