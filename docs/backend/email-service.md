# Email Service — emailService.go

**File:** `Go_Refined_Code/handlers/emailService.go`

Sends breach notification emails to affected users via the [Resend](https://resend.com) email API. This feature runs once at application startup if `SEND_BREACH_EMAILS=true` is set.

---

## Environment Variables Required

| Variable | Description |
|---|---|
| `RESEND_API_KEY` | API key from the Resend dashboard |
| `RESEND_FROM_EMAIL` | Sender address (must be verified in Resend) |

If either variable is missing, the function logs a warning and exits without sending any emails.

---

## Functions

### `SendBreachNotification(apiKey, fromEmail, toEmail, username string) error`

Sends a single HTML email to one user informing them that their credentials were exposed.

**What it does:**
1. Builds an HTML email body containing the username and a link to the profile page.
2. Creates a `POST /emails` request to `https://api.resend.com`.
3. Sends the request with the API key as a `Bearer` token.
4. Returns an error if the Resend API responds with a 4xx or 5xx status.

The email subject is: `"Important: Your account credentials were exposed"`

The email body instructs the user to:
- Change their password immediately on the profile page.
- Change their password on any other site where they used the same password.

---

### `SendBreachNotificationsToAll(db *sql.DB)`

Iterates over every entry in `breachedCredentials` and sends a notification to each user found in the database.

**Flow:**
1. Reads `RESEND_API_KEY` and `RESEND_FROM_EMAIL` from the environment. Exits early if either is missing.
2. For each username in the breach list, queries the `users` table for their email address.
3. Skips usernames that are not registered in the database (`sql.ErrNoRows`).
4. Calls `SendBreachNotification` for each matched user.
5. Logs success or failure for each notification.

---

## Testing

In tests, the `resendBaseURL` variable is overridden to point to a local test HTTP server. This allows the email sending logic to be tested without a real Resend account or network connection.

---

## When Is This Called?

`main.go` starts this in a background goroutine at startup:

```go
if os.Getenv("SEND_BREACH_EMAILS") == "true" {
    go apiHandlers.SendBreachNotificationsToAll(database.DB)
}
```

The goroutine does not block the server from starting.
