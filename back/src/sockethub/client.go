package sockethub

import (
	"github.com/gorilla/websocket"
	"log"
)

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Client struct {
	channels []string
	socket   *websocket.Conn
	hub      *Hub
	read     chan []byte
	send     chan []byte
}

func (conn *Client) writer() {
	defer func() {
		err := conn.socket.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for message := range conn.send {
		if err := conn.socket.WriteMessage(websocket.BinaryMessage, message); err != nil {
			log.Println(err)
		}
	}
}

func (conn *Client) reader() {
	defer func() {
		conn.hub.unregister <- conn
		err := conn.socket.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		messagesType, data, err := conn.socket.ReadMessage()
		if messagesType != websocket.BinaryMessage {
			log.Println("Expected BinaryMessage, got: ", messagesType)
			continue
		}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		conn.hub.readBuffer <- ReadMessage{data, conn}
	}
}

func (conn *Client) Join(room string) {
	conn.channels = append(conn.channels, room)
}

func (conn *Client) Leave(channel string) {
	for i, v := range conn.channels {
		if v == channel {
			conn.channels = remove(conn.channels, i)
			break
		}
	}
}

func (conn *Client) IsInChannel(channel string) bool {
	for _, v := range conn.channels {
		if v == channel {
			return true
		}
	}
	return false
}

func (conn *Client) Emit(data []byte) {
	conn.send <- data
}
