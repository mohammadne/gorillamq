package auth

import (
	"testing"

	"go.uber.org/zap"
)

func TestEmptyAuth(t *testing.T) {
	auth := NewAuth(zap.NewNop(), &Config{
		Username: " ",
		Password: " ",
	})

	valid := ""
	invalid := "take:me"

	if !auth.Authenticate(valid) {
		t.Error("failed to check empty token")
	}

	if !auth.Authenticate(invalid) {
		t.Error("failed to check free token")
	}
}

func TestAuth(t *testing.T) {
	auth := NewAuth(zap.NewNop(), &Config{
		Username: " ",
		Password: " ",
	})

	valid := "root:password"
	invalid := "take:me"

	if !auth.Authenticate(valid) {
		t.Error("failed to check true token")
	}

	if auth.Authenticate(invalid) {
		t.Error("failed to check free token")
	}
}
