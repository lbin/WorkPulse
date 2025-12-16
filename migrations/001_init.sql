-- WorkPulse initial schema
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS orgs (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        text NOT NULL,
  status      text NOT NULL DEFAULT 'active',
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz NULL
);

CREATE TABLE IF NOT EXISTS teams (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id         uuid NOT NULL REFERENCES orgs(id),
  parent_team_id uuid NULL REFERENCES teams(id),
  name           text NOT NULL,
  path           text NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at     timestamptz NULL,
  UNIQUE (org_id, name)
);

CREATE INDEX IF NOT EXISTS idx_teams_org_id ON teams(org_id);
CREATE INDEX IF NOT EXISTS idx_teams_parent ON teams(parent_team_id);

CREATE TABLE IF NOT EXISTS users (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  password_hash text NOT NULL,
  email        text NOT NULL,
  display_name text NOT NULL,
  avatar_url   text NULL,
  status       text NOT NULL DEFAULT 'active',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  deleted_at   timestamptz NULL,
  UNIQUE (org_id, email)
);

CREATE TABLE IF NOT EXISTS memberships (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  team_id      uuid NOT NULL REFERENCES teams(id),
  user_id      uuid NOT NULL REFERENCES users(id),
  role_in_team text NOT NULL DEFAULT 'member',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  deleted_at   timestamptz NULL,
  UNIQUE (team_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_memberships_user ON memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_memberships_team ON memberships(team_id);

CREATE TABLE IF NOT EXISTS audit_logs (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  actor_user_id uuid NOT NULL REFERENCES users(id),
  action        text NOT NULL,
  entity_type   text NOT NULL,
  entity_id     uuid NOT NULL,
  diff          jsonb NULL,
  created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_org_time ON audit_logs(org_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);

-- OKR
CREATE TABLE IF NOT EXISTS okr_cycles (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid NOT NULL REFERENCES orgs(id),
  team_id     uuid NULL REFERENCES teams(id),
  type        text NOT NULL CHECK (type IN ('month','quarter','halfyear','year')),
  name        text NOT NULL,
  start_date  date NOT NULL,
  end_date    date NOT NULL,
  status      text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','closed')),
  created_by  uuid NOT NULL REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz NULL,
  UNIQUE (org_id, team_id, type, name)
);

CREATE INDEX IF NOT EXISTS idx_cycle_org_team_time ON okr_cycles(org_id, team_id, start_date);

CREATE TABLE IF NOT EXISTS objectives (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id               uuid NOT NULL REFERENCES orgs(id),
  cycle_id             uuid NOT NULL REFERENCES okr_cycles(id),
  team_id              uuid NULL REFERENCES teams(id),
  owner_user_id        uuid NOT NULL REFERENCES users(id),
  parent_objective_id  uuid NULL REFERENCES objectives(id),
  title                text NOT NULL,
  description          text NULL,
  weight               numeric(6,3) NOT NULL DEFAULT 1.0,
  status               text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','closed')),
  tags                 jsonb NOT NULL DEFAULT '[]'::jsonb,
  schema_version       int NOT NULL DEFAULT 1,
  payload              jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz NOT NULL DEFAULT now(),
  deleted_at           timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_obj_cycle ON objectives(cycle_id);
CREATE INDEX IF NOT EXISTS idx_obj_owner ON objectives(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_obj_team ON objectives(team_id);
CREATE INDEX IF NOT EXISTS idx_obj_parent ON objectives(parent_objective_id);
CREATE INDEX IF NOT EXISTS gin_obj_payload ON objectives USING gin (payload jsonb_path_ops);

CREATE TABLE IF NOT EXISTS key_results (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  objective_id  uuid NOT NULL REFERENCES objectives(id),
  title         text NOT NULL,
  metric_type   text NOT NULL,
  target_value  numeric(18,6) NULL,
  current_value numeric(18,6) NULL,
  unit          text NULL,
  weight        numeric(6,3) NOT NULL DEFAULT 1.0,
  confidence    int NOT NULL DEFAULT 70,
  status        text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','closed')),
  schema_version int NOT NULL DEFAULT 1,
  metric_payload jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_kr_objective ON key_results(objective_id);
CREATE INDEX IF NOT EXISTS gin_kr_metric_payload ON key_results USING gin (metric_payload jsonb_path_ops);

CREATE TABLE IF NOT EXISTS okr_alignments (
  id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id    uuid NOT NULL REFERENCES orgs(id),
  from_type text NOT NULL CHECK (from_type IN ('objective','kr')),
  from_id   uuid NOT NULL,
  to_type   text NOT NULL CHECK (to_type IN ('objective','kr')),
  to_id     uuid NOT NULL,
  relation  text NOT NULL DEFAULT 'align',
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (from_type, from_id, to_type, to_id)
);

CREATE INDEX IF NOT EXISTS idx_align_from ON okr_alignments(from_type, from_id);
CREATE INDEX IF NOT EXISTS idx_align_to ON okr_alignments(to_type, to_id);

-- Work Items
CREATE TABLE IF NOT EXISTS work_items (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  team_id       uuid NULL REFERENCES teams(id),
  owner_user_id uuid NOT NULL REFERENCES users(id),
  type          text NOT NULL CHECK (type IN ('REPORT_ENTRY','TASK','MILESTONE','DELIVERABLE','ACTION_ITEM')),
  title         text NOT NULL,
  description   text NULL,
  status        text NOT NULL DEFAULT 'todo' CHECK (status IN ('todo','doing','done','blocked','archived')),
  priority      int NOT NULL DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
  effort_minutes int NULL,
  start_at      timestamptz NULL,
  end_at        timestamptz NULL,
  due_at        timestamptz NULL,
  source_type   text NULL CHECK (source_type IN ('report','project','meeting','manual')),
  source_id     uuid NULL,
  reason_code   text NULL,
  tags          jsonb NOT NULL DEFAULT '[]'::jsonb,
  schema_version int NOT NULL DEFAULT 1,
  custom_fields jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_wi_org_time ON work_items(org_id, COALESCE(end_at, start_at, created_at));
CREATE INDEX IF NOT EXISTS idx_wi_owner_status ON work_items(owner_user_id, status);
CREATE INDEX IF NOT EXISTS idx_wi_team_status ON work_items(team_id, status);
CREATE INDEX IF NOT EXISTS idx_wi_type ON work_items(type);
CREATE INDEX IF NOT EXISTS idx_wi_source ON work_items(source_type, source_id);
CREATE INDEX IF NOT EXISTS gin_wi_custom ON work_items USING gin (custom_fields jsonb_path_ops);

CREATE TABLE IF NOT EXISTS work_item_okr_links (
  id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id              uuid NOT NULL REFERENCES orgs(id),
  work_item_id        uuid NOT NULL REFERENCES work_items(id) ON DELETE CASCADE,
  objective_id        uuid NULL REFERENCES objectives(id),
  key_result_id       uuid NULL REFERENCES key_results(id),
  link_type           text NOT NULL CHECK (link_type IN ('planned','actual','both')),
  contribution_weight numeric(6,3) NOT NULL DEFAULT 1.0,
  evidence_required   boolean NOT NULL DEFAULT false,
  created_at          timestamptz NOT NULL DEFAULT now(),
  CHECK (objective_id IS NOT NULL OR key_result_id IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_wi_okr
ON work_item_okr_links(
  work_item_id,
  COALESCE(objective_id, '00000000-0000-0000-0000-000000000000'::uuid),
  COALESCE(key_result_id,'00000000-0000-0000-0000-000000000000'::uuid),
  link_type
);

CREATE INDEX IF NOT EXISTS idx_wiokr_kr ON work_item_okr_links(key_result_id);
CREATE INDEX IF NOT EXISTS idx_wiokr_obj ON work_item_okr_links(objective_id);

-- Attachments / Comments
CREATE TABLE IF NOT EXISTS attachments (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid NOT NULL REFERENCES orgs(id),
  entity_type text NOT NULL,
  entity_id   uuid NOT NULL,
  file_name   text NOT NULL,
  file_url    text NOT NULL,
  file_type   text NULL,
  size_bytes  bigint NULL,
  uploaded_by uuid NOT NULL REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_attach_entity ON attachments(entity_type, entity_id);

CREATE TABLE IF NOT EXISTS comments (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id         uuid NOT NULL REFERENCES orgs(id),
  entity_type    text NOT NULL,
  entity_id      uuid NOT NULL,
  author_user_id uuid NOT NULL REFERENCES users(id),
  content        text NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comments_entity_time ON comments(entity_type, entity_id, created_at);

-- Reports
CREATE TABLE IF NOT EXISTS reports (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  team_id       uuid NULL REFERENCES teams(id),
  author_user_id uuid NOT NULL REFERENCES users(id),
  type          text NOT NULL CHECK (type IN ('daily','weekly','monthly','quarterly','halfyearly','yearly')),
  period_start  date NOT NULL,
  period_end    date NOT NULL,
  status        text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','submitted','archived')),
  title         text NULL,
  summary       text NULL,
  schema_version int NOT NULL DEFAULT 1,
  payload       jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz NULL,
  UNIQUE (author_user_id, type, period_start, period_end)
);

CREATE INDEX IF NOT EXISTS idx_reports_author_time ON reports(author_user_id, period_start DESC);
CREATE INDEX IF NOT EXISTS idx_reports_team_time ON reports(team_id, period_start DESC);

CREATE TABLE IF NOT EXISTS report_entries (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  report_id     uuid NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
  section       text NOT NULL CHECK (section IN ('done','plan','blocker')),
  content       text NOT NULL,
  effort_minutes int NULL,
  status        text NOT NULL DEFAULT 'done',
  work_item_id  uuid NULL REFERENCES work_items(id),
  order_no      int NOT NULL DEFAULT 0,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_report_entries_report ON report_entries(report_id, section, order_no);

-- Projects
CREATE TABLE IF NOT EXISTS projects (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  team_id       uuid NULL REFERENCES teams(id),
  owner_user_id uuid NOT NULL REFERENCES users(id),
  name          text NOT NULL,
  description   text NULL,
  status        text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','paused','closed')),
  start_date    date NULL,
  end_date      date NULL,
  schema_version int NOT NULL DEFAULT 1,
  payload       jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_proj_team_status ON projects(team_id, status);
CREATE INDEX IF NOT EXISTS idx_proj_owner ON projects(owner_user_id);
CREATE INDEX IF NOT EXISTS gin_proj_payload ON projects USING gin (payload jsonb_path_ops);

CREATE TABLE IF NOT EXISTS project_members (
  project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  user_id    uuid NOT NULL REFERENCES users(id),
  role       text NOT NULL DEFAULT 'member',
  PRIMARY KEY (project_id, user_id)
);

CREATE TABLE IF NOT EXISTS project_work_items (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  project_id   uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  work_item_id uuid NOT NULL REFERENCES work_items(id) ON DELETE CASCADE,
  relation     text NOT NULL DEFAULT 'contains',
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (project_id, work_item_id)
);

CREATE INDEX IF NOT EXISTS idx_projwi_proj ON project_work_items(project_id);

CREATE TABLE IF NOT EXISTS risks_issues (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  project_id    uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  title         text NOT NULL,
  description   text NULL,
  severity      int NOT NULL DEFAULT 3 CHECK (severity BETWEEN 1 AND 5),
  status        text NOT NULL DEFAULT 'open' CHECK (status IN ('open','mitigating','closed')),
  owner_user_id uuid NULL REFERENCES users(id),
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_risk_proj_status ON risks_issues(project_id, status);

-- Meetings
CREATE TABLE IF NOT EXISTS meetings (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  team_id       uuid NULL REFERENCES teams(id),
  host_user_id  uuid NOT NULL REFERENCES users(id),
  title         text NOT NULL,
  agenda        text NULL,
  start_at      timestamptz NOT NULL,
  end_at        timestamptz NOT NULL,
  location      text NULL,
  online_link   text NULL,
  minutes       text NULL,
  status        text NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled','done','archived')),
  schema_version int NOT NULL DEFAULT 1,
  payload       jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_meeting_team_time ON meetings(team_id, start_at DESC);
CREATE INDEX IF NOT EXISTS idx_meeting_host_time ON meetings(host_user_id, start_at DESC);

CREATE TABLE IF NOT EXISTS meeting_attendees (
  meeting_id uuid NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
  user_id    uuid NOT NULL REFERENCES users(id),
  role       text NOT NULL DEFAULT 'attendee',
  PRIMARY KEY (meeting_id, user_id)
);

CREATE TABLE IF NOT EXISTS meeting_action_items (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  meeting_id   uuid NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
  work_item_id uuid NOT NULL REFERENCES work_items(id) ON DELETE CASCADE,
  order_no     int NOT NULL DEFAULT 0,
  created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_meet_ai_meeting ON meeting_action_items(meeting_id, order_no);

-- Config
CREATE TABLE IF NOT EXISTS field_definitions (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id      uuid NOT NULL REFERENCES orgs(id),
  scope_type  text NOT NULL CHECK (scope_type IN ('org','team')),
  scope_id    uuid NULL,
  entity_type text NOT NULL CHECK (entity_type IN ('work_item','report','project','okr')),
  field_key   text NOT NULL,
  field_name  text NOT NULL,
  field_type  text NOT NULL CHECK (field_type IN ('text','number','select','multi_select','date','bool')),
  options     jsonb NULL,
  required    boolean NOT NULL DEFAULT false,
  active      boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  UNIQUE (scope_type, scope_id, entity_type, field_key)
);

CREATE TABLE IF NOT EXISTS report_templates (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id        uuid NOT NULL REFERENCES orgs(id),
  team_id       uuid NULL REFERENCES teams(id),
  type          text NOT NULL CHECK (type IN ('daily','weekly','monthly','quarterly','halfyearly','yearly')),
  name          text NOT NULL,
  schema_version int NOT NULL DEFAULT 1,
  template      jsonb NOT NULL,
  active        boolean NOT NULL DEFAULT true,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);
