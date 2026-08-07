package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/wallleap/hotify-bark-server/internal/gotifycompat"
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
		// Delete endpoints: single id or wipe-all, token-authenticated like
		// gotify's DELETE /message and DELETE /message/:id.
		router.Delete("/message", routeGotifyMessageDeleteAll)
		router.Delete("/message/:id", routeGotifyMessageDeleteOne)
		// Live subscription endpoint used by hotify-bridge /stream monitoring.
		router.Get("/stream", routeGotifyStreamUpgrade, fiberws.New(routeGotifyStream))

		// --- Device-scoped variants -------------------------------------
		// Each supports per-device history & live stream. Authentication still
		// uses the global client token (device isolation is not a credential).
		// Static path segments take priority over the legacy push_compat
		// `/:device_key/:body` routes, so these never clash.
		router.Get("/:device_key/version", routeGotifyDeviceVersion)
		router.Get("/:device_key/message", routeGotifyDeviceMessage)
		router.Delete("/:device_key/message", routeGotifyDeviceMessageDeleteAll)
		router.Delete("/:device_key/message/:id", routeGotifyDeviceMessageDeleteOne)
		router.Get("/:device_key/stream", routeGotifyDeviceStreamUpgrade, fiberws.New(routeGotifyStream))
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

// routeGotifyDeviceVersion serves the same version for a device-scoped probe.
// Like the global /version, it requires no token.
func routeGotifyDeviceVersion(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	return c.JSON(map[string]string{"version": gotifyService.Version()})
}

// messageQuery holds the shared limit/since pagination parsed from a request.
type messageQuery struct {
	limit int
	since uint64
}

// parseMessageQuery reads limit (default 100, max 200, min 1) and since
// (ID < since) from the request, following the global /message semantics.
func parseMessageQuery(c *fiber.Ctx) messageQuery {
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
	return messageQuery{limit: limit, since: since}
}

// messageResponse wraps a message list in the gotify paging envelope.
func messageResponse(messages []gotifycompat.Message, q messageQuery) map[string]interface{} {
	return map[string]interface{}{
		"paging": map[string]interface{}{
			"size":  len(messages),
			"limit": q.limit,
			"since": q.since,
		},
		"messages": messages,
	}
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

	q := parseMessageQuery(c)

	messages, err := gotifyService.Messages(q.limit, q.since)
	if err != nil {
		return c.Status(500).JSON(failed(500, "get messages failed: %v", err))
	}
	return c.JSON(messageResponse(messages, q))
}

// routeGotifyDeviceMessage is the device-scoped GET /message: only messages of
// the device_key in the URL path are returned.
func routeGotifyDeviceMessage(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}

	q := parseMessageQuery(c)
	device := c.Params("device_key")

	messages, err := gotifyService.MessagesByDevice(device, q.limit, q.since)
	if err != nil {
		return c.Status(500).JSON(failed(500, "get messages failed: %v", err))
	}
	return c.JSON(messageResponse(messages, q))
}

// routeGotifyMessageDeleteAll wipes the whole message history (gotify's
// DELETE /message), token-authenticated.
func routeGotifyMessageDeleteAll(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}
	if err := gotifyService.DeleteAllMessages(); err != nil {
		return c.Status(500).JSON(failed(500, "delete messages failed: %v", err))
	}
	return c.JSON(success())
}

// routeGotifyMessageDeleteOne removes a single message (gotify's
// DELETE /message/:id); 404 when the id does not exist.
func routeGotifyMessageDeleteOne(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(failed(400, "invalid message id: %v", err))
	}
	existed, err := gotifyService.DeleteMessage(id)
	if err != nil {
		return c.Status(500).JSON(failed(500, "delete message failed: %v", err))
	}
	if !existed {
		return c.Status(404).JSON(failed(404, "message not found"))
	}
	return c.JSON(success())
}

// routeGotifyDeviceMessageDeleteAll wipes only the given device's message
// history; other devices' messages are untouched.
func routeGotifyDeviceMessageDeleteAll(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}
	if err := gotifyService.DeleteAllMessagesByDevice(c.Params("device_key")); err != nil {
		return c.Status(500).JSON(failed(500, "delete messages failed: %v", err))
	}
	return c.JSON(success())
}

// routeGotifyDeviceMessageDeleteOne removes a single message that belongs to
// the device in the URL path; 404 when it does not exist or belongs to another
// device.
func routeGotifyDeviceMessageDeleteOne(c *fiber.Ctx) error {
	if gotifyService == nil {
		return c.Status(503).JSON(failed(503, "gotify compat not initialized"))
	}
	if !gotifyService.ValidateToken(gotifyToken(c)) {
		return c.Status(401).JSON(failed(401, "unauthorized"))
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(failed(400, "invalid message id: %v", err))
	}
	existed, err := gotifyService.DeleteMessageByDevice(c.Params("device_key"), id)
	if err != nil {
		return c.Status(500).JSON(failed(500, "delete message failed: %v", err))
	}
	if !existed {
		return c.Status(404).JSON(failed(404, "message not found"))
	}
	return c.JSON(success())
}

// routeGotifyStreamUpgrade validates the token before the WebSocket upgrade so
// that unauthorized clients get a 401 handshake response (gotify behaviour).
func routeGotifyStreamUpgrade(c *fiber.Ctx) error {
	return routeGotifyStreamUpgradeDevice(c, "")
}

// routeGotifyDeviceStreamUpgrade is the device-scoped /stream upgrade: the
// device_key is stashed for the stream handler, then the same token validation
// as the global endpoint applies.
func routeGotifyDeviceStreamUpgrade(c *fiber.Ctx) error {
	return routeGotifyStreamUpgradeDevice(c, c.Params("device_key"))
}

func routeGotifyStreamUpgradeDevice(c *fiber.Ctx, device string) error {
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
	if device != "" {
		c.Locals("device_key", device)
	}
	return c.Next()
}

// routeGotifyStream serves the live WebSocket stream of bare gotify message
// JSON frames (no event/socketConnected envelope), matching what the
// hotify-bridge subscriber expects. When a device_key is present in the fiber
// locals, only that device's messages are streamed. It replies to client pings
// and reaches out with server pings so idle NAT'd connections survive. All
// writes happen in this goroutine (no concurrent writes on the connection).
func routeGotifyStream(conn *fiberws.Conn) {
	device := ""
	if v := conn.Locals("device_key"); v != nil {
		device, _ = v.(string)
	}
	var ch <-chan gotifycompat.Message
	var unsubscribe func()
	if device != "" {
		ch, unsubscribe = gotifyService.SubscribeByDevice(device)
	} else {
		ch, unsubscribe = gotifyService.Subscribe()
	}
	defer unsubscribe()
	if barkMetrics != nil {
		barkMetrics.SetActiveStreams(float64(gotifyService.SubscriberCount()))
		defer barkMetrics.SetActiveStreams(float64(gotifyService.SubscriberCount()))
	}

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
