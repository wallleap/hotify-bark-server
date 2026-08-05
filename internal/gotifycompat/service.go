package gotifycompat

import (
	"fmt"
	"time"

	"github.com/mritd/logger"
)

// Config controls the gotify-compatible monitoring interface.
type Config struct {
	DataDir     string // where gotify.db lives; empty/unusable falls back to memory
	ClientToken string // optional pre-seeded client token; empty → auto-generated & persisted
	Version     string // value served on GET /version; empty → "dev"
	MaxMessages int    // retained message count; <=0 → defaultMaxMessages
}

const (
	defaultMaxMessages = 1000
	defaultAppID       = 1
)

// Service ties together token auth, message persistence and the WebSocket hub.
type Service struct {
	store     Store
	hub       *Hub
	tokenHash []byte
	autoToken string
	version   string
}

// Init creates the service. It never fails outright for storage reasons:
// persistence is best-effort and degrades to an in-memory store when the data
// directory cannot be used.
func Init(cfg Config) (*Service, error) {
	max := cfg.MaxMessages
	if max <= 0 {
		max = defaultMaxMessages
	}

	store, err := openStore(cfg.DataDir, max)
	if err != nil {
		logger.Errorf("gotify compat: data store unavailable (%v); falling back to in-memory", err)
		store = newMemoryStore(max)
	}

	version := cfg.Version
	if version == "" {
		version = "latest"
	}

	svc := &Service{
		store:   store,
		hub:     newHub(),
		version: version,
	}

	if cfg.ClientToken != "" {
		hash := hashToken(cfg.ClientToken)
		if err := store.SetTokenHash(hash); err != nil {
			return nil, fmt.Errorf("persist token hash: %w", err)
		}
		svc.tokenHash = hash
		// Operator-supplied token: never logged.
		return svc, nil
	}

	// Auto-generated token: prefer the persisted one for stability across restarts.
	hash, err := store.TokenHash()
	if err != nil {
		return nil, fmt.Errorf("load token hash: %w", err)
	}
	if len(hash) == 0 {
		tok, err := generateToken()
		if err != nil {
			return nil, fmt.Errorf("generate token: %w", err)
		}
		h := hashToken(tok)
		// best-effort persistence of the plaintext for recovery & stability
		if err := store.SetTokenHash(h); err != nil {
			return nil, fmt.Errorf("persist token hash: %w", err)
		}
		if ptErr := store.SetAutoToken(tok); ptErr != nil {
			logger.Warnf("monitor: could not persist auto client token: %v", ptErr)
		}
		hash = h
		svc.autoToken = tok
	}
	svc.tokenHash = hash
	return svc, nil
}

// ValidateToken reports whether the supplied token matches the stored hash.
func (s *Service) ValidateToken(token string) bool {
	return tokensEqual(hashToken(token), s.tokenHash)
}

// ClientToken returns the plaintext token only when it was auto-generated in
// this process's lifetime (so it can be logged exactly once). It returns ""
// for operator-supplied or pre-existing tokens.
func (s *Service) ClientToken() string {
	return s.autoToken
}

// Version returns the value served on GET /version.
func (s *Service) Version() string {
	return s.version
}

// Messages returns up to limit stored messages newest-first, optionally
// filtered to ID < since.
func (s *Service) Messages(limit int, since uint64) ([]Message, error) {
	return s.store.Recent(limit, since)
}

// Subscribe registers a live message consumer. The returned channel is closed
// by the returned unsubscribe func.
func (s *Service) Subscribe() (<-chan Message, func()) {
	return s.hub.Subscribe()
}

// Publish persists a message and fan-outs it to all subscribers. Persistence
// errors are returned to the caller, which must never block or alter the
// originating push.
func (s *Service) Publish(title, body string, priority int, extras map[string]interface{}) error {
	if extras == nil {
		extras = map[string]interface{}{}
	}
	m := Message{
		AppID:    defaultAppID,
		Title:    title,
		Message:  body,
		Priority: priority,
		Extras:   extras,
		Date:     time.Now().Format(time.RFC3339Nano),
	}
	id, err := s.store.Add(&m)
	if err != nil {
		return err
	}
	m.ID = id
	s.hub.Publish(m)
	return nil
}
