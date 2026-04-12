# Spec: OAuth Login for Admin Page

## Purpose

Replace HTTP Basic Auth with OAuth 2.0 (GitHub provider) for the admin panel, providing proper session management, multi-user support, and production-grade authentication.

## Requirements

### Requirement: OAuthAuthenticator SHALL implement the existing Authenticator interface

#### Scenario: OAuthAuthenticator satisfies the interface

- **GIVEN:** The `Authenticator` interface defines `Middleware()` and `CurrentUser()` methods
- **WHEN:** `OAuthAuthenticator` is instantiated with OAuth config
- **THEN:** It satisfies the `Authenticator` interface
- **AND:** It can be swapped in place of `BasicAuthenticator` without changing handlers

### Requirement: GET /api/v1/auth/login SHALL redirect to GitHub OAuth authorization

#### Scenario: User clicks "Sign in with GitHub"

- **GIVEN:** `GITHUB_CLIENT_ID` is configured
- **WHEN:** `GET /api/v1/auth/login` is requested
- **THEN:** Response is a 302 redirect to `https://github.com/login/oauth/authorize`
- **AND:** The redirect includes `client_id`, `redirect_uri`, `scope=read:user`, and a random `state` parameter
- **AND:** The `state` parameter is stored in a short-lived cookie for CSRF protection

### Requirement: GET /api/v1/auth/callback SHALL exchange the authorization code for a session

#### Scenario: GitHub redirects back with a valid code

- **GIVEN:** User authorized the app on GitHub
- **WHEN:** `GET /api/v1/auth/callback?code=abc&state=xyz` is requested
- **THEN:** The backend exchanges the code for an access token with GitHub
- **AND:** Fetches the user's GitHub profile (username, avatar)
- **AND:** Checks the username against `OAUTH_ALLOWED_USERS`
- **AND:** Creates a signed JWT session token
- **AND:** Sets the JWT in an HTTP-only, Secure, SameSite=Lax cookie
- **AND:** Redirects to `/admin`

#### Scenario: State parameter mismatch (CSRF protection)

- **GIVEN:** The `state` parameter does not match the stored cookie
- **WHEN:** `GET /api/v1/auth/callback?code=abc&state=invalid` is requested
- **THEN:** Response status is 403 with error "invalid state parameter"

#### Scenario: User not in allowlist

- **GIVEN:** GitHub user "stranger" is not in `OAUTH_ALLOWED_USERS`
- **WHEN:** The callback completes with valid credentials for "stranger"
- **THEN:** Response status is 403 with error "user not authorized"

### Requirement: GET /api/v1/auth/me SHALL return current user information

#### Scenario: Authenticated user

- **GIVEN:** Valid JWT session cookie is present
- **WHEN:** `GET /api/v1/auth/me` is requested
- **THEN:** Response is 200 with `{"username": "...", "avatar_url": "...", "provider": "github"}`

#### Scenario: No session

- **GIVEN:** No session cookie is present
- **WHEN:** `GET /api/v1/auth/me` is requested
- **THEN:** Response status is 401

### Requirement: POST /api/v1/auth/logout SHALL clear the session

#### Scenario: User logs out

- **GIVEN:** Valid session exists
- **WHEN:** `POST /api/v1/auth/logout` is requested
- **THEN:** The session cookie is cleared (MaxAge=-1)
- **AND:** Response status is 200

### Requirement: OAuth session middleware SHALL protect admin routes

#### Scenario: Valid JWT grants access

- **GIVEN:** Request has a valid, non-expired JWT in the session cookie
- **WHEN:** A request to `/api/v1/admin/posts` is made
- **THEN:** The request proceeds to the handler
- **AND:** User info is available via `CurrentUser(r)`

#### Scenario: Expired JWT returns 401

- **GIVEN:** JWT has expired (older than 24 hours without refresh)
- **WHEN:** A request to `/api/v1/admin/posts` is made
- **THEN:** Response status is 401

### Requirement: JWT tokens SHALL be signed and validated securely

#### Scenario: Token signing

- **GIVEN:** `SESSION_SECRET` is configured
- **WHEN:** A JWT is created after OAuth callback
- **THEN:** The token is signed with HMAC-SHA256 using `SESSION_SECRET`
- **AND:** The token contains `sub` (GitHub user ID), `username`, `exp` (24h), `iat` claims

#### Scenario: Tampered token is rejected

- **GIVEN:** A JWT with a modified payload
- **WHEN:** The middleware validates the token
- **THEN:** Validation fails and response status is 401

### Requirement: AUTH_MODE env var SHALL select the authenticator

#### Scenario: OAuth mode (production)

- **GIVEN:** `AUTH_MODE=oauth`
- **WHEN:** The server starts
- **THEN:** `OAuthAuthenticator` is used for admin routes

#### Scenario: Basic mode (local development)

- **GIVEN:** `AUTH_MODE=basic` (or unset)
- **WHEN:** The server starts
- **THEN:** `BasicAuthenticator` is used (existing behavior preserved)
