package main

import (
	"encoding/json"
	"github.com/diyor28/not-agar/gamengine"
	"github.com/frankenbeanies/uuid4"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
)

var gameMap = gamengine.GameMap{
	GameId:  uuid4.New().String(),
	Players: []gamengine.Player{},
	Foods:   []gamengine.Food{},
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func parseRequest(conn *gamengine.Connection) (gamengine.ServerRequest, error) {
	var request gamengine.ServerRequest
	err := conn.Socket.ReadJSON(&request)
	return request, err
}

func reader(conn *websocket.Conn) {
	connection := gameMap.AddConnection(conn)
	for {
		request, err := parseRequest(connection)
		if err != nil {
			log.Println(err)
			return
		}
		go gameMap.HandleEvent(request, connection)
	}
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	reader(ws)
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	//setupResponse(&w, r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	var data gamengine.Player
	_ = json.NewDecoder(r.Body).Decode(&data)
	result := gameMap.CreatePlayer(data.Nickname, false)
	_ = json.NewEncoder(w).Encode(result)
}

func main() {
	runtime.GOMAXPROCS(4)
	go gameMap.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ws", websocketEndpoint)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/players", createPlayer).Methods("POST", "OPTIONS")
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":3100", router))
}
