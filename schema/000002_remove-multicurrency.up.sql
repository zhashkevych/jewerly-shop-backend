ALTER TABLE products ADD COLUMN price DECIMAL(10, 2);
UPDATE products p SET price=pr.usd FROM prices pr WHERE p.price_id = pr.id;
ALTER TABLE products ALTER COLUMN price SET NOT NULL;

ALTER TABLE products DROP COLUMN price_id;
DROP TABLE prices;