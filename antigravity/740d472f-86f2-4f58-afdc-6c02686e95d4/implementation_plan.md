# Debugging SSE Notifications Plan

The user is not seeing toast popups when tagged. This plan focuses on adding instrumentation to the backend to trace the SSE lifecycle from connection to dispatch.

## Proposed Changes

### [Component: Backend Logging]

#### [MODIFY] [notification_hub.go](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/server/router/api/v1/notification_hub.go)
- Add `fmt.Printf` or log calls in `NotificationStreamHandler` to confirm connection establishment and user ID.
- Add logs in `Register` and `Unregister` to track active connections.
- Add logs in `NotifyUser` to see if it even finds any connections for the target user.

#### [MODIFY] [memo_service.go](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/server/router/api/v1/memo_service.go)
- Add logs in `dispatchMemoMentions` to confirm it's calling `NotifyUser` with the correct `userID`.

### [Component: Frontend Verification]

#### [MODIFY] [index.html](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/web/index.html)
- Switch `data-on-load` to a more explicit Datastar pattern if needed, or add a console log for debugging.

## Verification Plan

### Automated Tests
- None, focusing on manual instrumentation.

### Manual Verification
1. Rebuild the server: `go build -o bin/memos ./bin/memos`.
2. Run the server and check the output logs.
3. Refresh the browser and look for "SSE connection established for user X" in the server logs.
4. Perform a mention and check for "Dispatching SSE notification to user Y" and "Notification sent to Z connections" in the logs.
5. Check browser console for any errors related to Datastar or SSE.
