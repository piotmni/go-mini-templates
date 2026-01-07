# go-auth-device-cli CLI Authentication Specification

## Overview

Implement OAuth 2.0 Device Authorization Grant (RFC 8628) for CLI authentication using Better Auth as the identity provider.

## References

- [Better Auth Device Authorization](https://www.better-auth.com/docs/plugins/device-authorization) - v1.3.x
- [RFC 8628 - OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628)
- [GitHub CLI auth patterns](https://github.com/cli/cli) - reference implementation
- [Stripe CLI auth patterns](https://github.com/stripe/stripe-cli) - reference implementation

## Authentication Flow

```
User runs: go-auth-device auth login
    |
POST /api/auth/device/authorize
    -> {device_code, user_code, verification_uri, interval}
    |
Display: "Visit https://example.com/device and enter code: ABCD-1234"
    |
Open browser automatically (skip with --no-browser)
    |
Poll: POST /api/auth/device/token (every `interval` seconds)
    |
Handle errors: authorization_pending, slow_down, access_denied, expired_token
    |
On success: {access_token, refresh_token, expires_in}
    |
Store tokens in OS keyring, metadata in config file
    |
Verify: GET /api/auth/session -> display user info
```

## Project Structure

```
go-auth-device-cli/
├── cmd/
│   ├── root.go              # cobra root command
│   └── auth/
│       ├── login.go         # auth login command
│       ├── logout.go        # auth logout command
│       └── status.go        # auth status command
├── internal/
│   ├── auth/
│   │   ├── device_flow.go   # RFC 8628 implementation
│   │   ├── token.go         # token types and validation
│   │   └── client.go        # HTTP client for auth endpoints
│   ├── config/
│   │   ├── config.go        # TOML config file management
│   │   └── keyring.go       # secure credential storage
│   └── browser/
│       └── browser.go       # cross-platform browser opener
├── pkg/
│   └── api/
│       └── client.go        # authenticated API client
├── go.mod
└── SPEC.md
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/zalando/go-keyring` | Cross-platform secure credential storage |
| `github.com/pelletier/go-toml/v2` | Config file parsing |
| `github.com/pkg/browser` | Cross-platform browser opening |
| `github.com/briandowns/spinner` | Terminal spinner |
| `github.com/fatih/color` | Colored terminal output |

## Configuration

**Config path:** `~/.config/go-auth-device-cli/config.toml`

```toml
[auth]
hostname = "https://auth.example.com"
user_email = "user@example.com"
expires_at = "2024-01-15T10:30:00Z"
```

**Keyring entries:**
- Service: `go-auth-device-cli`
- Keys: `access_token`, `refresh_token`

## API Endpoints

### Request Device Code

```
POST {hostname}/api/auth/device/authorize
Content-Type: application/json

{
  "client_id": "<CLIENT_ID>",
  "scope": "openid profile email"
}

Response:
{
  "device_code": "abc123",
  "user_code": "ABCD-1234",
  "verification_uri": "https://auth.example.com/device",
  "verification_uri_complete": "https://auth.example.com/device?user_code=ABCD-1234",
  "expires_in": 600,
  "interval": 5
}
```

### Poll for Token

```
POST {hostname}/api/auth/device/token
Content-Type: application/json

{
  "grant_type": "urn:ietf:params:oauth:grant-type:device_code",
  "device_code": "abc123",
  "client_id": "<CLIENT_ID>"
}

Success Response:
{
  "access_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600
}

Error Response:
{
  "error": "authorization_pending",
  "error_description": "..."
}
```

### Error Handling

| Error | Action |
|-------|--------|
| `authorization_pending` | Continue polling |
| `slow_down` | Increase interval by 5 seconds |
| `access_denied` | Abort, user denied |
| `expired_token` | Abort, timeout reached |

## CLI Commands

### `go-wh-repeater auth login`

```
Flags:
  --hostname string    Auth server URL (default from config or prompt)
  --no-browser         Don't open browser automatically
  --timeout duration   Polling timeout (default 5m)
```

### `go-auth-device-cli auth logout`

Clears tokens from keyring and removes auth section from config.

### `go-auth-device-cli auth status`

Displays current authentication state:
- Logged in as: user@example.com
- Server: https://auth.example.com
- Token expires: 2024-01-15 10:30:00

## Worklist

### Phase 1: Project Setup
- [ ] Initialize Go module `go mod init github.com/piotmni/better-auth-device-login/go-auth-device-cli`
- [ ] Add Cobra CLI scaffold with root command
- [ ] Create `auth` command group with `login`, `logout`, `status` subcommands
- [ ] Setup config directory creation (`~/.config/go-wh-repeater/`)

### Phase 2: Core Auth Implementation
- [ ] Implement `DeviceCodeRequest` - POST to device/authorize endpoint
- [ ] Implement `TokenPollLoop` - polling with interval and backoff
- [ ] Implement token response parsing and validation
- [ ] Add error handling for all RFC 8628 error states

### Phase 3: Credential Storage
- [ ] Implement keyring wrapper (`internal/config/keyring.go`)
- [ ] Implement TOML config read/write (`internal/config/config.go`)
- [ ] Add token expiry checking

### Phase 4: UX
- [ ] Add terminal spinner during polling
- [ ] Add colored output for success/error states
- [ ] Implement browser auto-open with fallback message
- [ ] Add `--no-browser` flag support

### Phase 5: Commands Integration
- [ ] Complete `auth login` command with full flow
- [ ] Complete `auth logout` command
- [ ] Complete `auth status` command

### Phase 6: API Client
- [ ] Create authenticated HTTP client wrapper
- [ ] Add automatic token refresh middleware
- [ ] Export client for use by other commands

## Open Questions

| Item | Status |
|------|--------|
| Better Auth server URL | **NEEDS INPUT** |
| OAuth client_id | **NEEDS INPUT** |
| Required scopes | Default: `openid profile email` |

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| CLI framework | Cobra | Industry standard, gh/stripe use it |
| Keyring library | go-keyring | Cross-platform, maintained |
| Config format | TOML | Human-readable, Stripe pattern |
| UI complexity | Minimal (spinner + color) | Fast to implement, sufficient for auth |
