package sockethub

import (
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	Data     []byte
	channels []string
}

type ReadMessage struct {
	Data   []byte
	Client *Client
}

type EventHandler struct {
	Event    string
	Callback func(data interface{}, client *Client)
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast  chan BroadcastMessage
	readBuffer chan ReadMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	onMessage  func(data []byte, client *Client)
}

func NewHub() *Hub {
	h := Hub{
		readBuffer: make(chan ReadMessage),
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
				isCorrectRecipient := false
				for _, c := range message.channels {
					isCorrectRecipient = isCorrectRecipient || client.IsInChannel(c)
				}
				if !isCorrectRecipient {
					continue
				}
				client.send <- message.Data
			}
		case message := <-h.readBuffer:
			go h.onMessage(message.Data, message.Client)
		}
	}
}

func (h *Hub) OnMessage(callback func(data []byte, client *Client)) {
	h.onMessage = callback
}

func (h *Hub) AddConnection(ws *websocket.Conn) *Client {
	client := &Client{socket: ws, send: make(chan []byte), hub: h}
	go client.reader()
	go client.writer()
	h.register <- client
	return client
}

func (h *Hub) Emit(data []byte, channel string) {
	h.broadcast <- BroadcastMessage{data, []string{channel}}
}
