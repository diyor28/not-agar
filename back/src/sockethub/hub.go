package sockethub

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type BroadcastMessage struct {
	Message
	roomId string
}

type EventHandler struct {
	Event    string
	Callback func(data interface{}, roomId string)
}

type Hub struct {
	connections   []*Client
	eventHandlers []EventHandler
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast  chan BroadcastMessage
	readBuffer chan BroadcastMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	h := Hub{
		readBuffer: make(chan BroadcastMessage),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	return &h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if client.roomId != message.roomId {
					continue
				}
				client.send <- message.Message
			}
		case message := <-h.readBuffer:
			go h.receivedEvent(message)
		}
	}
}

func (h *Hub) AddConnection(ws *websocket.Conn, roomId string) {
	client := &Client{roomId: roomId, socket: ws, send: make(chan Message), hub: h}
	go client.reader()
	go client.writer()
	h.register <- client
}

func (h *Hub) receivedEvent(message BroadcastMessage) {
	for _, event := range h.eventHandlers {
		if event.Event == message.Event {
			event.Callback(message.Data, message.roomId)
		}
	}
}

func (h *Hub) On(event string, callback func(data interface{}, roomId string)) {
	h.eventHandlers = append(h.eventHandlers, EventHandler{event, callback})
}

func (h *Hub) Emit(event string, data interface{}, roomId string) {
	h.broadcast <- BroadcastMessage{Message{event, data}, roomId}
}
