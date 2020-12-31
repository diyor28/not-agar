package main

import (
	"encoding/json"
	"github.com/diyor28/not-agar/gamengine"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
)

var gameMap = gamengine.NewGameMap()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 2048,
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	gameMap.Hub.AddConnection(ws, uuid)
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
	router.HandleFunc("/player-ws/{uuid}/", websocketEndpoint)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/players", createPlayer).Methods("POST", "OPTIONS")
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":3100", router))
}
