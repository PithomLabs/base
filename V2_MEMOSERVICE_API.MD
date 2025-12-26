# V2 MemoService API Analysis

## Context

User "ading" receives PERMISSION_DENIED errors when viewing notifications that reference memos they don't have access to (private memos created by other users).

## Current Behavior Analysis

### GetMemo Permission Logic (server/router/api/v1/memo_service.go:210-242)

```go
func (s *APIV1Service) GetMemo(ctx context.Context, request *v1pb.GetMemoRequest) (*v1pb.Memo, error) {
    // ... get memo from store ...
    
    if memo.Visibility != store.Public {
        user, err := s.GetCurrentUser(ctx)
        if user == nil {
            return nil, status.Errorf(codes.PermissionDenied, "permission denied")
        }
        if memo.Visibility == store.Private && memo.CreatorID != user.ID {
            return nil, status.Errorf(codes.PermissionDenied, "permission denied")  // ← THIS IS THE ERROR
        }
    }
    // ...
}
```

### Permission Rules
| Visibility | Who Can View |
|------------|--------------|
| PUBLIC     | Anyone |
| PROTECTED  | Any authenticated user |
| PRIVATE    | Only the creator |

### Root Cause
The notification system creates inbox entries for users who are mentioned in comments, but the **underlying memo might be PRIVATE**. When the mentioned user tries to view the notification, they can't access the memo content.

## "V2 MemoService API" - What Did User Mean?

There is **NO v2 API** in the codebase. The user likely meant:
1. Using a different/alternative approach
2. Creating a specialized API for notifications
3. Modifying the permission model

## Proposed Solutions

### Option A: Skip Permission Check for Notification Context (Recommended)

**Concept**: Add a special "notification preview" endpoint or modify activity conversion to include memo snippet directly.

**Pros**:
- No security compromise (only show snippet, not full content)
- Clean separation of concerns

**Cons**:
- Requires backend changes

**Implementation**:
1. Store memo snippet in Activity payload during creation (already has `relatedMemoId`)
2. Add `memo_snippet` field to `ActivityMemoCommentPayload` in proto
3. Backend stores snippet when creating activity
4. Frontend shows snippet from activity without fetching memo

### Option B: Internal Memo Fetch for Activity Conversion

**Concept**: When converting activity to API response, fetch memo internally (bypassing user permission) but only return minimal info.

**Pros**: 
- Minimal changes
- Activity already contains memo names

**Cons**:
- Slightly complex permission model

**Implementation**:
1. Modify `activity_service.go::convertActivityPayloadFromStore` to use internal store directly
2. Only include memo UID and snippet (already done in current code)
3. Frontend shows activity without needing to fetch full memo

### Option C: Change Frontend to Not Fetch Memo Details

**Concept**: The notification message doesn't NEED the full memo. Just show:
"[User] mentioned you in a memo"

**Pros**:
- Fastest to implement
- No backend changes needed

**Cons**:
- Less informative (no memo preview)
- Already partially implemented but still shows permission error

**Implementation**:
1. Modify `MemoCommentMessage.tsx` to not call `memoStore.getOrFetchMemoByName()`
2. Use only the information from the activity payload (which we already have)
3. Show generic message without memo content

## Recommendation

**I recommend Option C** as the immediate fix:

1. The activity payload already contains `memo` and `relatedMemo` names
2. The sender information is already available
3. We don't need to show the memo content in the notification
4. This eliminates all permission errors

**Future Enhancement**: Option A for richer notifications (showing snippet)

## Proposed Changes

### File: `web/src/components/Inbox/MemoCommentMessage.tsx`

**Current Flow**:
```typescript
const activity = await activityServiceClient.getActivity({...});
if (activity.payload?.memoComment) {
  const memo = await memoStore.getOrFetchMemoByName(...);  // ← PERMISSION ERROR HERE
  setRelatedMemo(memo);
}
```

**New Flow**:
```typescript
const activity = await activityServiceClient.getActivity({...});
if (activity.payload?.memoComment) {
  // Don't fetch full memo - just use the name from activity
  const memoName = activity.payload.memoComment.relatedMemo;
  // Display notification using only activity data
}
```

**UI Change**:
- Before: "{user} commented on {memo.content snippet}"
- After: "{user} mentioned you in a memo" (clickable to view)

When clicked, navigate to the memo page - if user has access they'll see it, if not they'll get proper error page.

## Verification Plan

### Manual Testing
1. Login as "ibm2100", create a memo comment mentioning @ading
2. Login as "ading"
3. Navigate to /notifications
4. **Expected**: See notification message without spinner/error
5. Click notification
6. **Expected**: Navigate to memo (may show permission error page if private - this is acceptable)

---

**Awaiting your approval to proceed with Option C implementation.**
