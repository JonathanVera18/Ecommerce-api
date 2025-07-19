-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Order details
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    subtotal_amount DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    shipping_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    
    -- Payment information
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(20),
    payment_id VARCHAR(255),
    paid_at TIMESTAMP,
    
    -- Shipping information
    shipping_first_name VARCHAR(100) NOT NULL,
    shipping_last_name VARCHAR(100) NOT NULL,
    shipping_email VARCHAR(255),
    shipping_phone VARCHAR(20),
    shipping_street VARCHAR(255) NOT NULL,
    shipping_city VARCHAR(100) NOT NULL,
    shipping_state VARCHAR(100) NOT NULL,
    shipping_country VARCHAR(100) NOT NULL,
    shipping_postal_code VARCHAR(20) NOT NULL,
    
    -- Billing information (optional)
    billing_first_name VARCHAR(100),
    billing_last_name VARCHAR(100),
    billing_email VARCHAR(255),
    billing_phone VARCHAR(20),
    billing_street VARCHAR(255),
    billing_city VARCHAR(100),
    billing_state VARCHAR(100),
    billing_country VARCHAR(100),
    billing_postal_code VARCHAR(20),
    
    -- Tracking information
    tracking_number VARCHAR(100),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    
    -- Additional information
    notes TEXT,
    internal_notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create carts table
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
CREATE INDEX IF NOT EXISTS idx_carts_customer_id ON carts(customer_id);
CREATE INDEX IF NOT EXISTS idx_carts_deleted_at ON carts(deleted_at);

-- Add constraints
ALTER TABLE orders ADD CONSTRAINT chk_orders_status CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded'));
ALTER TABLE orders ADD CONSTRAINT chk_orders_payment_status CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded', 'cancelled'));
ALTER TABLE orders ADD CONSTRAINT chk_orders_payment_method CHECK (payment_method IN ('card', 'paypal', 'bank_transfer', 'cash_on_delivery') OR payment_method IS NULL);
ALTER TABLE orders ADD CONSTRAINT chk_orders_amounts CHECK (total_amount >= 0 AND subtotal_amount >= 0 AND tax_amount >= 0 AND shipping_amount >= 0 AND discount_amount >= 0);
