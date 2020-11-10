CREATE TABLE text_blocks
(
    "id"        serial       NOT NULL UNIQUE,
    "english"   varchar(255) NOT NULL,
    "russian"   varchar(255) NOT NULL,
    "ukrainian" varchar(255) NOT NULL
);

CREATE TABLE homepage_images
(
    "id" serial       NOT NULL UNIQUE,
    "image_id" int REFERENCES images (id) NOT NULL ON DELETE CASCADE
);