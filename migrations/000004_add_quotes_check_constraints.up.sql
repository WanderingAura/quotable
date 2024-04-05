ALTER TABLE quotes ADD CONSTRAINT quotes_content_check CHECK (LENGTH(content) < 300);

ALTER TABLE quotes ADD CONSTRAINT quotes_author_check CHECK (LENGTH(author) < 100);

ALTER TABLE quotes ADD CONSTRAINT quotes_source_title_check CHECK (LENGTH(source_title) < 300);

ALTER TABLE quotes ADD CONSTRAINT quotes_source_type_title_check CHECK (LENGTH(source_type) < 300);

ALTER TABLE quotes ADD CONSTRAINT quotes_tags_length_check CHECK (array_length(tags, 1) BETWEEN 0 AND 10);
