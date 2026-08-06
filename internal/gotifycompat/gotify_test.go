package gotifycompat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// buildTestService creates a bbolt-backed service in a temp dir.
func buildTestService(t *testing.T, tokenOverride string) *Service {
	t.Helper()
	dir := t.TempDir()
	svc, err := Init(Config{
		DataDir:     dir,
		ClientToken: tokenOverride,
		Version:     "1.2.3",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	return svc
}

func msgOf(title, body string, priority int) Message {
	return Message{
		AppID:    1,
		Title:    title,
		Message:  body,
		Priority: priority,
		Extras:   map[string]interface{}{"device_key": "k1"},
		Date:     time.Now().Format(time.RFC3339Nano),
	}
}

func TestInitGeneratesAndPersistsToken(t *testing.T) {
	svc := buildTestService(t, "")
	if svc.ClientToken() == "" {
		t.Fatal("expected a generated token")
	}
	if !svc.ValidateToken(svc.ClientToken()) {
		t.Fatal("generated token must validate")
	}
	if svc.ValidateToken("wrong-token") {
		t.Fatal("wrong token must not validate")
	}
	if svc.TokenSource() != TokenSourceGenerated {
		t.Fatalf("want TokenSourceGenerated, got %v", svc.TokenSource())
	}
}

// TestInitWithUnusableDataDirFallsBackToMemory guards the "first boot logs
// nothing" symptom: when the data dir cannot be used, Init must still succeed
// on the in-memory store, generate a token and keep the interface usable.
func TestInitWithUnusableDataDirFallsBackToMemory(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatalf("write blocker: %v", err)
	}
	// MkdirAll fails because a regular file sits where a dir would be created.
	svc, err := Init(Config{DataDir: filepath.Join(blocker, "sub"), ClientToken: ""})
	if err != nil {
		t.Fatalf("Init must not fail when the store is unusable: %v", err)
	}
	if svc.ClientToken() == "" {
		t.Fatal("expected a generated token even on memory fallback")
	}
	if !svc.ValidateToken(svc.ClientToken()) {
		t.Fatal("generated token must validate")
	}
	if svc.TokenSource() != TokenSourceGenerated {
		t.Fatalf("want TokenSourceGenerated, got %v", svc.TokenSource())
	}
	if err := svc.Publish("t", "b", 0, nil); err != nil {
		t.Fatalf("Publish on memory store: %v", err)
	}
}

func TestTokenSourceAfterReopen(t *testing.T) {
	dir := t.TempDir()
	svc1, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if svc1.TokenSource() != TokenSourceGenerated {
		t.Fatalf("want Generated on first boot, got %v", svc1.TokenSource())
	}
	if err := svc1.store.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	svc2, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if svc2.TokenSource() != TokenSourcePersisted {
		t.Fatalf("want Persisted after reopen, got %v", svc2.TokenSource())
	}
	if svc2.ClientToken() != "" {
		t.Fatal("persisted token must not be re-logged on restart")
	}
}

func TestTokenSourceOperator(t *testing.T) {
	svc := buildTestService(t, "secret")
	if svc.TokenSource() != TokenSourceOperator {
		t.Fatalf("want TokenSourceOperator, got %v", svc.TokenSource())
	}
}

func TestPublishAndMessagesOrdering(t *testing.T) {
	svc := buildTestService(t, "")
	for i := 0; i < 5; i++ {
		if err := svc.Publish("t", "body", 0, map[string]interface{}{"n": i}); err != nil {
			t.Fatalf("Publish: %v", err)
		}
	}
	msgs, err := svc.Messages(100, 0)
	if err != nil {
		t.Fatalf("Messages: %v", err)
	}
	if len(msgs) != 5 {
		t.Fatalf("want 5 messages, got %d", len(msgs))
	}
	// newest first
	for i := 0; i < len(msgs)-1; i++ {
		if msgs[i].ID <= msgs[i+1].ID {
			t.Fatalf("messages not ordered desc: %d then %d", msgs[i].ID, msgs[i+1].ID)
		}
	}
	// monotonic ids
	for i, m := range msgs {
		if m.ID != uint64(5-i) {
			t.Fatalf("msg[%d].ID = %d, want %d", i, m.ID, 5-i)
		}
	}
}

func TestMessagesSinceFilter(t *testing.T) {
	svc := buildTestService(t, "")
	for i := 0; i < 10; i++ {
		if err := svc.Publish("t", "body", 0, nil); err != nil {
			t.Fatalf("Publish: %v", err)
		}
	}
	// since=7 → only ids 1..6
	msgs, err := svc.Messages(100, 7)
	if err != nil {
		t.Fatalf("Messages: %v", err)
	}
	if len(msgs) != 6 {
		t.Fatalf("want 6 messages with since=7, got %d", len(msgs))
	}
	if msgs[0].ID != 6 || msgs[5].ID != 1 {
		t.Fatalf("unexpected range: got %d..%d", msgs[0].ID, msgs[5].ID)
	}
}

func TestMessagesLimit(t *testing.T) {
	svc := buildTestService(t, "")
	for i := 0; i < 10; i++ {
		_ = svc.Publish("t", "body", 0, nil)
	}
	msgs, _ := svc.Messages(3, 0)
	if len(msgs) != 3 {
		t.Fatalf("want 3 messages, got %d", len(msgs))
	}
	msgs, _ = svc.Messages(0, 0)
	if len(msgs) != 0 {
		t.Fatalf("want 0 messages for limit=0, got %d", len(msgs))
	}
}

func TestPruningRetainsNewest(t *testing.T) {
	svc, err := Init(Config{DataDir: t.TempDir(), ClientToken: "", MaxMessages: 3})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	for i := 0; i < 10; i++ {
		_ = svc.Publish("t", "body", 0, nil)
	}
	msgs, _ := svc.Messages(100, 0)
	if len(msgs) != 3 {
		t.Fatalf("want 3 retained, got %d", len(msgs))
	}
	if msgs[0].ID != 10 || msgs[2].ID != 8 {
		t.Fatalf("unexpected retained ids: %d,%d", msgs[0].ID, msgs[2].ID)
	}
}

func TestPersistenceAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	svc1, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	tok := svc1.ClientToken()
	for i := 0; i < 3; i++ {
		_ = svc1.Publish("t", "body", 0, nil)
	}
	if err := svc1.store.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	svc2, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("reopen Init: %v", err)
	}
	// token stable, not re-generated
	if svc2.ClientToken() != "" {
		t.Fatal("token should not be re-generated/re-logged on restart")
	}
	if !svc2.ValidateToken(tok) {
		t.Fatal("persisted token must still validate after reopen")
	}
	msgs, _ := svc2.Messages(100, 0)
	if len(msgs) != 3 {
		t.Fatalf("want 3 persisted messages, got %d", len(msgs))
	}
	if msgs[0].ID != 3 {
		t.Fatalf("id should continue monotonically, got %d", msgs[0].ID)
	}
	// ids keep increasing after reopen
	_ = svc2.Publish("t", "body", 0, nil)
	msgs, _ = svc2.Messages(1, 0)
	if msgs[0].ID != 4 {
		t.Fatalf("id after reopen should be 4, got %d", msgs[0].ID)
	}
}

