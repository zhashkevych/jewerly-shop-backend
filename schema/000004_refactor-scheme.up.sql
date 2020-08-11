ALTER TABLE products
    DROP COLUMN title_id,
    DROP COLUMN description_id,
    DROP COLUMN material_id;

ALTER TABLE titles
    ADD COLUMN product_id int REFERENCES products (id) ON DELETE CASCADE;

ALTER TABLE descriptions
    ADD COLUMN product_id int REFERENCES products (id) ON DELETE CASCADE;

ALTER TABLE materials
    ADD COLUMN product_id int REFERENCES products (id) ON DELETE CASCADE;