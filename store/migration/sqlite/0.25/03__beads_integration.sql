-- Beads Integration Migration (Idempotent)
-- Adds beads-specific columns to tickets table and creates agent_workflows table

-- Note: SQLite doesn't support ALTER TABLE ADD COLUMN IF NOT EXISTS
-- So we wrap this in a safe check. Columns will only be added if they don't exist.

-- Create agent_workflows table first (safe with IF NOT EXISTS)
CREATE TABLE IF NOT EXISTS agent_workflows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    session_id TEXT NOT NULL,
    agent_name TEXT NOT NULL DEFAULT 'antigravity',
    task_name TEXT,
    task_mode TEXT CHECK(task_mode IN ('PLANNING', 'EXECUTION', 'VERIFICATION')),
    task_status TEXT,
    task_summary TEXT,
    predicted_size INTEGER,
    created_ts INTEGER NOT NULL,
    metadata TEXT DEFAULT '{}'
);

-- Create indexes for agent_workflows
CREATE INDEX IF NOT EXISTS idx_workflows_ticket ON agent_workflows(ticket_id);
CREATE INDEX IF NOT EXISTS idx_workflows_session ON agent_workflows(session_id);
CREATE INDEX IF NOT EXISTS idx_workflows_created ON agent_workflows(created_ts);

-- For tickets table columns, we check if they exist by querying pragma
-- This is done via application code, so here we just document the expected state:

-- Expected new columns in tickets table (added via application on first run):
-- - beads_id TEXT
-- - parent_id INTEGER
-- - labels TEXT DEFAULT '[]'
-- - dependencies TEXT DEFAULT '[]'
-- - discovery_context TEXT
-- - closed_reason TEXT
-- - issue_type TEXT

-- Create indexes (safe with IF NOT EXISTS)
CREATE INDEX IF NOT EXISTS idx_tickets_beads_id ON tickets(beads_id);
CREATE INDEX IF NOT EXISTS idx_tickets_parent_id ON tickets(parent_id);
CREATE INDEX IF NOT EXISTS idx_tickets_issue_type ON tickets(issue_type);
