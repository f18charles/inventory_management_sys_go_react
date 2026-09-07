CREATE TABLE sales (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    sale_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'cancelled', 'refunded')),
    total_amount BIGINT NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'card', 'mobile_money', 'bank_transfer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sales_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_user_id ON sales(user_id);
CREATE INDEX idx_sales_status ON sales(status);