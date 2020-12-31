package sockethub

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type ServerResponse struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ServerRequest struct {
	Event string      `json:"Event"` // move
	Data  interface{} `json:"Data"`
}

type EventHandler struct {
	Event    string
	Callback func(data interface{}, client *Client)
}

type Hub struct {
	connections   []*Client
	eventHandlers []EventHandler
}

func NewHub() *Hub {
	h := Hub{}
	return &h
}

func (h *Hub) Run() {
	for {
		var validConnections []*Client
		var wg sync.WaitGroup
		for _, conn := range h.connections {
			wg.Add(1)
			go func(client *Client) {
				defer wg.Done()
				request, err := client.ReadJSON()
				if err != nil {
					log.Println(err)
				} else {
					validConnections = append(validConnections, client)
				}
				h.receivedEvent(request, client)
			}(conn)
		}
		wg.Wait()
		h.connections = validConnections
	}
}

func (h *Hub) AddConnection(ws *websocket.Conn, roomId string) {
	conn := &Client{roomId, ws, &sync.Mutex{}}
	h.connections = append(h.connections, conn)
}

func (h *Hub) receivedEvent(request ServerRequest, client *Client) {
	for _, event := range h.eventHandlers {
		if event.Event == request.Event {
			event.Callback(request.Data, client)
		}
	}
}

func (h *Hub) On(event string, callback func(data interface{}, client *Client)) {
	h.eventHandlers = append(h.eventHandlers, EventHandler{event, callback})
}

func (h *Hub) CloseChannel(roomId string) {

}

func (h *Hub) Emit(event string, data interface{}, roomId string) {
	var validConnections []*Client
	for _, conn := range h.connections {
		if roomId == "" || conn.roomId == roomId {
			if err := conn.Emit(event, data); err != nil {
				log.Println(err)
			} else {
				validConnections = append(validConnections, conn)
			}
		}
	}
	h.connections = validConnections
}

func (h *Hub) Publish(event string, data interface{}) {
	h.Emit(event, data, "")
}
