-- Make the slug column unique and not null in the "event" table
ALTER TABLE events
ALTER COLUMN slug SET NOT NULL;

ALTER TABLE events
ADD CONSTRAINT unique_slug UNIQUE (slug);
