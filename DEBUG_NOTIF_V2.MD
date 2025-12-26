# Debugging Notification V2

## Issue
The `notifications` table is empty after testing ticket comments with mentions. Logic seems to not be triggering or failing silently.

## Investigation Steps
1.  **Code Review**: Check call sites of `dispatchTicketMentions` in `ticket_service.go`.
2.  **Logic Verification**: Verify mention parsing and user lookup in `dispatchTicketMentions`.
3.  **Reproduction**: Attempt to reproduce via API/Server test.

## Findings
- **Bug Identified**: `dispatchTicketMentions` and `dispatchMemoMentions` only look up users by `Nickname`.
- **Root Cause**: The `FindUser` implementation generates an exact match query on the `nickname` column. If a user is mentioned by their `username` (which is standard behavior for `@handle`), the lookup fails if their `nickname` column is empty or different.
- **Fix**: Update the lookup logic to try `Username` first, and if not found, try `Nickname` (or check both).

## Verification
- Code has been patched to query by `Username` first, then `Nickname`.
- Backend code compiles successfully.
