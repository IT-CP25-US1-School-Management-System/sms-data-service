
-- Migration: restore subject_group_id column on governments
-- Down: re-add the column and recreate the foreign key referencing departments(id)
ALTER TABLE "governments"
	ADD COLUMN IF NOT EXISTS "subject_group_id" int;

-- Recreate foreign key. Use a generated constraint name by omitting the CONSTRAINT clause.
ALTER TABLE "governments"
	ADD FOREIGN KEY ("subject_group_id") REFERENCES "departments" ("id");

-- Note: if in your schema the intended reference is different (e.g. subjects_groups), update the REFERENCES accordingly.

