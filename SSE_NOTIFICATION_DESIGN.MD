# SSE Notification Design with Datastar

## Executive Summary

This document proposes a **real-time notification system** using **Server-Sent Events (SSE)** with the [datastar-go](https://github.com/starfederation/datastar-go) SDK. Instead of the current poll-and-fetch model, notifications will be **pushed instantly** from server to client when a user is mentioned.

---

## Configuration Decisions

| Option | Choice |
|--------|--------|
| Frontend Approach | Pure Datastar (minimal React changes) |
| Sound Effects | Disabled |
| Toast Position | Top-right |

---

## How SSE Works With Existing Functionality

The SSE system **enhances and works alongside** the current notification system - it does NOT replace everything:

| Scenario | Current Behavior | With SSE Enhancement |
|----------|------------------|----------------------|
| Page Load | Fetch all notifications via API | **Same** - initial load via existing API |
| New Notification | Need to refresh page to see | **Real-time push** - appears instantly |
| Badge Count | Static until page refresh | **Live update** via SSE signals |
| Toast Popup | None | **New** - top-right toast when mentioned |

### Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                     Notifications Page                            │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │ Initial Load: Existing API (/api/v1/inboxes)                │ │
│  │ Real-time Updates: SSE stream appends new notifications     │ │
│  └─────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                     Any Page (Navigation Bar)                     │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │ SSE updates notification badge in real-time                 │ │
│  │ Toast popup appears top-right when mentioned                │ │
│  └─────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

### Key Points
- ✅ Existing notifications page still works as before
- ✅ SSE **adds** real-time updates without replacing existing API calls
- ✅ Toast popup (top-right) provides immediate feedback
- ✅ Badge updates in real-time across all pages
- ✅ Graceful degradation if SSE connection fails (falls back to existing behavior)

## Current Problems

| Problem | Impact |
|---------|--------|
| Notifications require page refresh | Poor UX |
| Fetching memo details causes PERMISSION_DENIED | Broken functionality |
| Multiple API calls to resolve notification data | Slow, complex |
| No real-time feedback when mentioned | Users miss notifications |

---

## Proposed Solution: Datastar SSE Architecture

### What is Datastar?

Datastar is a lightweight hypermedia framework that:
1. **Backend → Frontend**: Push DOM patches via SSE (no polling)
2. **Frontend → Backend**: Send signals/data via HTTP requests
3. **Zero npm dependencies**: Single `<script>` tag inclusion

### Key Datastar Features We'll Use

```go
sse := datastar.NewSSE(w, r)

// 1. Patch DOM elements directly
sse.PatchElements(`<div id="notification-badge">5</div>`)

// 2. Update client-side signals (state)
sse.MarshalAndPatchSignals(map[string]any{
    "unreadCount": 5,
    "notifications": []Notification{...},
})

// 3. Execute JavaScript
sse.ExecuteScript(`playNotificationSound()`)
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                           BACKEND (Go)                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  User creates comment with @mention                                  │
│           │                                                          │
│           ▼                                                          │
│  ┌─────────────────────┐    ┌─────────────────────────────────────┐ │
│  │ MemoService.Create  │───▶│ NotificationHub.Dispatch(userID)    │ │
│  └─────────────────────┘    └─────────────────────────────────────┘ │
│                                        │                             │
│                                        ▼                             │
│                              ┌─────────────────────┐                 │
│                              │  SSE Connection Map │                 │
│                              │  userID → []SSE     │                 │
│                              └─────────────────────┘                 │
│                                        │                             │
│                                        ▼                             │
│                              sse.PatchElements(...)                  │
│                              sse.MarshalAndPatchSignals(...)         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ SSE Stream
                                        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         FRONTEND (React + Datastar)                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │ <div data-on-load="$$get('/api/v1/notifications/stream')">      ││
│  │                                                                  ││
│  │   <!-- Dynamically patched by server -->                        ││
│  │   <div id="notification-list">                                  ││
│  │     <NotificationItem sender="ibm2100" message="..." />         ││
│  │   </div>                                                        ││
│  │                                                                  ││
│  │   <!-- Badge updated in real-time -->                           ││
│  │   <span id="notification-badge" data-text="$unreadCount"></span>││
│  │                                                                  ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Detailed Design

### 1. Backend: NotificationHub

**File**: `server/router/api/v1/notification_hub.go`

```go
package v1

import (
    "sync"
    "github.com/starfederation/datastar-go/datastar"
)

type NotificationHub struct {
    mu          sync.RWMutex
    connections map[int32][]*datastar.SSE // userID -> active SSE connections
}

var hub = &NotificationHub{
    connections: make(map[int32][]*datastar.SSE),
}

// Register adds a new SSE connection for a user
func (h *NotificationHub) Register(userID int32, sse *datastar.SSE) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.connections[userID] = append(h.connections[userID], sse)
}

// Unregister removes an SSE connection
func (h *NotificationHub) Unregister(userID int32, sse *datastar.SSE) {
    h.mu.Lock()
    defer h.mu.Unlock()
    // Remove specific SSE from slice...
}

