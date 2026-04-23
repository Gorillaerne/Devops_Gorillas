package tests

import (
	"devops_gorillas/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendBreachNotification_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("unexpected Authorization header: %s", r.Header.Get("Authorization"))
		}
		var payload handlers.resendPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode payload: %v", err)
		}
		if len(payload.To) == 0 || payload.To[0] != "victim@example.com" {
			t.Errorf("unexpected To: %v", payload.To)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	handlers.resendBaseURL = srv.URL
	t.Cleanup(func() { handlers.resendBaseURL = "https://api.resend.com" })

	err := handlers.SendBreachNotification("test-key", "noreply@example.com", "victim@example.com", "Benthe1954")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendBreachNotification_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	handlers.resendBaseURL = srv.URL
	t.Cleanup(func() { handlers.resendBaseURL = "https://api.resend.com" })

	err := handlers.SendBreachNotification("bad-key", "from@example.com", "to@example.com", "Jack1969")
	if err == nil {
		t.Error("expected error on 401 response")
	}
}

func TestSendBreachNotificationsToAll_SkipsWhenNoAPIKey(t *testing.T) {
	db := newUsersDB(t)
	// Should log and return without error when RESEND_API_KEY is not set
	t.Setenv("RESEND_API_KEY", "")
	handlers.SendBreachNotificationsToAll(db) // must not panic
}

func TestSendBreachNotificationsToAll_SendsToKnownUsers(t *testing.T) {
	notified := map[string]bool{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload handlers.resendPayload
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if len(payload.To) > 0 {
			notified[payload.To[0]] = true
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	handlers.resendBaseURL = srv.URL
	t.Cleanup(func() { handlers.resendBaseURL = "https://api.resend.com" })

	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "noreply@example.com")

	db := newUsersDB(t)
	seedUser(t, db, "Benthe1954", "benthe@example.com", "somepass")
	seedUser(t, db, "Jack1969", "jack@example.com", "somepass")

	handlers.SendBreachNotificationsToAll(db)

	if !notified["benthe@example.com"] {
		t.Error("expected benthe@example.com to be notified")
	}
	if !notified["jack@example.com"] {
		t.Error("expected jack@example.com to be notified")
	}
}
