# Security Critical Fixes - P0 Issues

## Phase 1: Issue Creation
- [x] Propose beads issues for P0 findings
- [x] Create approved beads issues (base-y88, base-nnw, base-16n)
- [ ] Sync beads state

## Phase 2: base-y88 - Foreign Key Constraints
- [/] Enable foreign key constraints in SQLite
- [ ] Test data integrity enforcement
- [ ] Verify no broken references exist
- [ ] Run quality gates

## Phase 3: P0-002 - IDOR in Ticket Operations
- [ ] Add authorization checks to UpdateTicket
- [ ] Add authorization checks to DeleteTicket
- [ ] Add visibility controls to GetTicket
- [ ] Add assignee validation to CreateTicket
- [ ] Write tests for authorization
- [ ] Run quality gates

## Phase 4: P0-003 - SQL Injection Prevention
- [ ] Audit all dynamic SQL query building
- [ ] Implement query parameter validation
- [ ] Add input sanitization for filter fields
- [ ] Test with injection payloads
- [ ] Run quality gates

## Phase 5: Verification
- [ ] Integration testing of all fixes
- [ ] Security regression testing
- [ ] Documentation updates
- [ ] Git commit and push
