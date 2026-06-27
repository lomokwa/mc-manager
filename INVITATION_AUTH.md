# Invitation-Based User Authentication

Steps to implement invite-link registration so you can send specific people a link to create their password.

---

## 1. Add a Database (SQLite or Postgres)

Add a DB driver to the project:

```bash
# SQLite (simple, no external service)
go get github.com/mattn/go-sqlite3

# OR Postgres
go get github.com/lib/pq
```

## 2. Create Tables

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP
);
```

## 3. Backend Endpoints to Add

### `POST /api/invitations` (Admin only)

Creates an invitation and returns the registration link.

```go
import (
    "crypto/rand"
    "encoding/hex"
    "time"
)

func CreateInvitation(c *gin.Context) {
    // Generate secure random token
    tokenBytes := make([]byte, 32)
    if _, err := rand.Read(tokenBytes); err != nil {
        // handle error
    }
    token := hex.EncodeToString(tokenBytes)

    // Store in DB with expiry (e.g., 72 hours)
    expiresAt := time.Now().Add(72 * time.Hour)
    // INSERT INTO invitations (token, email, expires_at) VALUES (?, ?, ?)

    // Return the link
    frontendURL := os.Getenv("FRONTEND_URL") // e.g., "https://yourdomain.com"
    link := fmt.Sprintf("%s/register?token=%s", frontendURL, token)

    c.JSON(http.StatusOK, gin.H{"link": link, "expires_at": expiresAt})
}
```

### `GET /api/invitations/:token` (Public)

Validates the token is real, unexpired, and unused.

```go
func ValidateInvitation(c *gin.Context) {
    token := c.Param("token")

    // SELECT * FROM invitations WHERE token = ? AND used_at IS NULL AND expires_at > NOW()
    // If not found → 404
    // If found → 200 with email (if set)
}
```

### `POST /api/register` (Public)

Accepts token + username + password, creates the user.

```go
func Register(c *gin.Context) {
    // Parse: token, username, password from body
    // Validate invitation token (same check as above)
    // Hash password with bcrypt
    // INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)
    // UPDATE invitations SET used_at = NOW() WHERE token = ?
    // Return success
}
```

**Dependency:**
```bash
go get golang.org/x/crypto/bcrypt
```

### `POST /api/login` (Public)

Standard login endpoint returning a JWT.

```go
func Login(c *gin.Context) {
    // Parse: username, password from body
    // SELECT password_hash FROM users WHERE username = ?
    // bcrypt.CompareHashAndPassword(hash, password)
    // Generate JWT with user ID + expiry
    // Return token
}
```

**Dependency:**
```bash
go get github.com/golang-jwt/jwt/v5
```

## 4. Replace API Key Middleware with JWT Auth

Update `middleware/auth.go` to validate JWTs instead of (or in addition to) the static API key:

```go
func ValidateJWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        // Strip "Bearer " prefix
        // Parse and validate JWT
        // Set user info in context: c.Set("userID", claims.UserID)
        c.Next()
    }
}
```

## 5. Frontend Changes (React)

### New route: `/register`

- Reads `?token=` from URL query params
- On mount, calls `GET /api/invitations/:token` to validate
- If valid: shows a form (username + password + confirm password)
- On submit: calls `POST /api/register` with token + credentials
- On success: redirect to `/login`

### New route: `/login`

- Username + password form
- Calls `POST /api/login`
- Stores JWT in memory (or httpOnly cookie)
- Redirects to main app

### Auth context

- Wrap app in an auth provider that holds the JWT
- Attach JWT to all API requests via `Authorization: Bearer <token>`
- Redirect unauthenticated users to `/login`

## 6. Flow Summary

```
1. You (admin) call POST /api/invitations → get a link
2. Send link to friend via Discord/email/etc.
3. Friend opens link → /register?token=abc123
4. Frontend validates token, shows password form
5. Friend sets username + password → account created
6. Friend logs in with credentials → gets JWT
7. JWT used for all subsequent API calls
```

## 7. Security Checklist

- [ ] Use `crypto/rand` for token generation (NOT `math/rand`)
- [ ] Tokens expire in 24–72 hours
- [ ] Tokens are single-use (mark `used_at` after registration)
- [ ] Hash passwords with bcrypt (cost ≥ 12)
- [ ] Rate-limit `/register` and `/login` endpoints
- [ ] JWTs expire (e.g., 24 hours) with refresh token rotation
- [ ] Store JWT secret in env var, not in code
- [ ] HTTPS in production
