-- name: CreateSuperUserIfNotExists :one
INSERT INTO super_users (username, email, password, role, permissions)
SELECT
    $1::VARCHAR,
    $2::VARCHAR,
    $3::VARCHAR,
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
    SELECT 1 FROM super_users WHERE username = $1::VARCHAR OR email = $2::VARCHAR
)
RETURNING id, email, password, role, permissions::text[], created_at, updated_at, username;
