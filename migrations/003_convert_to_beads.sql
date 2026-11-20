-- Migration: 003_convert_to_beads.sql
-- Idempotent conversion: removes old sneaker & jersey seed records and prevents image duplicates

-- 1. Ensure unique constraint on product_images (product_id, url)
DO $$
BEGIN
    -- Delete any existing duplicates before adding constraint
    DELETE FROM product_images a
    USING product_images b
    WHERE a.id > b.id
      AND a.product_id = b.product_id
      AND a.url = b.url;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uq_product_images_pid_url'
    ) THEN
        ALTER TABLE product_images ADD CONSTRAINT uq_product_images_pid_url UNIQUE (product_id, url);
    END IF;
END $$;

-- 2. Remove legacy sneaker and jersey products and associated records
DELETE FROM products WHERE id IN (
    'volt-runner',
    'court-classic-low',
    'night-high',
    'retro-wave',
    'volt-runner-pro',
    'street-court',
    'dune-trainer',
    'blackout-mid',
    'harambee-home',
    'royal-blues-home',
    'red-stripes-classic',
    'lions-away'
);

-- 3. Ensure valid bcrypt hash for admin: admin@amazingoutfits.co.ke / admin123
UPDATE users
SET password_hash = '$2a$10$wGPIp73PHOF3f9kVc2BrDuSlFZkNSbCw9yr1sO8AJf2wwNrUtZo8e'
WHERE email = 'admin@amazingoutfits.co.ke';
