CREATE TABLE inventories (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reorder_level INT NOT NULL DEFAULT 0 CHECK (reorder_level >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_inventories_product_id_unique ON inventories(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventories_deleted_at ON inventories(deleted_at);