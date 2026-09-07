CREATE TABLE suppliers (
    id UUID PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_person TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_suppliers_email_unique ON suppliers(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_deleted_at ON suppliers(deleted_at);