package metrics

import "testing"

func TestStatusClass(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{200, "2xx"},
		{204, "2xx"},
		{301, "3xx"},
		{400, "4xx"},
		{429, "4xx"},
		{500, "5xx"},
		{503, "5xx"},
		{99, "other"},
		{600, "5xx"},
	}
	for _, tt := range tests {
		if got := statusClass(tt.code); got != tt.want {
			t.Fatalf("statusClass(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestNew(t *testing.T) {
	reg := New()
	if reg == nil {
		t.Fatal("New() returned nil")
	}
	if reg.Handler() == nil {
		t.Fatal("Handler() returned nil")
	}
	if reg.Middleware() == nil {
		t.Fatal("Middleware() returned nil")
	}
	// Gauge should be settable without panicking.
	reg.SetActiveStreams(3)
}
