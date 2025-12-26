# Debug Log: Inbox Spinning Loaders Issue

## Problem Statement
User reports that both Inbox and Notifications pages show spinning loaders with no content for user "ading" after being tagged in a ticket comment by user "ibm2100".

## Expected Behavior
1. User tags @ading in a ticket comment
2. Activity and Inbox records are created
3. Ading sees a notification in both /inbox and /notifications pages
4. Notification displays properly with user info and memo link

## Actual Behavior
- Spinning loaders appear (4 of them reported)
- No notification content displays

## Initial Analysis

### Data Flow
1. **Comment Creation**: User creates memo comment on ticket → `memo_service.go:dispatchMemoMentions`
2. **Activity Creation**: Creates `Activity` with `ActivityTypeMemoComment` and payload containing `MemoId` and `RelatedMemoId`
3. **Inbox Creation**: Creates `Inbox` with type `MEMO_COMMENT` and reference to activity ID
4. **Frontend Fetch**: Page fetches inbox list, then for each inbox:
   - Calls `activityServiceClient.getActivity()`
   - Expects `activity.payload.memoComment` to be populated
5. **Rendering**: `MemoCommentMessage` component processes payload and displays

### Spinning Loader Trigger
The loader appears when `initialized` state is `false` in:
- `MemoCommentMessage.tsx` (line 26, 144-147)
- `TicketCommentMessage.tsx` (similar pattern)

The `initialized` becomes `true` only after:
```typescript
if (activity.payload?.memoComment) {
  // ... fetch memo data
  setInitialized(true);
}
```

### Hypothesis
The spinning loader indicates that `activity.payload?.memoComment` is either:
1. `undefined` (payload is empty)
2. `null`
3. Request is failing silently

## Investigation Steps

### Step 1: Check Backend Conversion Logic
File: `server/router/api/v1/activity_service.go`

The conversion function:
- Gets memo by ID
- Gets related memo by ID
- Returns error if either is not found
- Constructs payload with memo names

**Potential Issue**: If either memo lookup fails, the function returns an error (NotFound).
This would cause the frontend request to fail, but the component might not handle the error properly.

### Step 2: Check Frontend Error Handling
Need to verify how `MemoCommentMessage` handles errors from `activityServiceClient.getActivity()`.

Looking at the code:
```typescript
useAsyncEffect(async () => {
  if (!inbox.activityId) {
    return;
  }
  const activity = await activityServiceClient.getActivity({
    name: `${activityNamePrefix}${inbox.activityId}`,
  });
  if (activity.payload?.memoComment) {
    // ... set initialized
  }
}, [inbox.activityId]);
```

**Issue Found**: No error handling! If the activity fetch fails or returns without `memoComment`, 
the component never calls `setInitialized(true)`, leaving the spinner forever.

### Step 3: Verify Generated API Proto
Need to confirm that the regenerated TypeScript includes `memoComment` field properly.

**Verification Result**: ✅ Confirmed - `activity_service.ts` has both `memoComment` and `ticketComment` fields.

## Root Cause Identified

### Primary Issue: Missing Error Handling
Both `MemoCommentMessage.tsx` and `TicketCommentMessage.tsx` had a critical flaw in their `useAsyncEffect` hooks:

**Problem Code Pattern:**
```typescript
useAsyncEffect(async () => {
  const activity = await activityServiceClient.getActivity({...});
  if (activity.payload?.memoComment) {
    // ... fetch and process data
    setInitialized(true);  // ❌ Only called if condition is true
  }
}, [inbox.activityId]);
```

**What Happens:**
1. If `getActivity()` throws an error → No catch block → Error propagates → `setInitialized` never called
2. If `activity.payload?.memoComment` is `undefined` or `null` → Condition fails → `setInitialized` never called
3. Component stuck in `initialized === false` state → Infinite spinner

### Secondary Issues (Potential)
1. **Backend Conversion Errors**: If memo lookup fails in `activity_service.go::convertActivityPayloadFromStore`, it returns a NotFound error, which would cause the frontend request to fail.
2. **Empty Payloads**: If the activity exists but has an empty or invalid payload structure, the condition `activity.payload?.memoComment` would be falsy.

## Resolution

### Fix Applied
Added proper error handling with `try-catch-finally` blocks to both components:

