# Notifications Feature Proposal and Analysis

## Objective
Implement a new "Notifications" page in the sidebar to display notifications when the logged-in user is tagged in:
1.  **Comments** (Existing Memo comments).
2.  **Ticket Description** (New functionality).

## Analysis

### Existing System (`Inbox` & `Activity`)
*   **Memos** already has an `inbox` and `activity` system.
*   **Memo Mentions**: `server/router/api/v1/memo_service.go` (`dispatchMemoMentions`) handles `@username` parsing in memo content and creates `Activity` (Type: `MEMO_COMMENT`) and `Inbox` records.
*   **Tickets**: Currently, `server/router/api/v1/ticket_service.go` does **not** handle mentions in the `Description` field.
*   **Schema**:
    *   `activity` table: Stores activity type and payload (JSON).
    *   `inbox` table: Stores message type and activity ID.
    *   Protobuf definitions (`proto/store/`) define the payloads and enums.

### Proposed Changes

#### 1. Schema & Protobuf Updates
We need to extend the `Activity` and `Inbox` definitions to support Ticket mentions.

*   **`proto/store/activity.proto`**: Add `ActivityTicketCommentPayload` to `ActivityPayload`.
    ```protobuf
    message ActivityTicketCommentPayload {
      int32 ticket_id = 1;
    }
    message ActivityPayload {
      ActivityMemoCommentPayload memo_comment = 1;
      ActivityTicketCommentPayload ticket_comment = 2; // NEW
    }
    ```
*   **`proto/store/inbox.proto`**: Add `TICKET_COMMENT` to `InboxMessage.Type`.
    ```protobuf
    message InboxMessage {
      enum Type {
        // ...
        TICKET_COMMENT = 3; // NEW
      }
      // ...
    }
    ```
*   **`proto/api/v1/inbox_service.proto`**: Update `Inbox.Type` enum to include `TICKET_COMMENT`.
    ```protobuf
    message Inbox {
      enum Type {
        // ...
        TICKET_COMMENT = 3; // NEW
      }
    }
    ```
*   **`store/activity.go`**: Add `ActivityTypeTicketComment` constant.
    ```go
    const (
        ActivityTypeTicketComment ActivityType = "TICKET_COMMENT"
    )
    ```

#### 2. Backend Implementation
*   **`server/router/api/v1/ticket_service.go`**:
    *   Implement `dispatchTicketMentions(ctx context.Context, ticket *store.Ticket) error`.
        *   Logic: Parse `ticket.Description` for `@username` patterns (reuse regex from `memo_service.go`).
        *   For each mentioned user:
            *   Create `Activity` with `Type: TICKET_COMMENT` and `Payload: { TicketID: ticket.ID }`.
            *   Create `Inbox` with `Type: TICKET_COMMENT` and `ActivityID: activity.ID`.
    *   Call `dispatchTicketMentions` in `CreateTicket` and `UpdateTicket` (after successful save).
*   **`server/router/api/v1/inbox_service.go`**:
    *   Ensure `convertInboxFromStore` correctly maps the new types (likely handled automatically if utilizing generated code, but verify Enum mappings).

#### 3. Frontend Implementation
*   **`web/src/components/Navigation.tsx`**:
    *   Add a new "Notifications" link in the sidebar (below "Tickets").
    *   Icon: `BellIcon` or similar.
    *   Path: `/notifications`.
*   **`web/src/pages/Notifications.tsx`** (New Page):
    *   Fetch notifications using `useInboxStore` (or `InboxService`).
    *   Render a list of notifications.
    *   **Memo Comments**: Link to the Memo.
    *   **Ticket Mentions**: Link to `/tickets/:id` (e.g. `Routes.TICKET_DETAIL`).
    *   Include "Mark as Read" / "Archive" functionality (supported by `InboxService`).

#### 4. Data Migration
*   **SQL**: No schema changes required for `activity` or `inbox` tables as they use `TEXT` headers for payloads handling JSON/Protobufs.
*   **Code Generation**: Requires running `task generate` (or `buf generate` / `go generate ./...`) to update Go and TS bindings from modified `.proto` files.

## Execution Plan
1.  **Modify Protobufs**: Update `activity.proto`, `inbox.proto`, `inbox_service.proto`.
2.  **Generate Code**: Run generation tools.
3.  **Update Constants**: Update `store/activity.go`.
4.  **Implement Backend Logic**: Update `ticket_service.go` to dispatch mentions.
5.  **Implement Frontend**: Create `Notifications` page and update Sidebar.
