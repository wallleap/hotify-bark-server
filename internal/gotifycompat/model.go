package gotifycompat

// Message mirrors the gotify message wire format consumed by hotify-bridge.
// Date is intentionally an opaque string (RFC3339Nano) so it passes through
// untouched into the bridge's Push Kit payloads (the bridge treats it as an
// opaque string and never parses it).
type Message struct {
	ID       uint64                 `json:"id"`
	AppID    uint                   `json:"appid"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Priority int                    `json:"priority"`
	Extras   map[string]interface{} `json:"extras"`
	Date     string                 `json:"date"`
}
