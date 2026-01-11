package validator

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"invalid - no @", "userexample.com", false},
		{"invalid - no domain", "user@", false},
		{"invalid - no user", "@example.com", false},
		{"invalid - empty", "", false},
		{"invalid - spaces", "user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Email(tt.email); got != tt.want {
				t.Errorf("Email(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"pending", "pending", true},
		{"in-progress", "in-progress", true},
		{"completed", "completed", true},
		{"invalid status", "done", false},
		{"empty", "", false},
		{"uppercase", "PENDING", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Status(tt.status); got != tt.want {
				t.Errorf("Status(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestNonEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"non-empty", "hello", true},
		{"with spaces", "  hello  ", true},
		{"empty", "", false},
		{"only spaces", "   ", false},
		{"only tabs", "\t\t", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NonEmpty(tt.s); got != tt.want {
				t.Errorf("NonEmpty(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
