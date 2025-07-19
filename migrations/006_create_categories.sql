-- Create categories table
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    image_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for parent_id
CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- Create index for slug
CREATE INDEX idx_categories_slug ON categories(slug);

-- Create index for active categories
CREATE INDEX idx_categories_active ON categories(is_active);

-- Add category_id to products table
ALTER TABLE products ADD COLUMN category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL;

-- Create index for category_id in products
CREATE INDEX idx_products_category_id ON products(category_id);
