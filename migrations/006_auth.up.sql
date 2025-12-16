-- Add credential support and keep teams in sync with org units
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS password_hash text NOT NULL DEFAULT '';

-- Ensure org units get created for any existing teams
INSERT INTO org_units (id, org_id, parent_unit_id, name, path, unit_type, created_at, updated_at)
SELECT t.id, t.org_id, t.parent_team_id, t.name, t.path, 'team', t.created_at, t.updated_at
FROM teams t
ON CONFLICT (id) DO NOTHING;
