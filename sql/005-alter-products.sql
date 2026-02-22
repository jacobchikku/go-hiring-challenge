-- 2. Modify the existing products table to link to categories
ALTER TABLE products 
ADD COLUMN category_id INTEGER REFERENCES categories(id);