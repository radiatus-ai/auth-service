package jwt

import (
	"testing"
	"time"
)

const (
	testSecret = "test-secret-key"
	testUserID = "test-user-id"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Generate a token
	token, err := GenerateToken(testUserID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check if the user ID in the claims matches the original
	if claims.UserID != testUserID {
		t.Errorf("UserID mismatch. Expected %s, got %s", testUserID, claims.UserID)
	}
}

func TestExpiredToken(t *testing.T) {
	// Generate a token that expires immediately
	token, err := GenerateToken(testUserID, testSecret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate the expired token
	_, err = ValidateToken(token, testSecret)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}

func TestInvalidToken(t *testing.T) {
	// Try to validate an invalid token
	_, err := ValidateToken("invalid-token", testSecret)
	if err == nil {
		t.Error("Expected an error for invalid token, got nil")
	}
}

func TestGetUserIDFromToken(t *testing.T) {
	// Generate a token
	token, err := GenerateToken(testUserID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Get user ID from the token
	userID, err := GetUserIDFromToken(token, testSecret)
	if err != nil {
		t.Fatalf("Failed to get user ID from token: %v", err)
	}

	// Check if the retrieved user ID matches the original
	if userID != testUserID {
		t.Errorf("UserID mismatch. Expected %s, got %s", testUserID, userID)
	}
}

func TestGetUserIDFromInvalidToken(t *testing.T) {
	// Try to get user ID from an invalid token
	_, err := GetUserIDFromToken("invalid-token", testSecret)
	if err == nil {
		t.Error("Expected an error for invalid token, got nil")
	}
}
