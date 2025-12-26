# Debugging Notifications

## Issue
User reports that tagging a user in a ticket comment results in "4 spinning icons" in the Notifications page for the tagged user.

## Analysis
- "Spinning icon" typically comes from `MemoCommentMessage.tsx` or `TicketCommentMessage.tsx` when `initialized` state is false.
- This state depends on `activityServiceClient.getActivity` resolving and returning a valid payload.
- Since it was a "comment of ticket", it is likely a `MEMO_COMMENT` type notification, rendered by `MemoCommentMessage`.
- I recently regenerated the frontend protobuf clients (`buf generate`).

## Hypothesis
1. **Frontend Type Mismatch**: The generated `activity_service.ts` might have different property names (e.g., casing) than what `MemoCommentMessage` expects (`memoComment`).
2. **Backend/Frontend Sync**: The backend might return data that the new frontend client fails to decode (e.g. if field IDs shifted/changed, though I only added fields).
3. **Data Missing**: The activity payload might be empty or malformed.

## Plan
1. Check `web/src/types/proto/api/v1/activity_service.ts` to confirm fields are `memoComment` and `ticketComment`.
2. Check `web/src/components/Inbox/MemoCommentMessage.tsx` to see what it expects.
3. Verify `activity.pb.go` again to ensure field tags match.
