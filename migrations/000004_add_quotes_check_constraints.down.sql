ALTER TABLE quotes DROP CONSTRAINT IF EXISTS quotes_content_check;

ALTER TABLE quotes DROP CONSTRAINT IF EXISTS quotes_source_title_check;

ALTER TABLE quotes DROP CONSTRAINT IF EXISTS quotes_source_type_title_check;

ALTER TABLE quotes DROP CONSTRAINT IF EXISTS quotes_author_check;

ALTER TABLE quotes DROP CONSTRAINT IF EXISTS quotes_tags_length_check;