# Invitation-Based User Registration

MC Manager uses an invitation-based registration system. Only users with a valid invitation link can create an account.

## Flow Overview

```
┌─────────┐         ┌──────────┐         ┌──────────┐
│  Admin  │         │  Server  │         │   User   │
└────┬────┘         └────┬─────┘         └────┬─────┘
     │                   │                    │
     │  POST /api/admin/ │                    │
     │  invitations      │                    │
     │──────────────────>│                    │
     │                   │                    │
     │  { link, token,   │                    │
     │    expires_at }   │                    │
     │<──────────────────│                    │
     │                   │                    │
     │  Send link via    │                    │
     │  Discord/email    │                    │
     │───────────────────────────────────────>│
     │                   │                    │
     │                   │  GET /api/         │
     │                   │  invitations/:token│
     │                   │<───────────────────│
     │                   │                    │
     │                   │  { success: true } │
     │                   │───────────────────>│
     │                   │                    │
     │                   │  POST /api/register│
     │                   │  {token, username, │
     │                   │   password}        │
     │                   │<───────────────────│
     │                   │                    │
     │                   │  201 Created       │
     │                   │───────────────────>│
     │                   │                    │
     │                   │  POST /api/login   │
     │                   │  {username,        │
     │                   │   password}        │
     │                   │<───────────────────│
     │                   │                    │
     │                   │  { token: "jwt" }  │
     │                   │───────────────────>│
```

## API Endpoints

### Create Invitation (Admin)

```
POST /api/admin/invitations
Header: X-API-Key: <admin-api-key>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "a1b2c3d4...",
    "link": "http://localhost:5173/register?token=a1b2c3d4...",
    "expires_at": "2024-01-02T15:04:05Z"
  }
}
```

Invitations expire after **24 hours** and are single-use.

### Validate Invitation (Public)

```
GET /api/invitations/:token
```

**Success (200):**
```json
{ "success": true }
```

**Invalid/Expired (404):**
```json
{ "success": false, "error": "invalid or expired invitation" }
```

### Register (Public)

```
POST /api/register
Content-Type: application/json

{
  "token": "a1b2c3d4...",
  "username": "player1",
  "password": "securepassword"
}
```

**Success (201):**
```json
{ "success": true }
```

**Errors (400):**
- Invalid or expired invitation token
- Username already taken
- Missing required fields

### Login (Public)

```
POST /api/login
Content-Type: application/json

{
  "username": "player1",
  "password": "securepassword"
}
```

**Success (200):**
```json
{
  "success": true,
  "data": { "token": "eyJhbGciOi..." }
}
```

**Error (401):**
```json
{ "success": false, "error": "invalid username or password" }
```

## Using the JWT

All protected endpoints require the JWT in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOi...
```

For WebSocket connections (which can't set headers), pass the token as a query parameter:

```
ws://localhost:8080/api/console?token=eyJhbGciOi...
```

### Token Details

- **Algorithm:** HS256
- **Expiry:** 24 hours
- **Claims:** `user_id`, `username`, `exp`

## Frontend Integration

The React client handles auth via:

- **`AuthContext`** — stores the JWT in localStorage, provides `login()`, `register()`, and `logout()`
- **`ProtectedRoute`** — redirects unauthenticated users to `/login`
- **`/login`** — username + password form
- **`/register?token=xxx`** — validates the invitation token, shows registration form

After login, the token is automatically attached to all API requests and WebSocket connections.

## Database Schema

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

## Security

- Invitation tokens are generated with `crypto/rand` (64 hex chars)
- Passwords are hashed with bcrypt (cost 12)
- JWTs are signed with a secret from the `JWT_SECRET` env var
- Login and registration are rate-limited (10 req/sec, burst 20)
- Tokens are single-use — marked with `used_at` after registration
