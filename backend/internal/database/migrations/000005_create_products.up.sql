CREATE TABLE products (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    unit_price BIGINT NOT NULL CHECK (unit_price >= 0),
    cost_price BIGINT NOT NULL CHECK (cost_price >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_is_active ON products(is_active);