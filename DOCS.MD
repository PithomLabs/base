# Memos Project Documentation

## 1. Project Overview
**Memos** is a privacy-first, lightweight note-taking service. It allows users to capture thoughts, ideas, and tasks in a social-media-like feed. The project is designed to be self-hosted, with a focus on simplicity and durability. This fork/version (`memos-0244`) includes an extension for a **Jira-like Ticket System**.

## 2. Architecture & Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Echo (Web Framework), gRPC for some internal definitions.
- **Database**: SQLite (default), supports MySQL/PostgreSQL abstractly via `store` interfaces.
- **Entry Point**: `bin/memos/main.go` wires up the server, database driver, and API routes.

### Frontend
- **Framework**: React with Vite.
- **Language**: TypeScript (`.tsx`).
- **Styling**: Tailwind CSS + MUI Joy UI.
- **State Management**: MobX (seen in `memoStore` usage).
- **Communication**: REST API (`/api/v1`), leveraging `long`, `protobuf` definitions for types.

## 3. Directory Structure

- **`bin/`**: Application entry points.
  - `memos/main.go`: The main executable.
- **`server/`**: HTTP server and API logic.
  - `router/api/v1/`: Request handlers for REST endpoints (e.g., `memo_service.go`, `ticket_service.go`).
  - `server.go`: Server lifecycle management.
- **`store/`**: Data persistence layer.
  - **Entity definitions** (structs): `memo.go`, `user.go`, `ticket.go`.
  - **Database drivers**: `db/sqlite`, `db/mysql`, etc.
  - **Migrations**: `migration/sqlite/`.
- **`web/`**: Single Page Application (frontend).
  - `src/pages/`: Route components (e.g., `Tickets.tsx`).
  - `src/components/`: Reusable UI components.
  - `src/store/`: Frontend state stores.
- **`plugin/`**: Integrations (Webhooks, Identity Providers, Storage).
- **`proto/`**: Protobuf definitions for API contracts.

## 4. Core Concepts & Data Model

### Memos
The core entity.
- **Storage**: `memo` table.
- **Attributes**: `Content` (Markdown), `Visibility` (Public/Private/Protected), `Pinned`, `Resources` (attachments).
- **Identity**: Uses `UID` (ShortUUID) for public addressing and integer `ID` for internal relations.
- **Versioning**: Implicit via `UpdatedTs`.

### Users
- **Storage**: `user` table.
- **Roles**: `HOST`, `ADMIN`, `USER`.
- **Auth**: Password hash, Session tokens, or SSO via `idp` table.

### Resources
- **Storage**: `resource` table.
- **Usage**: Images/files attached to memos. Can be local blob or external storage (S3).

### Tickets (New Feature)
A dedicated task management entity, separate from Memos but integrated into the ecosystem.
- **Storage**: `tickets` table (SQLite).
- **Attributes**:
  - `Title`, `Description` (supports linking to Memos via `/m/{uid}`).
  - `Status`: `OPEN`, `IN_PROGRESS`, `CLOSED`.
  - `Priority`: `LOW`, `MEDIUM`, `HIGH`.
  - `Type`: `TASK`, `BUG`, `STORY`.
  - `Tags`: Explicit string array (JSON in DB).
- **Assignee**: Linked to `user` ID.

## 5. Feature Deep Dive: Tagging Systems

The project employs two distinct tagging patterns:

1.  **Memo Tags (Content-Derived)**
    - **Mechanism**: Tags are **not** stored in a separate column. They are parsed dynamically from the Markdown `Content` using `#hashtag` syntax.
    - **Implementation**: The backend `gomark` parser extracts nodes of type `TAG`.
    - **Search**: `FindMemoPayload` allows filtering by `TagSearch`. The backend checks if the parsed content contains the tag.
    - **Management**: `RenameMemoTag` and `DeleteMemoTag` APIs iterate through memos and rewrite their content string.

2.  **Ticket Tags (Explicit)**
    - **Mechanism**: Tickets use a dedicated `tags` column in the database (stored as a JSON string array `[]`).
    - **Implementation**: The `Ticket` struct has a `Tags []string` field.
    - **UI**: Displayed as Chips in the Kanban/List view. Managed via form inputs (requires explicit Add/Remove actions).

## 6. Ticket Service Details
- **API**: `server/router/api/v1/ticket_service.go`
- **Endpoints**:
  - `POST /api/v1/tickets`: Create new ticket.
  - `GET /api/v1/tickets`: List with filters (`type`, `creatorId`).
  - `PATCH /api/v1/tickets/:id`: Update status, priority, assignee, etc.
  - `DELETE /api/v1/tickets/:id`: Remove ticket.
- **Frontend**: `Tickets.tsx` provides a list view and a Kanban board (`TicketKanban.tsx`). It supports drag-and-drop status updates.

## 7. Development Guide

### Prerequisites
- Go 1.22+
- Node.js & npm/pnpm

### Running Locally

**1. Backend**
```bash
# From root directory
go run bin/memos/main.go --mode dev --port 8081
```
*Server will start on localhost:8081.*

**2. Frontend**
```bash
cd web
npm install
npm run dev
```
*Frontend dev server usually starts on localhost:3000, proxying API requests to 8081.*

### Build
- **Frontend**: `npm run build` outputs to `dist/`.
- **Release**: `npm run release` builds and moves assets to `server/router/frontend/dist` for embedding in the Go binary.
- **Backend Build**: `go build -o memos bin/memos/main.go`.
