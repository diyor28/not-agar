package sockethub

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	roomId string
	Socket *websocket.Conn
	lock   *sync.Mutex
}

func (conn *Client) WriteJSON(v interface{}) error {
	conn.lock.Lock()
	err := conn.Socket.WriteJSON(v)
	conn.lock.Unlock()
	return err
}

func (conn *Client) ReadJSON() (ServerRequest, error) {
	var request ServerRequest
	conn.lock.Lock()
	err := conn.Socket.ReadJSON(&request)
	conn.lock.Unlock()
	return request, err
}

func (conn *Client) Emit(event string, data interface{}) error {
	return conn.WriteJSON(ServerResponse{event, data})
}
