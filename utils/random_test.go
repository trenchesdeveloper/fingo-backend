package utils

import (
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		length int
	}{
		{"Length10", 10},
		{"Length0", 0},
		{"Length1", 1},
		{"Length50", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("RandomString() length = %d, want %d", len(result), tt.length)
			}
			matched, _ := regexp.MatchString("^[a-z]*$", result)
			if !matched {
				t.Errorf("RandomString() contains non-lowercase letters: %s", result)
			}
		})
	}
}

func TestRandomEmail(t *testing.T) {
	emailRegex := `^[a-z]{6}@example\.com$`
	for i := 0; i < 100; i++ {
		email := RandomEmail()
		matched, err := regexp.MatchString(emailRegex, email)
		if err != nil {
			t.Fatalf("Failed to compile regex: %v", err)
		}
		if !matched {
			t.Errorf("RandomEmail() = %s, does not match regex %s", email, emailRegex)
		}
	}
}
