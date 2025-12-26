# Implementation Plan - Ensure Frontend Embedding

This plan outlines the steps to ensure that the frontend assets are always embedded into the Go binary when compiling `bin/memos/main.go`.

## Proposed Changes

### [Backend]

#### [MODIFY] [frontend.go](file:///home/chaschel/Documents/ibm/go/tix-gemini/server/router/frontend/frontend.go)

- Add a `go:generate` directive to automate the frontend build and sync process.
- This will ensure that running `go generate ./...` (or specifically on this package) will populate the `dist` directory before the Go compiler embeds it.

```go
//go:generate npm --prefix ../../../web install
//go:generate npm --prefix ../../../web run release
```

> [!NOTE]
> The `release` script in `web/package.json` is already configured to output to `../server/router/frontend/dist`.

## Verification Plan

### Automated Tests
1. Run `go generate ./server/router/frontend/frontend.go` to verify the frontend builds and populates the `dist` directory.
2. Run `go build -o memos ./bin/memos/main.go`.
3. Verify the binary contains the frontend assets:
   - Check binary size.
   - Run the binary and access the web interface (if possible in the environment).
   - Alternatively, use `strings memos | grep index.html` or similar to check for embedded assets.

### Manual Verification
- The user can run `npm run release` in the `web` directory followed by `go build ./bin/memos` to verify the single binary deployment.
