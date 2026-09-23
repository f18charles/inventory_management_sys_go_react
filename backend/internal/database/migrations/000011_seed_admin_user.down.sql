-- Removes the seeded admin user (matched by username to avoid UUID dependency).
DELETE FROM users WHERE username = 'admin' AND email = 'admin@example.com';
