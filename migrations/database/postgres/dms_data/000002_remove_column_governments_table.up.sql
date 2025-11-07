-- Migration: remove subject_group_id column from governments
-- Up: drop the column (use IF EXISTS to be safe)
ALTER TABLE "governments"
	DROP COLUMN IF EXISTS "subject_group_id";

-- Note: any dependent foreign key/constraint on this column will be removed by dropping the column.