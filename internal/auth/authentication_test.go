package auth

import (
	"net/http"
	"testing"

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

	token, err := MakeJWT(uuid.New(), secret)
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

	token, err := MakeJWT(userId, secret)
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expected    string
		expectError bool
	}{
		{
			name:        "Valid Bearer Token",
			headers:     http.Header{"Authorization": {"Bearer abc123"}},
			expected:    "abc123",
			expectError: false,
		},
		{
			name:        "Missing Authorization Header",
			headers:     http.Header{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "Malformed Authorization Header",
			headers:     http.Header{"Authorization": {"abc123"}},
			expected:    "",
			expectError: true,
		},
		{
			name:        "Bearer Without Token",
			headers:     http.Header{"Authorization": {"Bearer"}},
			expected:    "",
			expectError: true,
		},
		{
			name:        "Extra Spaces Around Token",
			headers:     http.Header{"Authorization": {"Bearer    abc123   "}},
			expected:    "abc123",
			expectError: false,
		},
		{
			name:        "Bearer Case Insensitivity",
			headers:     http.Header{"Authorization": {"bearer abc123"}},
			expected:    "abc123",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.headers)

			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if token != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, token)
			}
		})
	}
}