func TestTokenOverridePersistsHash(t *testing.T) {
	dir := t.TempDir()
	svc, err := Init(Config{DataDir: dir, ClientToken: "secret"})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if !svc.ValidateToken("secret") {
		t.Fatal("override token must validate")
	}
	if svc.ClientToken() != "" {
		t.Fatal("override token must not be logged")
	}
	// raw plaintext must NOT be in the store
	raw, _ := svc.store.AutoToken()
	if raw != "" {
		t.Fatal("operator token must not be stored in plaintext")
	}

	// reopen: still validates
	svc.store.(*bboltStore).Close()
	svc2, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if !svc2.ValidateToken("secret") {
		t.Fatal("hash must persist across restart")
	}
}

// TestOperatorOverrideClearsStaleAutoToken verifies that switching to an
// operator-supplied token wipes the plaintext of a previously auto-generated
// one from the store.
func TestOperatorOverrideClearsStaleAutoToken(t *testing.T) {
	dir := t.TempDir()
	svc1, err := Init(Config{DataDir: dir, ClientToken: ""})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if raw, _ := svc1.store.AutoToken(); raw == "" {
		t.Fatal("expected persisted plaintext autoToken after first boot")
	}
	if err := svc1.store.(*bboltStore).Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	svc2, err := Init(Config{DataDir: dir, ClientToken: "secret"})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if raw, _ := svc2.store.AutoToken(); raw != "" {
		t.Fatalf("stale autoToken not cleared: %q", raw)
	}
	if !svc2.ValidateToken("secret") {
		t.Fatal("operator token must validate")
	}
}

func TestHubFanout(t *testing.T) {
	svc := buildTestService(t, "")
	ch, unsubscribe := svc.Subscribe()
	defer unsubscribe()

	_ = svc.Publish("t", "body", 1, map[string]interface{}{"a": 1})
	select {
	case m := <-ch:
		if m.Title != "t" || m.Priority != 1 {
			t.Fatalf("unexpected message: %+v", m)
		}
		if _, ok := m.Extras["a"]; !ok {
			t.Fatal("extras lost")
		}
	case <-time.After(time.Second):
		t.Fatal("no broadcast received")
	}

	unsubscribe()
	_ = svc.Publish("t2", "body", 0, nil)
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("should not receive after unsubscribe")
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("expected closed channel after unsubscribe")
	}
}

func TestMessageJSONWireShape(t *testing.T) {
	svc := buildTestService(t, "")
	_ = svc.Publish("title", "body", 2, map[string]interface{}{"device_key": "k", "subtitle": "sub"})
	msgs, _ := svc.Messages(1, 0)
	b, err := json.Marshal(msgs[0])
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, field := range []string{"id", "appid", "title", "message", "priority", "extras", "date"} {
		if _, ok := got[field]; !ok {
			t.Fatalf("missing wire field %q in %s", field, b)
		}
	}
	// date is a string (RFC3339Nano), not a Go time.Time re-marshal
	var date string
	if err := json.Unmarshal(got["date"], &date); err != nil {
		t.Fatalf("date must be a string, got %s", got["date"])
	}
	if _, err := time.Parse(time.RFC3339Nano, date); err != nil {
		t.Fatalf("date not RFC3339Nano: %v", err)
	}
	// extras always an object
	if string(got["extras"]) != "{}" && !reflect.DeepEqual(got["extras"], json.RawMessage(`{"device_key":"k","subtitle":"sub"}`)) {
		t.Fatalf("unexpected extras: %s", got["extras"])
	}
}
