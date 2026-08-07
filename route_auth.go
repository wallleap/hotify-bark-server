package main

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberbasicauth "github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/mritd/logger"
)

// authFreeRouters are paths exempt from Basic Auth; they carry their own
// authentication (gotify token on /message & /stream, none needed on the
// probe endpoints).
var authFreeRouters = []string{"/ping", "/register", "/healthz", "/version", "/message", "/stream"}

// authFreeSuffixes are single device-key-scoped path suffixes exempt from
// Basic Auth, matching /:device_key/<suffix>. They carry the same per-route
// authentication as their global counterparts (gotify token / none).
var authFreeSuffixes = []string{"/version", "/message", "/stream"}

// basicAuthEnabled reports whether Basic Auth was configured at startup. It
// gates features that should not be exposed to an unauthenticated network
// (e.g. the device count on GET /info).
var basicAuthEnabled bool

// isAuthFreePath reports whether p is (or is under) a Basic-Auth whitelisted
// path, with exact-or-subpath semantics so lookalikes like /messageevil are
// NOT whitelisted. Device-scoped paths of the form /<device>/message,
// /<device>/stream and /<device>/version are also whitelisted.
func isAuthFreePath(urlPrefix, p string) bool {
	for _, item := range authFreeRouters {
		base := path.Join(urlPrefix, item)
		if p == base || strings.HasPrefix(p, base+"/") {
			return true
		}
	}
	// /<device>/<suffix> — one non-empty segment then a whitelisted suffix.
	trimmed := strings.TrimPrefix(p, urlPrefix)
	trimmed = strings.TrimPrefix(trimmed, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 || parts[0] == "" {
		return false
	}
	for _, suffix := range authFreeSuffixes {
		if parts[1] == strings.TrimPrefix(suffix, "/") {
			return true
		}
	}
	return false
}

func routerAuth(user, passwd string, router fiber.Router, urlPrefix string) {
	basicAuthEnabled = user != "" && passwd != ""
	if user == "" && passwd == "" {
		logger.Warn("************************************************************")
		logger.Warn("Hotify-Bark Server Has NO Basic Auth.")
		logger.Warn("PUBLIC deployments should set BARK_SERVER_BASIC_AUTH_USER/PASSWORD.")
		logger.Warn("/push, /register, /mcp* and /:device_key are OPEN to everyone.")
		logger.Warn("************************************************************")
		return
	}

	logger.Info("Hotify-Bark Server Has Basic Auth Enabled.")
	basicAuth := fiberbasicauth.New(fiberbasicauth.Config{
		Users: map[string]string{user: passwd},
		Realm: "Coffee Time",
		Unauthorized: func(c *fiber.Ctx) error {
			if isAuthFreePath(urlPrefix, c.Path()) {
				return c.Next()
			}
			return c.Status(418).SendString("I'm a teapot")
		},
	})

	router.Use("/+", basicAuth)
}
