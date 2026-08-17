CREATE TABLE suppliers(
    id TEXT PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_person TEXT NOT NULL,
    phone VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
);