-- Meetings: schedule, action items, links to OKR/Projects/Reports
CREATE TABLE IF NOT EXISTS meetings (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id           uuid NOT NULL REFERENCES orgs(id),
  team_id          uuid NULL REFERENCES teams(id),
  title            text NOT NULL,
  agenda           text NULL,
  scheduled_at     timestamptz NOT NULL,
  duration_minutes int NOT NULL DEFAULT 60,
  facilitator_user_id uuid NULL REFERENCES users(id),
  notes            text NULL,
  attendee_ids     jsonb NOT NULL DEFAULT '[]'::jsonb,
  status           text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','scheduled','completed','cancelled')),
  schema_version   int NOT NULL DEFAULT 1,
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_meetings_org_time ON meetings(org_id, scheduled_at DESC);
CREATE INDEX IF NOT EXISTS idx_meetings_team ON meetings(team_id);

CREATE TABLE IF NOT EXISTS meeting_actions (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  meeting_id   uuid NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
  title        text NOT NULL,
  owner_user_id uuid NULL REFERENCES users(id),
  due_date     date NULL,
  status       text NOT NULL DEFAULT 'todo' CHECK (status IN ('todo','doing','done','archived')),
  related_task_id uuid NULL REFERENCES tasks(id),
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_meeting_actions_meeting ON meeting_actions(meeting_id);
CREATE INDEX IF NOT EXISTS idx_meeting_actions_owner ON meeting_actions(owner_user_id);

CREATE TABLE IF NOT EXISTS meeting_links (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  meeting_id   uuid NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
  target_type  text NOT NULL CHECK (target_type IN ('okr_objective','okr_key_result','project','report')),
  target_id    uuid NOT NULL,
  relation     text NOT NULL DEFAULT 'related',
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (meeting_id, target_type, target_id, relation)
);

CREATE INDEX IF NOT EXISTS idx_meeting_links_meeting ON meeting_links(meeting_id);
CREATE INDEX IF NOT EXISTS idx_meeting_links_target ON meeting_links(target_type, target_id);
