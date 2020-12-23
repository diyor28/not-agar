package main

import (
	"encoding/json"
	"github.com/diyor28/not-agar/gamengine"
	"github.com/frankenbeanies/uuid4"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
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

func parseRequest(conn *websocket.Conn) (gamengine.Player, int, error) {
	var request gamengine.ServerRequest
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return gamengine.Player{}, messageType, err
	}
	_ = json.Unmarshal(p, &request)
	var moveEvent = gamengine.MoveEvent{}
	if err = mapstructure.Decode(request.Data, &moveEvent); err != nil {
		log.Println(err)
		return gamengine.Player{}, messageType, err
	}
	if moveEvent.Uuid == "" {
		return gamengine.Player{}, messageType, err
	}
	player := gameMap.GetPlayer(moveEvent.Uuid)
	if player.Uuid == "" {
		return gamengine.Player{}, messageType, err
	}
	gameMap.UpdatePlayer(moveEvent)
	return player, messageType, err
}

func reader(conn *websocket.Conn) {
	connection := gameMap.AddConnection(conn)
	for {
		receivedPlayer, _, err := parseRequest(conn)
		if err != nil {
			log.Println(err)
			return
		}
		if err := connection.WriteJSON(gameMap.ServerResponse(&receivedPlayer)); err != nil {
			log.Println(err)
			return
		}
	}
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
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
