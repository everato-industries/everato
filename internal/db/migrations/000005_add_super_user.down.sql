-- Drop the table first since it depends on the ENUM types
DROP TABLE IF EXISTS super_users;

-- Then drop the ENUM types
DROP TYPE IF EXISTS PERMISSIONS;
DROP TYPE IF EXISTS SUPER_USER_ROLE;
