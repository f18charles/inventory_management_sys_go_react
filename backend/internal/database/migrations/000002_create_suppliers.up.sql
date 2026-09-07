CREATE TABLE suppliers (
    id UUID PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_person TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);