package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usememos/memos/store"
)

func (d *DB) CreateTicket(ctx context.Context, create *store.Ticket) (*store.Ticket, error) {
	tagsBytes, err := json.Marshal(create.Tags)
	if err != nil {
		return nil, err
	}
	labelsBytes, err := json.Marshal(create.Labels)
	if err != nil {
		return nil, err
	}
	depsBytes, err := json.Marshal(create.Dependencies)
	if err != nil {
		return nil, err
	}

	stmt := `
		INSERT INTO tickets (
			beads_id,
			title,
			description,
			status,
			priority,
			creator_id,
			assignee_id,
			created_ts,
			updated_ts,
			type,
			tags,
			issue_type,
			labels,
			parent_id,
			dependencies,
			discovery_context,
			closed_reason
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	if err := d.db.QueryRowContext(
		ctx,
		stmt,
		create.BeadsID,
		create.Title,
		create.Description,
		create.Status,
		create.Priority,
		create.CreatorID,
		create.AssigneeID,
		create.CreatedTs,
		create.UpdatedTs,
		create.Type,
		string(tagsBytes),
		create.IssueType,
		string(labelsBytes),
		create.ParentID,
		string(depsBytes),
		create.DiscoveryContext,
		create.ClosedReason,
	).Scan(&create.ID); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListTickets(ctx context.Context, find *store.FindTicket) ([]*store.Ticket, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if find.ID != nil {
		where = append(where, "id = ?")
		args = append(args, *find.ID)
	}
	if find.CreatorID != nil {
		where = append(where, "creator_id = ?")
		args = append(args, *find.CreatorID)
	}
	if find.Type != nil {
		where = append(where, "type = ?")
		args = append(args, *find.Type)
	}
	if find.Description != nil {
		where = append(where, "description = ?")
		args = append(args, *find.Description)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			beads_id,
			title,
			description,
			status,
			priority,
			creator_id,
			assignee_id,
			created_ts,
			updated_ts,
			type,
			tags,
			issue_type,
			labels,
			parent_id,
			dependencies,
			discovery_context,
			closed_reason
		FROM tickets
		WHERE %s
		ORDER BY created_ts DESC
	`, strings.Join(where, " AND "))

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.Ticket, 0)
	for rows.Next() {
		var ticket store.Ticket
		var tagsStr, labelsStr, depsStr string
		if err := rows.Scan(
			&ticket.ID,
			&ticket.BeadsID,
			&ticket.Title,
			&ticket.Description,
			&ticket.Status,
			&ticket.Priority,
			&ticket.CreatorID,
			&ticket.AssigneeID,
			&ticket.CreatedTs,
			&ticket.UpdatedTs,
			&ticket.Type,
			&tagsStr,
			&ticket.IssueType,
			&labelsStr,
			&ticket.ParentID,
			&depsStr,
			&ticket.DiscoveryContext,
			&ticket.ClosedReason,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsStr), &ticket.Tags); err != nil {
			ticket.Tags = []string{}
		}
		if err := json.Unmarshal([]byte(labelsStr), &ticket.Labels); err != nil {
			ticket.Labels = []string{}
		}
		if err := json.Unmarshal([]byte(depsStr), &ticket.Dependencies); err != nil {
			ticket.Dependencies = []int32{}
		}
		list = append(list, &ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) GetTicket(ctx context.Context, find *store.FindTicket) (*store.Ticket, error) {
	list, err := d.ListTickets(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (d *DB) UpdateTicket(ctx context.Context, update *store.UpdateTicket) (*store.Ticket, error) {
	set, args := []string{}, []interface{}{}
	if update.Title != nil {
		set = append(set, "title = ?")
		args = append(args, *update.Title)
	}
	if update.Description != nil {
		set = append(set, "description = ?")
		args = append(args, *update.Description)
	}
	if update.Status != nil {
		set = append(set, "status = ?")
		args = append(args, *update.Status)
	}
	if update.Priority != nil {
		set = append(set, "priority = ?")
		args = append(args, *update.Priority)
	}
	if update.AssigneeID != nil {
		set = append(set, "assignee_id = ?")
		args = append(args, *update.AssigneeID)
	}
	if update.UpdatedTs != nil {
		set = append(set, "updated_ts = ?")
		args = append(args, *update.UpdatedTs)
	}

	if update.Type != nil {
		set = append(set, "type = ?")
		args = append(args, *update.Type)
	}
	if update.Tags != nil {
		tagsBytes, err := json.Marshal(update.Tags)
		if err != nil {
			return nil, err
		}
		set = append(set, "tags = ?")
		args = append(args, string(tagsBytes))
	}
	if update.IssueType != nil {
		set = append(set, "issue_type = ?")
		args = append(args, *update.IssueType)
	}
	if update.Labels != nil {
		labelsBytes, err := json.Marshal(update.Labels)
		if err != nil {
			return nil, err
		}
		set = append(set, "labels = ?")
		args = append(args, string(labelsBytes))
	}
	if update.Dependencies != nil {
		depsBytes, err := json.Marshal(update.Dependencies)
		if err != nil {
			return nil, err
		}
		set = append(set, "dependencies = ?")
		args = append(args, string(depsBytes))
	}
	if update.ClosedReason != nil {
		set = append(set, "closed_reason = ?")
		args = append(args, *update.ClosedReason)
	}

	args = append(args, update.ID)
	stmt := fmt.Sprintf(`
		UPDATE tickets
		SET %s
		WHERE id = ?
		RETURNING id, beads_id, title, description, status, priority, creator_id, assignee_id, created_ts, updated_ts, type, tags, issue_type, labels, parent_id, dependencies, discovery_context, closed_reason
	`, strings.Join(set, ", "))

	var ticket store.Ticket
	var tagsStr, labelsStr, depsStr string
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&ticket.ID,
		&ticket.BeadsID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Status,
		&ticket.Priority,
		&ticket.CreatorID,
		&ticket.AssigneeID,
		&ticket.CreatedTs,
		&ticket.UpdatedTs,
		&ticket.Type,
		&tagsStr,
		&ticket.IssueType,
		&labelsStr,
		&ticket.ParentID,
		&depsStr,
		&ticket.DiscoveryContext,
		&ticket.ClosedReason,
	); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(tagsStr), &ticket.Tags); err != nil {
		ticket.Tags = []string{}
	}
	if err := json.Unmarshal([]byte(labelsStr), &ticket.Labels); err != nil {
		ticket.Labels = []string{}
	}
	if err := json.Unmarshal([]byte(depsStr), &ticket.Dependencies); err != nil {
		ticket.Dependencies = []int32{}
	}

	return &ticket, nil
}

func (d *DB) DeleteTicket(ctx context.Context, delete *store.DeleteTicket) error {
	stmt := `DELETE FROM tickets WHERE id = ?`
	result, err := d.db.ExecContext(ctx, stmt, delete.ID)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
