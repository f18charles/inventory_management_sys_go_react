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
    '$2a$12$3z1QkTCi7SovEAhX.X.fAePM.NQIY.EaFPEP2wPFaS4C6oO2uV2Wy',
    'admin',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (username) DO NOTHING;
