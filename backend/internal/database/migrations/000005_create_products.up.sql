CREATE TABLE products  (
    id UUID PRIMARY KEY,
    category_id uuid NOT NULL,
    sku VARCHAR UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    unit_price INTEGER NOT NULL,
    cost_price INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);