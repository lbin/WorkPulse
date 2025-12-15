-- Reports enhancements: content payload, status workflow, links, review metadata
ALTER TABLE reports
        ADD COLUMN IF NOT EXISTS content jsonb NOT NULL DEFAULT '{}'::jsonb,
        ADD COLUMN IF NOT EXISTS submitted_at timestamptz NULL,
        ADD COLUMN IF NOT EXISTS reviewed_at timestamptz NULL,
        ADD COLUMN IF NOT EXISTS reviewer_user_id uuid NULL REFERENCES users(id),
        ADD COLUMN IF NOT EXISTS review_comment text NULL;

-- Expand status values to cover submit/review flow
ALTER TABLE reports DROP CONSTRAINT IF EXISTS reports_status_check;
ALTER TABLE reports
        ADD CONSTRAINT reports_status_check
        CHECK (status IN ('draft','submitted','approved','rejected','archived'));

-- Links from reports to external entities
CREATE TABLE IF NOT EXISTS report_links (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id       uuid NOT NULL REFERENCES orgs(id),
  report_id    uuid NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
  target_type  text NOT NULL CHECK (target_type IN ('okr_objective','okr_key_result','project','meeting')),
  target_id    uuid NOT NULL,
  relation     text NOT NULL DEFAULT 'related',
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (report_id, target_type, target_id, relation)
);

CREATE INDEX IF NOT EXISTS idx_report_links_report ON report_links(report_id);
CREATE INDEX IF NOT EXISTS idx_report_links_target ON report_links(target_type, target_id);
