# Breach Detection — breachList.go

**File:** `Go_Refined_Code/handlers/breachList.go`

Implements a simple credential breach check. When a user logs in, their username and password are compared against a list of known-leaked credentials. If they match, the login response includes `"breached": true` and the frontend locks the user to the profile page until they change their password.

---

## How It Works

### The Breach List

`breachedCredentials` is a Go `map[string]string` that maps usernames to their leaked plain-text passwords:

```go
var breachedCredentials = map[string]string{
    "Benthe1954":    "^Jt^pLkzW2",
    "Jack1969":      "_yRw7uqk4h",
    // ... 18 more entries
}
```

There are 20 entries in total. These represent a simulated production database breach used for testing and demonstration.

### `isBreached(username, password string) bool`

Looks up the username in the map and compares the password. Returns `true` only if both match exactly.

```go
func isBreached(username, password string) bool {
    leaked, ok := breachedCredentials[username]
    return ok && leaked == password
}
```

This function is called inside `HandleAPILogin` after the user has successfully authenticated.

---

## Breach Response Flow

```
User logs in with breached credentials
    └── HandleAPILogin verifies password ✓
        └── isBreached() → true
            └── Response: { "breached": true, "token": "..." }
                └── Frontend stores "breachWarning=1" in sessionStorage
                    └── reuseable_functions.js forces redirect to /profile
                        └── User must change password before using the site
```

---

## Email Notifications

If `SEND_BREACH_EMAILS=true` is set, `main.go` calls `SendBreachNotificationsToAll` in a background goroutine on startup. This looks up the email address for each breached username and sends a notification via the Resend API.

See [Email Service](email-service.md) for details.

---

## Limitations

- The breach list is hard-coded in source. In a real system, this would be loaded from a database or an external feed (e.g. Have I Been Pwned API).
- The check only triggers at login time, not on subsequent requests.
- Plain-text passwords are stored in the breach list so they can be compared against user input before bcrypt hashing. This is intentional for the simulation — in production these would typically be hashes checked against a breach database.
