package sockethub

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Client struct {
	IsClosed bool
	channels []string
	socket   *websocket.Conn
	hub      *Hub
	read     chan []byte
	send     chan []byte
}

func (conn *Client) writer() {
	defer func() {
		if conn.IsClosed {
			return
		}
		err := conn.socket.Close()
		if err != nil {
			log.Println("error while closing: ", err)
		}
	}()
	for message := range conn.send {
		if conn.IsClosed {
			break
		}
		if err := conn.socket.WriteMessage(websocket.BinaryMessage, message); err != nil {
			log.Println(err)
			break
		}
	}
}

func (conn *Client) reader() {
	defer func() {
		conn.IsClosed = true
		conn.hub.unregister <- conn
		err := conn.socket.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		messagesType, data, err := conn.socket.ReadMessage()
		if messagesType == websocket.CloseMessage {
			log.Println("Closing connection")
			break
		}
		if err != nil {
			log.Printf("error: %v", err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
		}
		if messagesType != websocket.BinaryMessage {
			log.Println("Expected BinaryMessage, got: ", messagesType)
			continue
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

func (conn *Client) Emit(data []byte) error {
	if conn.IsClosed {
		return errors.New("connection is closed")
	}
	conn.send <- data
	return nil
}