// NotifyUser sends a notification to all connections for a user
func (h *NotificationHub) NotifyUser(userID int32, notification Notification) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    for _, sse := range h.connections[userID] {
        // Patch the notification list with new item
        sse.PatchElements(renderNotificationHTML(notification))
        
        // Update unread count signal
        sse.MarshalAndPatchSignals(map[string]any{
            "unreadCount": getUnreadCount(userID),
        })
        
        // Optional: Play sound
        sse.ExecuteScript(`playNotificationSound()`)
    }
}
```

### 2. Backend: SSE Endpoint

**File**: `server/router/api/v1/notification_service.go`

```go
// NotificationStream handles SSE connections for real-time notifications
func (s *APIV1Service) NotificationStream(w http.ResponseWriter, r *http.Request) {
    user, err := s.GetCurrentUser(r.Context())
    if err != nil || user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Create SSE writer
    sse := datastar.NewSSE(w, r)
    
    // Register connection
    hub.Register(user.ID, sse)
    defer hub.Unregister(user.ID, sse)
    
    // Send initial state
    notifications := s.getRecentNotifications(user.ID)
    sse.MarshalAndPatchSignals(map[string]any{
        "notifications": notifications,
        "unreadCount":   len(filterUnread(notifications)),
    })
    
    // Keep connection alive until client disconnects
    <-r.Context().Done()
}
```

### 3. Backend: Dispatch on Mention

**File**: `server/router/api/v1/memo_service.go` (modify `dispatchMemoMentions`)

```go
func (s *APIV1Service) dispatchMemoMentions(ctx context.Context, memo *store.Memo) error {
    // ... existing mention parsing logic ...
    
    for _, mentionedUser := range mentionedUsers {
        // Existing: Create activity and inbox records
        s.createActivityAndInbox(mentionedUser, memo)
        
        // NEW: Push real-time notification via SSE
        hub.NotifyUser(mentionedUser.ID, Notification{
            Type:       "MEMO_COMMENT",
            SenderName: currentUser.Nickname,
            MemoName:   memo.Name,
            Snippet:    getSnippet(memo.Content),
            Timestamp:  time.Now(),
        })
    }
    return nil
}
```

### 4. Frontend: Datastar Integration

**Option A: Pure Datastar (No React changes)**

Add Datastar script and use `data-*` attributes:

```html
<!-- In index.html -->
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0/bundles/datastar.js"></script>

<!-- Notification component -->
<div data-on-load="$$get('/api/v1/notifications/stream')">
  <div id="notification-list">
    <!-- Server patches this with new notifications -->
  </div>
  <span id="notification-badge" data-text="$unreadCount"></span>
</div>
```

**Option B: React + Datastar Hybrid**

Create a React component that initializes Datastar SSE:

```tsx
// NotificationStream.tsx
import { useEffect, useRef } from "react";

const NotificationStream = () => {
  const containerRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    // Datastar will auto-connect when element with data-on-load is in DOM
    // Or manually use EventSource for more control
    const eventSource = new EventSource("/api/v1/notifications/stream");
    
    eventSource.onmessage = (event) => {
      // Handle incoming SSE events
      // Datastar SDK handles this automatically via data-* attributes
    };
    
    return () => eventSource.close();
  }, []);
  
  return (
    <div 
      ref={containerRef}
      data-on-load="$$get('/api/v1/notifications/stream')"
    >
      <div id="notification-list"></div>
    </div>
  );
};
```

### 5. Notification HTML Rendering (Server-Side)

```go
func renderNotificationHTML(n Notification) string {
    return fmt.Sprintf(`
        <div id="notification-%d" class="notification-item unread">
            <div class="notification-icon">💬</div>
            <div class="notification-content">
                <strong>%s</strong> mentioned you in a memo
                <p class="snippet">%s</p>
                <span class="timestamp">%s</span>
            </div>
            <button data-on-click="$$post('/api/v1/notifications/%d/archive')">
                Mark Read
            </button>
        </div>
    `, n.ID, n.SenderName, n.Snippet, n.Timestamp.Format("3:04 PM"), n.ID)
}
```

---

## Benefits Over Current System

| Feature | Current | Proposed SSE |
|---------|---------|--------------|
| Real-time | ❌ Poll on page load | ✅ Instant push |
| Permission Issues | ❌ Client fetches memo | ✅ Server sends snippet |
| Complexity | ❌ Multiple API calls | ✅ Single SSE stream |
| Notification Sound | ❌ Not possible | ✅ `ExecuteScript()` |
| Badge Updates | ❌ Manual refresh | ✅ Auto-update via signals |
| Offline Support | ❌ None | ✅ SSE auto-reconnects |

---

## Implementation Plan

### Phase 1: Backend SSE Infrastructure
1. Add `github.com/starfederation/datastar-go` dependency
2. Create `NotificationHub` for managing SSE connections
3. Add `/api/v1/notifications/stream` SSE endpoint
4. Modify `dispatchMemoMentions` to push via SSE

### Phase 2: Frontend Integration
1. Add Datastar script to `index.html`
2. Create `NotificationStream` component
3. Update notification badge to use SSE signals
4. Remove old `MemoCommentMessage` fetch logic

### Phase 3: Polish
1. Add notification sounds
2. Add toast popups for new notifications
3. Implement "Mark all as read" via SSE
4. Add connection retry logic

---

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `server/router/api/v1/notification_hub.go` | CREATE | SSE connection manager |
| `server/router/api/v1/notification_service.go` | CREATE | SSE endpoint handler |
| `server/router/api/v1/memo_service.go` | MODIFY | Push notifications on mention |
| `server/router/api/v1/v1.go` | MODIFY | Register SSE route |
| `web/index.html` | MODIFY | Add Datastar script |
| `web/src/components/NotificationStream.tsx` | CREATE | SSE client component |
| `web/src/components/Inbox/MemoCommentMessage.tsx` | DELETE | No longer needed |

---

## Verification Plan

### Automated Testing
- Unit tests for `NotificationHub` (register, unregister, notify)
- Integration test for SSE endpoint (mock HTTP connection)

### Manual Testing
1. Open two browser tabs logged in as different users
2. User A mentions @UserB in a memo comment
3. **Expected**: User B sees notification appear instantly (no refresh)
4. Click "Mark Read" - notification updates in real-time
5. Close and reopen tab - SSE reconnects automatically

---

**Awaiting your approval to proceed with implementation.**
