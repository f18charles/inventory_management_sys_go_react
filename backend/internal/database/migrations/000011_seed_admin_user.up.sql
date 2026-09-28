-- Seed: one admin user for local development.
-- Password is "admin1234" hashed with bcrypt (cost 12).
-- Change this immediately after first login in any real environment.
INSERT INTO users (id, first_name, last_name, username, email, pass_hash, role, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Admin',
    'User',
    'admin',
    'admin@example.com',
    '$2a$12$J61ahu6c5wMB/vPSG7kK7.k2XAFMQqWT79OsWazk2Lc6wFyNOSzx.',
    'admin',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (username) WHERE deleted_at IS NULL DO NOTHING;
