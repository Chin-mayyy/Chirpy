package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}
	if len(hashedPassword) == 0 {
		t.Errorf("Hashed password is empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf("Error checking password hash: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, "wrongpassword")
	if err == nil {
		t.Errorf("Expected error for wrong password")
	}
}

func TestMakeJWT(t *testing.T) {
	secret := "shakalakaboomboom"
	expiresIn := time.Minute * 3

	token, err := MakeJWT(uuid.New(), secret, expiresIn)
	if err != nil {
		t.Errorf("Expected no error got: %v", err)
		return
	}

	if token == "" {
		t.Errorf("Expected a token, got an empty string")
		return
	}
}

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	secret := "shakalakboomboom"
	expiresIn := time.Minute * 3

	token, err := MakeJWT(userId, secret, expiresIn)
	if err != nil {
		t.Errorf("expected no error got: %v", err)
		return
	}

	parsedId, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("expected no error got: %v", err)
	}

	if parsedId != userId {
		t.Errorf("Expected userID %v got: %v", parsedId, userId)
	}

	// Test with an invalid token
	_, err = ValidateJWT("invalidTOkenstring", secret)
	if err == nil {
		t.Errorf("Expected error for invalid token, got nil")
	}
}
