-- =====================================================
-- MIGRATION: Drop deprecated Product fields
-- Date: 2026-01-08
-- Reason: Refactor - Price/Stock/SKU moved to ProductItem
-- =====================================================

-- BACKUP FIRST (recommended)
-- pg_dump -U postgres -d product_service -t products > products_backup_$(date +%Y%m%d).sql

-- Drop deprecated columns from products table
ALTER TABLE products DROP COLUMN IF EXISTS price;
ALTER TABLE products DROP COLUMN IF EXISTS sku;
ALTER TABLE products DROP COLUMN IF EXISTS stock;

-- Verify schema after migration
-- \d products

-- =====================================================
-- ROLLBACK (if needed)
-- =====================================================
-- ALTER TABLE products ADD COLUMN price DECIMAL(15,2);
-- ALTER TABLE products ADD COLUMN sku VARCHAR(255) UNIQUE;
-- ALTER TABLE products ADD COLUMN stock INTEGER DEFAULT 0;
-- =====================================================

-- Expected final schema for products:
-- ✅ id (PK)
-- ✅ shop_id (FK)
-- ✅ category_id (FK, nullable, leaf category only)
-- ✅ name
-- ✅ description
-- ✅ base_price (giá tham chiếu)
-- ✅ is_active
-- ✅ sold_count
-- ✅ status
-- ✅ images (jsonb)
-- ✅ created_at
-- ✅ updated_at
