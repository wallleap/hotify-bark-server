package database

import (
	"path/filepath"
	"testing"
)

// TestMemBaseCRUD exercises the in-memory database used in tests/serverless-adjacent
// scenarios. MemBase holds a single package-level key/token pair keyed by the
// non-empty package cacheKey.
func TestMemBaseCRUD(t *testing.T) {
	db := NewMemBase()

	// save a token under the known key
	key, err := db.SaveDeviceTokenByKey(cacheKey, "token-v1")
	if err != nil {
		t.Fatalf("SaveDeviceTokenByKey failed: %v", err)
	}
	if key != cacheKey {
		t.Fatalf("want returned key %q, got %q", cacheKey, key)
	}

	tok, err := db.DeviceTokenByKey(cacheKey)
	if err != nil {
		t.Fatalf("DeviceTokenByKey failed: %v", err)
	}
	if tok != "token-v1" {
		t.Fatalf("want token %q, got %q", "token-v1", tok)
	}

	if n, _ := db.CountAll(); n != 1 {
		t.Fatalf("want CountAll=1, got %d", n)
	}

	// Updating the token should replace it, not append.
	if _, err := db.SaveDeviceTokenByKey(cacheKey, "token-v2"); err != nil {
		t.Fatalf("SaveDeviceTokenByKey(update) failed: %v", err)
	}
	tok, _ = db.DeviceTokenByKey(cacheKey)
	if tok != "token-v2" {
		t.Fatalf("want updated token %q, got %q", "token-v2", tok)
	}

	// Delete clears the token.
	if err := db.DeleteDeviceByKey(cacheKey); err != nil {
		t.Fatalf("DeleteDeviceByKey failed: %v", err)
	}
	if _, err := db.DeviceTokenByKey(cacheKey); err == nil {
		t.Fatal("want error after delete, got nil")
	}
}

func TestMemBaseErrors(t *testing.T) {
	db := NewMemBase()
	// Unknown key must error.
	if _, err := db.DeviceTokenByKey("nope"); err == nil {
		t.Fatal("want error for unknown key, got nil")
	}
	// Deleting an unknown key errors.
	if err := db.DeleteDeviceByKey("nope"); err == nil {
		t.Fatal("want error deleting unknown key, got nil")
	}
}

func TestEnvBaseCRUD(t *testing.T) {
	t.Setenv("BARK_KEY", "env-key")
	t.Setenv("BARK_DEVICE_TOKEN", "env-token")

	db := NewEnvBase()
	tok, err := db.DeviceTokenByKey("env-key")
	if err != nil {
		t.Fatalf("DeviceTokenByKey failed: %v", err)
	}
	if tok != "env-token" {
		t.Fatalf("want token %q, got %q", "env-token", tok)
	}

	// EnvBase mirrors the configured key/token and rejects anything else.
	if _, err := db.DeviceTokenByKey("other"); err == nil {
		t.Fatal("want error for mismatched key, got nil")
	}

	saved, err := db.SaveDeviceTokenByKey("env-key", "env-token")
	if err != nil || saved != "env-key" {
		t.Fatalf("SaveDeviceTokenByKey want (env-key,nil), got (%q,%v)", saved, err)
	}

	// EnvBase does not support deletes.
	if err := db.DeleteDeviceByKey("env-key"); err == nil {
		t.Fatal("want error from DeleteDeviceByKey, got nil")
	}
}

// TestBboltCRUD validates the default on-disk database against a scratch dir.
// Bbolt is a process-wide singleton here, so this is a single combined flow.
func TestBboltCRUD(t *testing.T) {
	dir := t.TempDir()
	db := NewBboltdb(filepath.Join(dir, "data"))

	key, err := db.SaveDeviceTokenByKey("", "token-device-1")
	if err != nil {
		t.Fatalf("SaveDeviceTokenByKey failed: %v", err)
	}
	if key == "" {
		t.Fatal("expected newly generated key for empty-key registration")
	}
	generated := key

	// SaveDeviceTokenByKey on an existing key keeps the key.
	key2, err := db.SaveDeviceTokenByKey(generated, "token-device-2")
	if err != nil {
		t.Fatalf("SaveDeviceTokenByKey(update) failed: %v", err)
	}
	if key2 != generated {
		t.Fatalf("want key preserved %q, got %q", generated, key2)
	}

	tok, err := db.DeviceTokenByKey(generated)
	if err != nil {
		t.Fatalf("DeviceTokenByKey failed: %v", err)
	}
	if tok != "token-device-2" {
		t.Fatalf("want updated token %q, got %q", "token-device-2", tok)
	}

	if n, _ := db.CountAll(); n != 1 {
		t.Fatalf("want CountAll=1, got %d", n)
	}

	if err := db.DeleteDeviceByKey(generated); err != nil {
		t.Fatalf("DeleteDeviceByKey failed: %v", err)
	}
	if _, err := db.DeviceTokenByKey(generated); err == nil {
		t.Fatal("want error after delete, got nil")
	}
}

func TestBboltErrors(t *testing.T) {
	db := NewBboltdb(t.TempDir())
	if _, err := db.DeviceTokenByKey("does-not-exist"); err == nil {
		t.Fatal("want error for missing key, got nil")
	}
}
