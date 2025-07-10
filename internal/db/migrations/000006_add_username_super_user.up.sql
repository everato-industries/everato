-- Add the username field in the super user table
ALTER TABLE super_users
ADD COLUMN username VARCHAR(255) NOT NULL UNIQUE;
