-- RBAC and org unit scaffolding
CREATE TABLE IF NOT EXISTS org_units (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id          uuid NOT NULL REFERENCES orgs(id),
  parent_unit_id  uuid NULL REFERENCES org_units(id),
  name            text NOT NULL,
  path            text NOT NULL,
  unit_type       text NOT NULL DEFAULT 'team',
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  deleted_at      timestamptz NULL,
  UNIQUE (org_id, name)
);

CREATE INDEX IF NOT EXISTS idx_org_units_org ON org_units(org_id);
CREATE INDEX IF NOT EXISTS idx_org_units_parent ON org_units(parent_unit_id);

-- seed org units from existing teams to maintain continuity
INSERT INTO org_units (id, org_id, parent_unit_id, name, path, unit_type, created_at, updated_at)
SELECT t.id, t.org_id, t.parent_team_id, t.name, t.path, 'team', t.created_at, t.updated_at
FROM teams t
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS roles (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid NOT NULL REFERENCES orgs(id),
  org_unit_id uuid NULL REFERENCES org_units(id),
  name        text NOT NULL,
  description text NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz NULL,
  UNIQUE (org_id, name, org_unit_id)
);

CREATE TABLE IF NOT EXISTS permissions (
  code        text PRIMARY KEY,
  description text NULL,
  scope       text NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id         uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_code text NOT NULL REFERENCES permissions(code),
  created_at      timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (role_id, permission_code)
);

-- membership enhancements
ALTER TABLE memberships
  ADD COLUMN IF NOT EXISTS org_unit_id uuid NULL REFERENCES org_units(id),
  ADD COLUMN IF NOT EXISTS role_id uuid NULL REFERENCES roles(id);

UPDATE memberships SET org_unit_id = team_id WHERE org_unit_id IS NULL;

-- allow associating roles directly to resources (project/report/etc)
CREATE TABLE IF NOT EXISTS role_bindings (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid NOT NULL REFERENCES orgs(id),
  role_id     uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  entity_type text NOT NULL,
  entity_id   uuid NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  UNIQUE (role_id, entity_type, entity_id)
);

CREATE INDEX IF NOT EXISTS idx_role_bindings_entity ON role_bindings(entity_type, entity_id);

-- base permissions for core modules
INSERT INTO permissions (code, description, scope) VALUES
  ('dashboard.view', 'View dashboard', 'org'),
  ('projects.view', 'View projects', 'org'),
  ('projects.manage', 'Create and edit projects', 'org'),
  ('okr.view', 'View OKRs', 'org'),
  ('okr.manage', 'Create and edit OKRs', 'org'),
  ('meetings.view', 'View meetings', 'org'),
  ('meetings.manage', 'Create and edit meetings', 'org'),
  ('reports.view', 'View reports', 'org'),
  ('reports.manage', 'Create/update/submit reports', 'org'),
  ('reports.review', 'Approve or reject reports', 'org'),
  ('workitems.manage', 'Create or modify work items', 'org'),
  ('analytics.view', 'View analytics dashboards', 'org')
ON CONFLICT (code) DO NOTHING;
