-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    short_description VARCHAR(500),
    sku VARCHAR(100) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    compare_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    
    -- Inventory
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    low_stock_level INTEGER DEFAULT 10,
    track_inventory BOOLEAN DEFAULT true,
    allow_backorders BOOLEAN DEFAULT false,
    
    -- Organization
    category VARCHAR(50) NOT NULL,
    tags VARCHAR(1000),
    brand VARCHAR(100),
    
    -- Physical properties
    weight DECIMAL(8,3),
    length DECIMAL(8,2),
    width DECIMAL(8,2),
    height DECIMAL(8,2),
    
    -- SEO
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    slug VARCHAR(255) UNIQUE NOT NULL,
    
    -- Status and visibility
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    featured BOOLEAN DEFAULT false,
    visible BOOLEAN DEFAULT true,
    
    -- Analytics
    view_count INTEGER DEFAULT 0,
    
    -- Seller
    seller_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_featured ON products(featured);
CREATE INDEX IF NOT EXISTS idx_products_visible ON products(visible);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_product_images_deleted_at ON product_images(deleted_at);

-- Add constraints
ALTER TABLE products ADD CONSTRAINT chk_products_category CHECK (category IN ('electronics', 'clothing', 'books', 'home', 'sports', 'toys', 'beauty', 'food', 'other'));
ALTER TABLE products ADD CONSTRAINT chk_products_status CHECK (status IN ('draft', 'active', 'inactive', 'deleted'));
ALTER TABLE products ADD CONSTRAINT chk_products_price CHECK (price >= 0);
ALTER TABLE products ADD CONSTRAINT chk_products_stock CHECK (stock_quantity >= 0);
