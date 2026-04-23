package tests

import (
	"devops_gorillas/handlers"
	"testing"
)

func TestIsBreached_KnownPair(t *testing.T) {
	if !handlers.isBreached("Benthe1954", "^Jt^pLkzW2") {
		t.Error("expected known pair to be breached")
	}
}

func TestIsBreached_WrongPassword(t *testing.T) {
	if handlers.isBreached("Benthe1954", "wrongpassword") {
		t.Error("expected wrong password to not be breached")
	}
}

func TestIsBreached_UnknownUser(t *testing.T) {
	if handlers.isBreached("unknownuser", "^Jt^pLkzW2") {
		t.Error("expected unknown user to not be breached")
	}
}

func TestIsBreached_AllPairs(t *testing.T) {
	for username, password := range handlers.breachedCredentials {
		if !handlers.isBreached(username, password) {
			t.Errorf("expected %s to be breached with correct password", username)
		}
	}
}
