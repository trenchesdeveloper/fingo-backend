package utils

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestGenerateHashedPassword(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{"ValidPassword", "mysecretpassword", false},
		{"EmptyPassword", "", false},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := GenerateHashedPassword(tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("GenerateHashedPassword() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError {
				// Verify that the hashed password can be correctly compared to the original password
				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password))
				if err != nil {
					t.Errorf("Generated hash does not match the password: %v", err)
				}
			}
		})
	}
}

func TestCompareHashedPassword(t *testing.T) {
	// Create a password and hash it
	password := "mysecretpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Define test cases
	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectError    bool
	}{
		{"ValidPassword", string(hashedPassword), password, false},
		{"InvalidPassword", string(hashedPassword), "wrongpassword", true},
		{"EmptyPassword", string(hashedPassword), "", true},
		{"EmptyHashedPassword", "", password, true},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompareHashedPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("CompareHashedPassword() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
