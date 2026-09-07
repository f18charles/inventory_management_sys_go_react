CREATE TABLE purchase_items (
    id UUID PRIMARY KEY,
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_cost BIGINT NOT NULL CHECK (unit_cost >= 0),
    sub_total BIGINT NOT NULL CHECK (sub_total >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX idx_purchase_items_product_id ON purchase_items(product_id);