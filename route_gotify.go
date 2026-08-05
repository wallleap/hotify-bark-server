package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/finb/bark-server/v2/internal/gotifycompat"
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
)

// gotifyService is the lazily-initialized gotify-compatible monitoring service.
var gotifyService *gotifycompat.Service

func init() {
	registerRoute("gotify", func(router fiber.Router) {
		// Probe endpoint used by hotify-bridge localhost autodetect.
		router.Get("/version", routeGotifyVersion)
		// Read/history endpoint used by hotify-bridge backfill & watermark init.
		router.Get("/message", routeGotifyMessage)
		// Live subscription endpoint used by hotify-bridge /stream monitoring.
		router.Get("/stream", routeGotifyStreamUpgrade, fiberws.New(routeGotifyStream))
	})
}

// gotifyToken extracts a client token following gotify's precedence:
// query → X-Gotify-Key header → Authorization: Bearer.
func gotifyToken(c *fiber.Ctx) string {
	if t := c.Query("token"); t != "" {
		return t
	}
	if t := c.Get("X-Gotify-Key"); t != "" {
		return t
	}
	if auth := c.Get(fiber.HeaderAuthorization); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// routeGotifyVersion serves {"version":"..."} like gotify's GET /version.
func routeGotifyVersion(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	return c.JSON(map[string]string{"version": gotifyService.Version()})
}

// routeGotifyMessage serves {"messages":[...]} newest-first, like gotify's
// GET /message. Supports limit (default 100, max 200) and since (ID < since).
func routeGotifyMessage(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}

	limit := c.QueryInt("limit", 100)
	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}
	var since uint64
	if s := c.Query("since"); s != "" {
		since, _ = strconv.ParseUint(s, 10, 64)
	}

	messages, err := gotifyService.Messages(limit, since)
	if err != nil {
		return c.Status(500).JSON(failed(500, "get messages failed: %v", err))
	}
	return c.JSON(map[string]interface{}{
		"paging": map[string]interface{}{
			"size":  len(messages),
			"limit": limit,
			"since": since,
		},
		"messages": messages,
	})
}

// routeGotifyStreamUpgrade validates the token before the WebSocket upgrade so
// that unauthorized clients get a 401 handshake response (gotify behaviour).
func routeGotifyStreamUpgrade(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !fiberws.IsWebSocketUpgrade(c) {
		return c.Status(400).JSON(failed(400, "websocket upgrade expected"))
	}
	token := gotifyToken(c)
	if !gotifyService.ValidateToken(token) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}
	return c.Next()
}

// routeGotifyStream serves the live WebSocket stream of bare gotify message
// JSON frames (no event/socketConnected envelope), matching what the
// hotify-bridge subscriber expects. It replies to client pings and reaches out
// with server pings so idle NAT'd connections survive. All writes happen in
// this goroutine (no concurrent writes on the connection).
func routeGotifyStream(conn *fiberws.Conn) {
	ch, unsubscribe := gotifyService.Subscribe()
	defer unsubscribe()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	// hotify-bridge sends a WebSocket ping every 20s; each received ping/pong
	// refreshes the read deadline so the connection survives indefinitely,
	// while a truly silent dead peer is reaped after 60s.
	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return conn.WriteControl(fiberws.PongMessage, []byte(appData), time.Now().Add(5*time.Second))
	})
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ping := time.NewTicker(45 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-done:
			return
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteControl(fiberws.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				return
			}
		case m, ok := <-ch:
			if !ok {
				return
			}
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(m); err != nil {
				return
			}
		}
	}
}
