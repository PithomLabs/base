# Analysis of Memos v0.24.4 for Ticket Management System

This document outlines the analysis of the Memos codebase (v0.24.4) and the plan to implement a basic Ticket Management System (Jira-like) as requested.

## 1. Project Structure Overview

The project is a standard Go + React application.

-   **Root**: `/home/chaschel/Documents/ibm/go/memos-0244`
-   **Backend**: Written in Go.
    -   `cmd/server`: Entry point.
    -   `server`: HTTP server setup and route definitions (`server/router`).
    -   `store`: Database interactions and data models. Supports SQLite, MySQL, and Postgres.
-   **Frontend**: Written in TypeScript/React.
    -   `web/src`: Source code.
    -   `web/src/pages`: Page components.
    -   `web/src/router`: Client-side routing.
    -   `web/src/components`: Reusable UI components.

## 2. Requirement Analysis

The goal is to add a basic ticket management system with:

-   **Ticket Form**: For adding/updating tickets (Title, Description, Status, etc.).
-   **Storage**: Native support for SQLite, MySQL, and Postgres.
-   **No Sprint Management**: Pure ticket tracking.
-   **Deliverable**: One form/page and a sidebar link.
-   **Constraints**: specific codebase (v0.24.4), native React elements only.

## 3. Implementation Plan

### 3.1 Backend Implementation

We need to add a new entity `Ticket` to the system.

#### 3.1.1 Database Schema (`store`)

We will introduce a `tickets` table. Since Memos supports multiple databases, we need to ensure the schema is compatible or defined for all. The `store` package abstracts the driver.

**Table: `tickets`**

| Column | Type | Description |
| :--- | :--- | :--- |
| `id` | INTEGER | Primary Key, Auto-increment |
| `title` | TEXT | Ticket title |
| `description` | TEXT | Ticket content/description (Markdown supported like memos) |
| `status` | TEXT | `open`, `in_progress`, `closed` (using Enum) |
| `priority` | TEXT | `low`, `medium`, `high` |
| `creator_id` | INTEGER | User ID of the creator |
| `assignee_id` | INTEGER | User ID of the assignee (optional) |
| `created_ts` | BIGINT | Timestamp of creation |
| `updated_ts` | BIGINT | Timestamp of last update |

**New File:** `store/ticket.go`

Define the `Ticket` struct and the `TicketStore` interface.

```go
type Ticket struct {
    ID          int32
    Title       string
    Description string
    Status      string
    Priority    string
    CreatorID   int32
    CreatedTs   int64
    UpdatedTs   int64
}

type FindTicket struct {
    ID *int32
    CreatorID *int32
}

type UpdateTicket struct {
    ID int32
    Title *string
    Status *string
    // ...
}

type DeleteTicket struct {
    ID int32
}
```

#### 3.1.2 Database Migration

Memos uses a migration system in `store/migrator.go`. We should add a strictly monotonic version number and the corresponding SQL for SQLite, MySQL, and Postgres to create the `tickets` table.

#### 3.1.3 API Routes (`server/router`)

We will add REST API endpoints for tickets.

**New File:** `server/router/api/v1/ticket_service.go`

Endpoints:
-   `POST /api/v1/tickets`: Create a ticket.
-   `GET /api/v1/tickets`: List tickets.
-   `PATCH /api/v1/tickets/:id`: Update a ticket.
-   `DELETE /api/v1/tickets/:id`: Delete a ticket.

**Registration:** `server/router/api/v1/v1.go`
Register the `TicketService` to the API group.

### 3.2 Frontend Implementation

We will add a new "Tickets" section to the UI.

#### 3.2.1 Routing (`web/src/router`)

Update `web/src/router/index.tsx` to include the `/tickets` route.

```tsx
// ...
const Tickets = lazy(() => import("@/pages/Tickets"));

export enum Routes {
  // ...
  TICKETS = "/tickets",
}

// ... in routes definition
{
    path: Routes.TICKETS,
    element: (
        <Suspense fallback={<Loading />}>
            <Tickets />
        </Suspense>
    ),
},
```

#### 3.2.2 Sidebar (`web/src/components/HomeSidebar`)

Update functionality in `HomeSidebar` or `Navigation` to show a link to `/tickets`.

#### 3.2.3 Page & Components (`web/src/pages/Tickets.tsx`)

Create a new page `Tickets.tsx`.

**UI Components to Reuse:**
-   **Layout**: Use `HomeLayout` or `RootLayout`.
-   **Buttons/Inputs**: Use components from `web/src/components/kit` (search for `Button`, `Input`).
-   **Dialogs**: Memos uses dialogs heavily (e.g., `CreateMemoDialog`). We can create a `TicketDialog` for creating/editing tickets.
-   **Lists**: Look at `MemoList` for inspiration on how to render a list of items.

**Features:**
1.  **List View**: A table or card list of existing tickets.
2.  **Create/Edit Modal**: A form with:
    -   Title (Input)
    -   Description (Textarea or Markdown Editor if reusable)
    -   Status (Select: Open, In Progress, Closed)
    -   Priority (Select: Low, Medium, High)

## 4. Execution Steps

1.  **Backend Model**: Create `store/ticket.go` and implement the `TicketStore` interface.
2.  **Migrations**: Add the `tickets` table definition.
3.  **Backend API**: Create `server/router/api/v1/ticket_service.go` and register routes.
4.  **Frontend Route**: Add `/tickets` to `web/src/router/index.tsx`.
5.  **Frontend Page**: Create `web/src/pages/Tickets.tsx` with basic list and create button.
6.  **Frontend Sidebar**: Add link to sidebar.
7.  **Frontend Form**: Implement the ticket creation/editing logic.

## 5. Notes

-   **Simplicity**: We will strictly follow the existing patterns. For example, if they use raw SQL in `store`, we do the same. If they use a specific ORM-like structure, we mimic it.
-   **Styles**: We will use the existing Tailwind CSS classes to match the look and feel.
-   **Integration**: The ticket system will be independent of the "Memos" (notes) concept but will share the user authentication and database connection.
