# Authentication Issues & Fixes

This document covers the OAuth/authentication issues encountered in this project and how they were resolved.

---

## Issue 1: `cannot index tokens (variable of type *services.TokenPair)`

**Error:**
```
internal/handlers/auth.go:450:25: cannot index tokens (variable of type *services.TokenPair)
internal/handlers/auth.go:451:25: cannot index tokens (variable of type *services.TokenPair)
```

**Cause:**  
The `OAuthCallback` handler was treating `tokens` as a `map[string]interface{}` and accessing values via `tokens["access_token"]`. However, `GenerateTokenPair()` returns `*services.TokenPair`, which is a **struct**, not a map. Structs cannot be indexed with bracket notation in Go.

**Fix:**  
Replaced map-style access with struct field access:

```go
// Before (broken)
url.QueryEscape(tokens["access_token"].(string))
url.QueryEscape(tokens["refresh_token"].(string))

// After (fixed)
url.QueryEscape(tokens.AccessToken)
url.QueryEscape(tokens.RefreshToken)
```

**File:** `apps/api/internal/handlers/auth.go`

---

## Issue 2: `no SESSION_SECRET environment variable is set`

**Error:**
```
goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail.
[400] GET /api/auth/oauth/google
```

**Cause:**  
The `gothic` package initializes its session store in an `init()` function, which runs **before** `main()`. The code was calling `os.Setenv("SESSION_SECRET", cfg.JWTSecret)` inside `main()` — too late. By the time the env var was set, gothic had already checked for it during `init()` and initialized with a `nil` store, causing every OAuth request to fail with a 400 error.

**Fix:**  
Instead of relying on the environment variable (which gothic only reads once at init), explicitly assign the session store after loading config:

```go
// Before (broken — env var set too late for gothic's init())
os.Setenv("SESSION_SECRET", cfg.JWTSecret)

// After (fixed — directly sets the store gothic uses)
gothic.Store = sessions.NewCookieStore([]byte(cfg.JWTSecret))
```

This also required adding these imports to `cmd/server/main.go`:
```go
"github.com/gorilla/sessions"
"github.com/markbates/goth/gothic"
```

**File:** `apps/api/cmd/server/main.go`

---

## Issue 3: `redirect_uri_mismatch` (Error 400 from Google)

**Error:**
```
Access blocked: This app's request is invalid
Error 400: redirect_uri_mismatch
```

**Cause:**  
The redirect URI sent to Google didn't match the one registered in Google Cloud Console. Two problems:

1. The `.env` file had `GOOGLE_REDIRECT_URL` set to `http://localhost:8080/api/auth/google/callback` (missing the `/oauth/` path segment).
2. The actual route registered in the app is `/api/auth/oauth/:provider/callback`, so the code constructs the redirect URI as `http://localhost:8080/api/auth/oauth/google/callback`.

Google requires **exact match** between the redirect URI in the request and the one registered in the Cloud Console.

**Fix:**

1. Updated `.env` to use the correct path:
   ```env
   GOOGLE_REDIRECT_URL="http://localhost:8080/api/auth/oauth/google/callback"
   ```

2. **Critically**, update the **Google Cloud Console** to match:
   - Go to [Google Cloud Console → Credentials](https://console.cloud.google.com/apis/credentials)
   - Edit your OAuth 2.0 Client ID
   - Under **Authorized redirect URIs**, set it to:
     ```
     http://localhost:8080/api/auth/oauth/google/callback
     ```
   - Save

**File:** `.env` + Google Cloud Console configuration

---

## Quick Checklist for OAuth Setup

- [ ] `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are set in `.env`
- [ ] `APP_URL` in `.env` matches your server address (e.g., `http://localhost:8080`)
- [ ] `OAUTH_FRONTEND_URL` in `.env` points to your admin app (e.g., `http://localhost:3001`)
- [ ] `JWT_SECRET` is set in `.env` (used for session cookie encryption)
- [ ] Google Cloud Console redirect URI is set to `{APP_URL}/api/auth/oauth/google/callback`
- [ ] For GitHub OAuth: redirect URI is `{APP_URL}/api/auth/oauth/github/callback`
