-- Migration: 002_fix_refresh_tokens.sql
-- Description: Drop stale company_id column from refresh_tokens.
--              The column was added by GORM AutoMigrate when the RefreshToken model
--              had a CompanyID field. The field was later removed from the model,
--              but AutoMigrate never drops columns, leaving a NOT NULL constraint
--              that breaks token creation for users without a company.
-- Date: 2026-02-18

BEGIN;

ALTER TABLE refresh_tokens DROP COLUMN IF EXISTS company_id;

COMMIT;
