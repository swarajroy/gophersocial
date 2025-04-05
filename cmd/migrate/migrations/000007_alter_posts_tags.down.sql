ALTER TABLE IF EXISTS posts
  ALTER COLUMN tags
  TYPE varchar(100)[]
  USING tags::varchar(100)[];