package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	const interval = time.Second
	tokenString, err := MakeJWT(uuid.New(), "secret", interval)
	if err != nil {
		t.Errorf("MakeJWT function failed: %v", err)
		return
	}

	_, err = ValidateJWT(tokenString, "secret")
	if err != nil {
		t.Errorf("ValidateJWT function failed: %v", err)
		return
	}

	_, err = ValidateJWT(tokenString, "not secret")
	if err == nil {
		t.Errorf("ValidateJWT function failed: accepted invalid secret")
		return
	}

	time.Sleep(interval * 2)
	_, err = ValidateJWT(tokenString, "secret")
	if err == nil {
		t.Errorf("ValidateJWT function failed: accepted expired token")
		return
	}
}
