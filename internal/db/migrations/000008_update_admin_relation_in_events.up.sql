ALTER TABLE events
DROP CONSTRAINT events_admin_id_fkey;

ALTER TABLE events
ADD CONSTRAINT events_admin_id_fkey
FOREIGN KEY (admin_id) REFERENCES super_users(id) ON DELETE CASCADE;
