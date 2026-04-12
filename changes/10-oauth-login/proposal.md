# Proposal: OAuth Login for Admin Page

## Why

The current admin panel uses HTTP Basic Auth with credentials passed on every request. This has several limitations:

- Credentials are stored in the frontend's memory and sent with every API call
- No session management — no token expiry, refresh, or revocation
- Adding additional admin users requires sharing a single set of credentials
- No audit trail of which user performed which action
- Basic Auth is not suitable for production-facing admin panels

Replacing Basic Auth with OAuth 2.0 (starting with GitHub as a provider) gives us proper session management, multi-user support, and a familiar login experience.

## What

Replace the existing `BasicAuthenticator` with an `OAuthAuthenticator` that implements the same `Authenticator` interface. GitHub OAuth is the primary provider, with the architecture supporting additional providers (Google, etc.) in the future.

## What Changes

### Backend
- Add `OAuthAuthenticator` implementing the existing `Authenticator` interface
- Implement OAuth 2.0 Authorization Code flow with GitHub as provider
- Add session management with signed JWT tokens stored in HTTP-only cookies
- Create `/api/v1/auth/login` (redirect to GitHub), `/api/v1/auth/callback` (handle OAuth callback), `/api/v1/auth/logout` (clear session), `/api/v1/auth/me` (current user info)
- Add an `admin_users` table to whitelist allowed GitHub usernames/IDs
- Replace Basic Auth middleware on admin routes with JWT session middleware
- Keep `BasicAuthenticator` as a fallback for local development (controlled by env var)

### Frontend
- Replace the username/password login form with a "Sign in with GitHub" button
- Store session state from HTTP-only cookie (no credential storage in JS)
- Add user avatar and name display in admin header
- Update auth context to use session-based auth instead of Basic Auth headers
- Add logout button that clears the session

### Configuration
- `OAUTH_PROVIDER=github` (default, extensible to `google`, etc.)
- `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` for OAuth app credentials
- `OAUTH_ALLOWED_USERS=username1,username2` for access control
- `SESSION_SECRET` for JWT signing
- `AUTH_MODE=oauth|basic` to select authenticator (basic for local dev)

## Approach

The existing `Authenticator` interface was designed for this swap. The `OAuthAuthenticator` implements the same `Middleware()` method, so handler code does not change. The selection between Basic and OAuth is a config-time decision in `main.go`.

Sessions use signed JWTs in HTTP-only, Secure, SameSite=Lax cookies — no tokens in localStorage. The JWT contains the GitHub user ID and username. Token expiry is set to 24 hours with refresh on activity.

Access control is whitelist-based: only GitHub users listed in `OAUTH_ALLOWED_USERS` can access admin routes, even if they successfully authenticate with GitHub.
