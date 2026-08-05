package gotifycompat

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

// Client token handling.
//
// Security model:
//   - Only the SHA-256 hash of the token is ever used for validation, and
//     comparison is constant-time to avoid timing attacks.
//   - An operator-supplied token (via BARK_SERVER_GOTIFY_CLIENT_TOKEN) is
//     never persisted in plaintext; only its hash is stored.
//   - A token that is auto-generated for convenience is persisted in plaintext
//     inside the 0600 gotify.db file so it stays stable across restarts
//     (otherwise the bridge config would break on every restart). It is logged
//     exactly once, at generation time.

const tokenBytes = 32

// hashToken returns the SHA-256 digest of a token.
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

// tokensEqual performs a constant-time comparison of two token hashes.
func tokensEqual(a, b []byte) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

// generateToken creates a cryptographically random, URL-safe client token.
func generateToken() (string, error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
