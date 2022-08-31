package sockethub

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	roomId string
	socket *websocket.Conn
	hub    *Hub
	read   chan Message
	send   chan Message
}

func (conn *Client) writer() {
	//fmt.Println("writer active")
	defer func() {
		conn.socket.Close()
	}()
	for data := range conn.send {
		//fmt.Println("writing data to socket", data.Event)
		if err := conn.socket.WriteJSON(data); err != nil {
			log.Println(err)
		}
	}
}

func (conn *Client) reader() {
	defer func() {
		conn.hub.unregister <- conn
		conn.socket.Close()
	}()
	for {
		var request Message
		err := conn.socket.ReadJSON(&request)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		conn.hub.readBuffer <- BroadcastMessage{request, conn.roomId}
	}
}

func (conn *Client) Emit(event string, data interface{}) {
	conn.send <- Message{event, data}
}
