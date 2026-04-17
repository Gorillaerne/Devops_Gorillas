package handlers

import "testing"

func TestIsBreached_KnownPair(t *testing.T) {
	if !isBreached("Benthe1954", "^Jt^pLkzW2") {
		t.Error("expected known pair to be breached")
	}
}

func TestIsBreached_WrongPassword(t *testing.T) {
	if isBreached("Benthe1954", "wrongpassword") {
		t.Error("expected wrong password to not be breached")
	}
}

func TestIsBreached_UnknownUser(t *testing.T) {
	if isBreached("unknownuser", "^Jt^pLkzW2") {
		t.Error("expected unknown user to not be breached")
	}
}

func TestIsBreached_AllPairs(t *testing.T) {
	for username, password := range breachedCredentials {
		if !isBreached(username, password) {
			t.Errorf("expected %s to be breached with correct password", username)
		}
	}
}
