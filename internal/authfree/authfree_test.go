package authfree

import "testing"

func TestIsAuthFreePath(t *testing.T) {
	tests := []struct {
		name      string
		urlPrefix string
		p         string
		want      bool
	}{
		// Global exact-match whitelist.
		{name: "ping exact", urlPrefix: "/", p: "/ping", want: true},
		{name: "register exact", urlPrefix: "/", p: "/register", want: true},
		{name: "message exact", urlPrefix: "/", p: "/message", want: true},
		{name: "stream exact", urlPrefix: "/", p: "/stream", want: true},
		{name: "version exact", urlPrefix: "/", p: "/version", want: true},
		{name: "info exact", urlPrefix: "/", p: "/info", want: true},
		{name: "info any-case-exact", urlPrefix: "/", p: "/info", want: true},
		// Global subpaths.
		{name: "message subpath", urlPrefix: "/", p: "/message/12", want: true},
		// Lookalike must NOT match (bare-prefix safety).
		{name: "messageevil rejected", urlPrefix: "/", p: "/messageevil", want: false},
		{name: "registerevil rejected", urlPrefix: "/", p: "/registerevil", want: false},
		{name: "infoevil rejected", urlPrefix: "/", p: "/infoevil", want: false},
		{name: "info subpath", urlPrefix: "/", p: "/info/extra", want: true},
		// Unrelated paths.
		{name: "push rejected", urlPrefix: "/", p: "/push", want: false},
		{name: "device push rejected", urlPrefix: "/", p: "/dev/body", want: false},
		{name: "mcp rejected", urlPrefix: "/", p: "/mcp", want: false},

		// Device-scoped 2-segment paths.
		{name: "device message", urlPrefix: "/", p: "/dev/message", want: true},
		{name: "device stream", urlPrefix: "/", p: "/dev/stream", want: true},
		{name: "device version", urlPrefix: "/", p: "/dev/version", want: true},
		{name: "device message id", urlPrefix: "/", p: "/dev/message/12", want: true},
		// The 3-segment message/<id> route treats ANY third segment as the id
		// parameter; it must stay exempt so the gotify token can gate it (the
		// long-term bug this guards: device-scoped delete id needing Basic
		// Auth while the global /message/:id does not).
		{name: "device message id any value", urlPrefix: "/", p: "/dev/message/evil", want: true},
		// Device-scoped lookalikes.
		{name: "device messageevil rejected", urlPrefix: "/", p: "/dev/messageevil", want: false},
		{name: "device message third segment rejected", urlPrefix: "/", p: "/dev/message/id/evil", want: false},
		{name: "device empty id rejected", urlPrefix: "/", p: "/dev/message/", want: false},
		// Device-scoped DELETE all (also exempt).
		{name: "device message delete-all", urlPrefix: "/", p: "/dev/message", want: true},
		// Root and single-segment.
		{name: "root rejected", urlPrefix: "/", p: "/", want: false},
		{name: "single segment rejected", urlPrefix: "/", p: "/dev", want: false},

		// urlPrefix support.
		{name: "prefixed ping", urlPrefix: "/bark", p: "/bark/ping", want: true},
		{name: "prefixed device message", urlPrefix: "/bark", p: "/bark/dev/message", want: true},
		{name: "prefixed lookalike rejected", urlPrefix: "/bark", p: "/bark/dev/messageevil", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthFreePath(tt.urlPrefix, tt.p); got != tt.want {
				t.Errorf("IsAuthFreePath(%q, %q) = %v, want %v", tt.urlPrefix, tt.p, got, tt.want)
			}
		})
	}
}