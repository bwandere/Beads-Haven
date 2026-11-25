-- Migration: 004_ensure_admin_role.sql
-- Repair the seeded admin account when earlier migrations left it with a non-admin role.
INSERT INTO users (email, password_hash, role)
VALUES ('admin@beadshaven.co.ke', '$2a$10$wGPIp73PHOF3f9kVc2BrDuSlFZkNSbCw9yr1sO8AJf2wwNrUtZo8e', 'admin')
ON CONFLICT (email) DO UPDATE SET role = EXCLUDED.role;
