# Task: Integrate Beads Issue Tracking with Ticket System

## Phase 1: Database Schema
- [x] Create migration `02__beads_integration.sql` for SQLite
- [x] Add beads-specific columns to tickets table
- [x] Create agent_workflows table
- [ ] Test migration on dev database

## Phase 2: Backend Core
- [x] Update `store/ticket.go` with new fields
- [x] Update `store/db/sqlite/ticket.go` implementation
- [x] Create `store/agent_workflow.go` interface
- [x] Implement agent workflow store for SQLite
- [x] Add priority/type mapping utilities

## Phase 3: BD CLI Integration
- [x] Create `server/service/beads.go` BD wrapper service
- [x] Implement `bd create` integration in ticket creation
- [/] Implement `bd update` integration in ticket updates
- [/] Implement `bd close` integration in ticket closure
- [x] Add `bd sync` auto-sync functionality

## Phase 4: API Layer
- [x] Update `ticket_service.go` to use BD CLI
- [/] Add dependency management endpoints
- [x] Add workflow logging endpoints
- [/] Add epic/subtask endpoints
- [x] Update route registration

## Phase 5: Frontend Updates
- [ ] Update priority dropdown to P0-P4
- [ ] Add labels multi-select
- [ ] Add dependencies picker
- [ ] Display beads_id prominently
- [ ] Show workflow history

## Phase 6: Testing & Verification
- [ ] Write backend tests
- [ ] Manual UI testing
- [ ] Database verification
- [ ] BD sync verification
- [ ] Integration tests

## Phase 7: Documentation
- [ ] Update README.md
- [ ] Add workflow examples
- [ ] Update Taskfile.yml
