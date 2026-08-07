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

// isAuthFreePath reports whether p is (or is under) a Basic-Auth whitelisted
// path, with exact-or-subpath semantics so lookalikes like /messageevil are
// NOT whitelisted.
func isAuthFreePath(urlPrefix, p string) bool {
	for _, item := range authFreeRouters {
		base := path.Join(urlPrefix, item)
		if p == base || strings.HasPrefix(p, base+"/") {
			return true
		}
	}
	return false
}

func routerAuth(user, passwd string, router fiber.Router, urlPrefix string) {
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
