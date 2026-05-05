// Package handlers emailService
package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// resendBaseURL is the Resend API endpoint. Overridden in tests.
var resendBaseURL = "https://api.resend.com" //nolint:gochecknoglobals

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// SendBreachNotification sends a single breach notification email via Resend.
func SendBreachNotification(apiKey, fromEmail, toEmail, username string) error {
	body := resendPayload{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: "Important: Your account credentials were exposed",
		HTML: fmt.Sprintf(`
<p>Dear %s,</p>
<p>We are writing to inform you that your account credentials were included in a data breach of our production database.</p>
<p>Your username and password were exposed. We strongly recommend that you:</p>
<ul>
  <li><strong>Change your password immediately</strong> at <a href="https://whoknows.example.com/profile">whoknows.example.com/profile</a></li>
  <li>Change your password on any other site where you used the same password</li>
</ul>
<p>We sincerely apologize for this incident.</p>
<p>— The ¿Who Knows? Team</p>
`, username),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resendBaseURL+"/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is controlled internally via resendBaseURL
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API returned status %d", resp.StatusCode)
	}

	return nil
}

// SendBreachNotificationsToAll looks up each breached username in the DB and
// sends a notification email. Requires RESEND_API_KEY and RESEND_FROM_EMAIL env vars.
func SendBreachNotificationsToAll(db *sql.DB) {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if apiKey == "" {
		log.Println("SendBreachNotificationsToAll: RESEND_API_KEY not set, skipping")
		return
	}
	if fromEmail == "" {
		log.Println("SendBreachNotificationsToAll: RESEND_FROM_EMAIL not set, skipping")
		return
	}

	for username := range breachedCredentials {
		var email string
		err := db.QueryRow("SELECT email FROM users WHERE username = ?", username).Scan(&email)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			log.Printf("SendBreachNotificationsToAll: DB error for %s: %v", username, err)
			continue
		}

		if err := SendBreachNotification(apiKey, fromEmail, email, username); err != nil {
			log.Printf("SendBreachNotificationsToAll: failed to notify %s: %v", username, err)
		} else {
			log.Printf("SendBreachNotificationsToAll: notified %s (%s)", username, email)
		}
	}
}
