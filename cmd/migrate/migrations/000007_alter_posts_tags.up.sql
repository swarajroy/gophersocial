ALTER TABLE IF EXISTS posts
  ALTER COLUMN tags
  TYPE text[]
  USING tags::text[];