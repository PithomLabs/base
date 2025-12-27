-- Add foreign key constraints to tickets table via table recreation
-- This migration is needed because foreign keys were not initially defined

PRAGMA foreign_keys = OFF;

-- Store existing data
CREATE TEMPORARY TABLE tickets_backup AS SELECT * FROM tickets;

-- Drop old table
DROP TABLE tickets;

-- Recreate with foreign keys
CREATE TABLE tickets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'OPEN',
  priority TEXT NOT NULL DEFAULT 'MEDIUM',
  type TEXT NOT NULL DEFAULT 'TASK',
  tags TEXT NOT NULL DEFAULT '[]',
  creator_id INTEGER NOT NULL,
  assignee_id INTEGER,
  created_ts BIGINT NOT NULL,
  updated_ts BIGINT NOT NULL,
  beads_id TEXT,
  parent_id INTEGER,
  labels TEXT DEFAULT '[]',
  dependencies TEXT DEFAULT '[]',
  discovery_context TEXT,
  closed_reason TEXT,
  issue_type TEXT,
  FOREIGN KEY (creator_id) REFERENCES user(id) ON DELETE CASCADE,
  FOREIGN KEY (assignee_id) REFERENCES user(id) ON DELETE SET NULL,
  FOREIGN KEY (parent_id) REFERENCES tickets(id) ON DELETE CASCADE
);

-- Restore data
INSERT INTO tickets SELECT * FROM tickets_backup;

-- Drop backup
DROP TABLE tickets_backup;

-- Recreate indexes
CREATE INDEX idx_tickets_creator_id ON tickets (creator_id);
CREATE INDEX idx_tickets_status ON tickets (status);
CREATE INDEX idx_tickets_assignee_id ON tickets (assignee_id);
CREATE UNIQUE INDEX idx_tickets_beads_id ON tickets(beads_id) WHERE beads_id IS NOT NULL;

PRAGMA foreign_keys = ON;
