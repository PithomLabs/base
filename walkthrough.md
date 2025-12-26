# Ticket Management System Walkthrough

I have successfully implemented the Jira-like Ticket Management System into Memos v0.24.4.

## Key Features Implemented

1.  **Backend (Go)**
    -   **Data Model**: Added `Ticket` struct with fields for `Title`, `Description`, `Status`, `Priority`, `Assignee`, etc.
    -   **Storage**: Implemented `TicketStore` for SQLite, MySQL, and Postgres.
    -   **API**: Created RESTful endpoints at `/api/v1/tickets` for CRUD operations.
    -   **Migration**: Added SQL migration scripts to create the `tickets` table automatically on startup.

2.  **Frontend (React/TypeScript)**
    -   **Route**: Added `/tickets` route.
    -   **Sidebar**: Added a "Tickets" link to the main navigation menu using the `Ticket` icon.
    -   **Page**: Created a comprehensive `Tickets` page with:
        -   List view table showing ID, Title, Status, Priority, Updated time, and Actions.
        -   "New Ticket" dialog with a form to create tickets.
        -   Edit capability by clicking on a ticket title.
        -   Delete capability.

## How to Verify

1.  **Build and Run Backend**:
    ```bash
    go build -o memos ./bin/memos/main.go
    ./memos --mode dev --driver sqlite --data ./data
    ```
    *Note: Ensure you have a data directory or adjust the flags as needed.*

2.  **Build and Run Frontend**:
    ```bash
    cd web
    npm install
    npm run dev
    ```

3.  **Usage**:
    -   Open the app in your browser (usually `http://localhost:3000` or port 8081).
    -   Login.
    -   Click the **Tickets** icon in the sidebar.
    -   Click **New Ticket** to create a ticket.
    -   See it appear in the list.
    -   Click on a ticket title to edit it.

## File Changes

-   `store/ticket.go`: New file defining the model.
-   `store/db/sqlite/ticket.go`, `store/db/mysql/ticket.go`, `store/db/postgres/ticket.go`: Database implementations.
-   `store/migration/*/0.25/00__tickets.sql`: Migrations.
-   `store/driver.go`, `store/store.go`: Interface updates.
-   `server/router/api/v1/ticket_service.go`: New API service.
-   `server/router/api/v1/v1.go`: Route registration.
-   `web/src/router/index.tsx`: Route definition.
-   `web/src/components/Navigation.tsx`: Sidebar link.
-   `web/src/pages/Tickets.tsx`: New page component.
