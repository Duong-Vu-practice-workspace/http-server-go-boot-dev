package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "test-secret"
	id := uuid.New()
	token, err := MakeJWT(id, secret, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatal(err)
	}
	if got != id {
		t.Fatalf("expected %v got %v", id, got)
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	wrong := "other-secret"
	id := uuid.New()
	token, err := MakeJWT(id, secret, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWT(token, wrong); err == nil {
		t.Fatal("expected error for invalid signature")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	secret := "s"
	id := uuid.New()
	token, err := MakeJWT(id, secret, -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWT(token, secret); err == nil {
		t.Fatal("expected error for expired token")
	}
}
