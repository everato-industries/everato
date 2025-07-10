-- name: CreateSuperUserIfNotExists :one
INSERT INTO super_users (username, name, email, password, role, permissions)
SELECT
    $1::VARCHAR,
    $2::VARCHAR,
    $3::VARCHAR,
    $4::VARCHAR,
    'SUPER_ADMIN'::SUPER_USER_ROLE,
    ARRAY[
        'MANAGE_EVENTS',
        'CREATE_EVENT',
        'EDIT_EVENT',
        'DELETE_EVENT',
        'VIEW_EVENT',
        'MANAGE_BOOKINGS',
        'CREATE_BOOKING',
        'EDIT_BOOKING',
        'DELETE_BOOKING',
        'VIEW_BOOKING',
        'MANAGE_USERS',
        'VIEW_REPORTS'
    ]::PERMISSIONS[]
WHERE NOT EXISTS (
    SELECT 1 FROM super_users WHERE username = $1::VARCHAR OR email = $3::VARCHAR
)
RETURNING id, email, password, role, permissions::text[], created_at, updated_at, username, name;

-- name: CreateAdminIfNotExists :one
INSERT INTO super_users (username, name, email, password, role, permissions)
SELECT
    $1,
    $2,
    $3,
    $4,
    $5::SUPER_USER_ROLE,
    $6::PERMISSIONS[]
WHERE NOT EXISTS (
    SELECT 1 FROM super_users WHERE username = $1 OR email = $3
)
RETURNING id, email, password, role, permissions, created_at, updated_at, username, name;

-- name: GetAdminUserRoles :many
SELECT e.enumlabel AS value
FROM pg_catalog.pg_type AS t
JOIN pg_catalog.pg_enum AS e ON t.oid = e.enumtypid
WHERE t.typname = 'super_user_role';

-- name: GetAdminPermissions :many
SELECT e.enumlabel AS value
FROM pg_catalog.pg_type AS t
JOIN pg_catalog.pg_enum AS e ON t.oid = e.enumtypid
WHERE t.typname = 'permissions';

-- name: GetAdminByEmail :one
SELECT id, email, password, role, permissions::text[], created_at, updated_at, username, name
FROM super_users
WHERE email = $1;

-- name: GetAdminByUsername :one
SELECT id, email, password, role, permissions::text[], created_at, updated_at, username, name
FROM super_users
WHERE username = $1;

-- name: GetAdminById :one
SELECT id, email, password, role, permissions::text[], created_at, updated_at, username, name
FROM super_users
WHERE id = $1;
