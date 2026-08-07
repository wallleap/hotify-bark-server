package logging

import "testing"

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want Level
		err  bool
	}{
		{"", LevelInfo, false},
		{"info", LevelInfo, false},
		{"DEBUG", LevelDebug, false},
		{" Warn ", LevelWarn, false},
		{"error", LevelError, false},
		{"bogus", "", true},
		{"123", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseLevel(tt.in)
			if tt.err {
				if err == nil {
					t.Fatalf("ParseLevel(%q) want error, got %v", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseLevel(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ParseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		in   string
		want Format
		err  bool
	}{
		{"", FormatConsole, false},
		{"console", FormatConsole, false},
		{"JSON", FormatJSON, false},
		{" json ", FormatJSON, false},
		{"yaml", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseFormat(tt.in)
			if tt.err {
				if err == nil {
					t.Fatalf("ParseFormat(%q) want error, got %v", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseFormat(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ParseFormat(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
