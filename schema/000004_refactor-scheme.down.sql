ALTER TABLE products
    ADD COLUMN title_id       int REFERENCES titles (id) ON DELETE CASCADE,
    ADD COLUMN description_id int REFERENCES descriptions (id) ON DELETE CASCADE,
    ADD COLUMN material_id    int REFERENCES materials (id) ON DELETE CASCADE;

ALTER TABLE titles
    DROP COLUMN product_id;

ALTER TABLE descriptions
    DROP COLUMN product_id;

ALTER TABLE materials
    DROP COLUMN product_id;