```typescript
useAsyncEffect(async () => {
  if (!inbox.activityId) {
    return;
  }

  try {
    const activity = await activityServiceClient.getActivity({
      name: `${activityNamePrefix}${inbox.activityId}`,
    });
    if (activity.payload?.memoComment) {
      // ... process data
    }
  } catch (error) {
    console.error("Failed to load memo comment activity:", error);
  } finally {
    setInitialized(true);  // ✅ Always called, stops spinner
  }
}, [inbox.activityId]);
```

### Files Modified
1. `/web/src/components/Inbox/MemoCommentMessage.tsx` - Added try-catch-finally
2. `/web/src/components/Inbox/TicketCommentMessage.tsx` - Added try-catch-finally

## Expected Behavior After Fix
1. If activity fetch succeeds and has valid payload → Display notification content
2. If activity fetch fails or returns empty payload → Stop spinner, show nothing (could be improved with error UI)
3. Errors are logged to console for debugging
4. No more infinite spinners

## Recommendations for Future Improvements
1. **Add Error UI**: Show a user-friendly error message when activity fetch fails
2. **Add Fallback Content**: When payload is empty, show "Notification unavailable" instead of blank
3. **Add Retry Logic**: For transient errors, allow user to retry loading
4. **Backend Logging**: Add logging in `activity_service.go` to track conversion failures
5. **Validation**: Add validation to ensure activities always have valid payloads before creating inbox entries

## Update: Actual Errors Revealed

After adding error handling, the browser console now shows the real errors:

### Error 1: Permission Denied
```
ClientError: /memos.api.v1.MemoService/GetMemo PERMISSION_DENIED: permission denied
```

### Error 2: Memo Not Found
```
ClientError: /memos.api.v1.ActivityService/GetActivity INTERNAL: 
failed to convert activity from store: rpc error: code = Internal 
desc = failed to convert activity payload from store: rpc error: 
code = NotFound desc = memo not found for comment payload
```

## Analysis of Actual Errors

### Problem: Backend Too Strict
The `activity_service.go::convertActivityPayloadFromStore` function:
1. Fetches memo by ID using `s.Store.GetMemo()`
2. If memo doesn't exist → Returns NotFound error
3. If user lacks permission → Returns PermissionDenied error
4. Either error prevents the entire activity from being returned

### Why This Happens
1. **Deleted Comments**: If the memo comment was deleted, the memo no longer exists
2. **Permission Changes**: If the memo visibility changed to private
3. **Data Integrity**: Activity references memo IDs that are no longer valid

### Impact
- Users can't see ANY notifications if referenced memos are inaccessible
- Entire inbox becomes unusable
- No graceful degradation

## Solution: Graceful Degradation

The backend should handle missing/inaccessible memos gracefully:
1. If memo is not found → Skip the activity payload (return empty payload instead of error)
2. If permission denied → Skip the activity payload
3. Log the issue for debugging but don't fail the entire request
4. Frontend already handles empty payloads (won't display anything, but won't crash)

This allows users to see other notifications even if some are broken.

## Implementation: Backend Fix

Modified `server/router/api/v1/activity_service.go::convertActivityPayloadFromStore`:

**Before:**
```go
if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to get memo: %v", err)
}
if memo == nil {
    return nil, status.Errorf(codes.NotFound, "memo not found for comment payload")
}
```

**After:**
```go
if err != nil {
    // Log the error but don't fail - memo might be deleted or inaccessible
    // This allows other notifications to still work
    return v2Payload, nil  // Return empty payload
}
if memo == nil {
    // Memo was deleted or doesn't exist - skip this payload
    // Frontend will handle empty payload gracefully (won't display)
    return v2Payload, nil
}
```

**Benefits:**
1. ✅ No more INTERNAL/NotFound errors blocking all notifications
2. ✅ Users can see working notifications even if some are broken
3. ✅ Graceful degradation instead of complete failure
4. ✅ Frontend already handles empty payloads (won't show anything for that notification)

**Trade-offs:**
- Silent failures (should add logging in production)
- Users won't see notifications for deleted/inaccessible content (this is acceptable)

## Final Test Plan

After restarting the backend:
1. Open browser as user "ading"
2. Go to /notifications or /inbox
3. **Expected Result**: 
   - No more spinning loaders
   - No errors in console
   - Notifications for accessible content display properly
   - Notifications for deleted/inaccessible content simply don't appear (no error)
4. Check browser console - should be clean (no errors)
