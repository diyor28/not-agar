package main

import (
	"encoding/json"
	"fmt"
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

func playerWS(w http.ResponseWriter, r *http.Request) {
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

func adminWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	gameMap.Hub.AddConnection(ws, "admin")
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
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
	processes := 4
	fmt.Println("Setting max processes:", processes)
	runtime.GOMAXPROCS(processes)
	go gameMap.Run()
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/players", createPlayer).Methods("POST", "OPTIONS")
	router.HandleFunc("/player-ws/{uuid}/", playerWS)
	router.HandleFunc("/admin", adminWS)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	fmt.Println("Listening on port 3100")
	err := http.ListenAndServe(":3100", router)
	log.Fatal(err)
}
