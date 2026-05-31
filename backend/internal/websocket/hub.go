package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uint
	Role   string
}

type outboundMessage struct {
	payload []byte
	roles   map[string]bool
	userID  *uint
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan outboundMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		// broadcast is intentionally unbuffered: the hub is the single writer gate,
		// while each client has its own buffered Send channel below.
		broadcast:  make(chan outboundMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("WS: Client Registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Println("WS: Client Unregistered")
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if !message.matches(client) {
					continue
				}
				select {
				case client.Send <- message.payload:
				default:
					// If the per-client buffer is full, the browser is too slow or gone.
					// Closing it protects the hub from being blocked by one stale connection.
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (m outboundMessage) matches(client *Client) bool {
	if m.userID != nil && client.UserID == *m.userID {
		return true
	}
	if len(m.roles) > 0 {
		return m.roles[strings.ToLower(client.Role)]
	}
	return m.userID == nil
}

func encodeMessage(topic string, data interface{}) []byte {
	msg := map[string]interface{}{
		"topic": topic,
		"data":  data,
	}
	bytes, _ := json.Marshal(msg)
	return bytes
}

func (h *Hub) Broadcast(topic string, data interface{}) {
	h.broadcast <- outboundMessage{payload: encodeMessage(topic, data)}
}

func (h *Hub) BroadcastToRoles(topic string, data interface{}, roles ...string) {
	if len(roles) == 0 {
		h.Broadcast(topic, data)
		return
	}
	roleMap := make(map[string]bool, len(roles))
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role != "" {
			roleMap[role] = true
		}
	}
	if len(roleMap) == 0 {
		return
	}
	h.broadcast <- outboundMessage{payload: encodeMessage(topic, data), roles: roleMap}
}

func (h *Hub) BroadcastToUser(topic string, data interface{}, userID uint) {
	if userID == 0 {
		return
	}
	h.broadcast <- outboundMessage{payload: encodeMessage(topic, data), userID: &userID}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, allowedOrigins ...string) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isAllowedOrigin(r, allowedOrigins...)
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Send is buffered so bursts of dashboard events do not block payment/webhook code.
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	if role, ok := r.Context().Value("role").(string); ok {
		client.Role = strings.ToLower(role)
	}
	switch userID := r.Context().Value("user_id").(type) {
	case uint:
		client.UserID = userID
	case int:
		if userID > 0 {
			client.UserID = uint(userID)
		}
	}
	client.Hub.register <- client

	go client.writePump()
	go client.readPump()
}

func isAllowedOrigin(r *http.Request, allowedOrigins ...string) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	for _, allowed := range allowedOrigins {
		allowed = strings.TrimRight(strings.TrimSpace(allowed), "/")
		if allowed != "" && strings.EqualFold(origin, allowed) {
			return true
		}
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return parsed.Scheme == "http" && (host == "localhost" || host == "127.0.0.1")
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
