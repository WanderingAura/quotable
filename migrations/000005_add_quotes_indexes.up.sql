CREATE INDEX IF NOT EXISTS quotes_content_idx ON quotes USING GIN (to_tsvector('simple', content));
CREATE INDEX IF NOT EXISTS quotes_tags_idx ON quotes USING GIN (tags);