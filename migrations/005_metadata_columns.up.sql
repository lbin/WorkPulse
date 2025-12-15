-- Normalize schema versioning and metadata support across entities
ALTER TABLE orgs ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE users ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE memberships ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE okr_cycles ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE okr_objectives ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE okr_key_results ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE okr_links ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE okr_alignments ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE work_items ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE work_item_okr_links ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE attachments ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE comments ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE reports ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE report_entries ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE report_links ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE projects ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE project_members ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE project_work_items ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE risks_issues ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE meetings ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE meeting_attendees ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE meeting_action_items ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE meeting_actions ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE meeting_links ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE field_definitions ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE report_templates ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE org_units ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE roles ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE role_bindings ADD COLUMN IF NOT EXISTS schema_version int NOT NULL DEFAULT 1, ADD COLUMN IF NOT EXISTS metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
