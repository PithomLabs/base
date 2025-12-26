package teststore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/memos/store"
)

func TestTicketStore(t *testing.T) {
	ctx := context.Background()
	ts := NewTestingStore(ctx, t)
	user, err := createTestingHostUser(ctx, ts)
	require.NoError(t, err)

	// Create Ticket
	ticketCreate := &store.Ticket{
		Title:       "Test Ticket",
		Description: "/m/valid-memo-id",
		Status:      store.TicketStatusOpen,
		Priority:    store.TicketPriorityHigh,
		Type:        "BUG",
		Tags:        []string{"backend", "urgent"},
		CreatorID:   user.ID,
		CreatedTs:   1600000000,
		UpdatedTs:   1600000000,
	}

	// Test invalid description
	invalidTicket := *ticketCreate
	invalidTicket.Description = "Invalid Description"
	err = invalidTicket.Validate()
	require.Error(t, err)
	ticket, err := ts.CreateTicket(ctx, ticketCreate)
	require.NoError(t, err)
	require.NotNil(t, ticket)
	require.Equal(t, "BUG", ticket.Type)
	require.Equal(t, []string{"backend", "urgent"}, ticket.Tags)

	// List Ticket (Filter by Type)
	typeBug := "BUG"
	list, err := ts.ListTickets(ctx, &store.FindTicket{
		Type: &typeBug,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(list))
	require.Equal(t, ticket.ID, list[0].ID)

	// Update Ticket
	newType := "STORY"
	newTags := []string{"frontend"}
	updated, err := ts.UpdateTicket(ctx, &store.UpdateTicket{
		ID:   ticket.ID,
		Type: &newType,
		Tags: newTags,
	})
	require.NoError(t, err)
	require.Equal(t, "STORY", updated.Type)
	require.Equal(t, []string{"frontend"}, updated.Tags)

	// Verify Update Persisted
	fetched, err := ts.GetTicket(ctx, &store.FindTicket{ID: &ticket.ID})
	require.NoError(t, err)
	require.Equal(t, "STORY", fetched.Type)

	// Delete
	err = ts.DeleteTicket(ctx, &store.DeleteTicket{ID: ticket.ID})
	require.NoError(t, err)
	list, err = ts.ListTickets(ctx, &store.FindTicket{CreatorID: &user.ID})
	require.NoError(t, err)
	require.Equal(t, 0, len(list))

	ts.Close()
}
