# Tasks: OAuth Login for Admin Page

## Backend

1. [x] Add `golang-jwt/jwt/v5` dependency
2. [x] Write tests for JWT creation and validation (sign, verify, expiry, tamper detection)
3. [x] Implement JWT helper functions (create token, validate token, extract claims)
4. [x] Write tests for `OAuthAuthenticator` middleware (valid session, expired, missing cookie)
5. [x] Implement `OAuthAuthenticator` with JWT session validation
6. [x] Write tests for auth handler (login redirect, callback, state mismatch, user not in allowlist)
7. [x] Implement auth handler (login, callback, me, logout)
8. [x] Add `AUTH_MODE` config and authenticator selection in `main.go`
9. [x] Create `admin_users` seed data or allowlist config
10. [x] Register `/api/v1/auth/*` routes

## Frontend

11. [x] Update `AdminLogin` to show "Sign in with GitHub" button
12. [x] Update `AuthProvider` to use `/auth/me` for session checking
13. [x] Remove Basic Auth header logic from API calls
14. [x] Add user avatar and name display in admin header
15. [x] Add logout button calling `/auth/logout`
16. [x] Write tests for OAuth login flow (redirect, session check, logout)

## Integration

17. [ ] Verify OAuth flow end-to-end with a test GitHub OAuth app
18. [x] Verify basic auth still works when `AUTH_MODE=basic`
19. [x] Update `.env.example` with new OAuth env vars
20. [x] Document OAuth setup steps in README
