package main

import (
	"github.com/diyor28/not-agar/cmd/gamengine"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
)

var gameMap = gamengine.NewGameMap(50)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 2048,
}

func playerWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := gameMap.Hub.AddConnection(ws)
	client.Join("anonymous")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Proto)
		next.ServeHTTP(w, r)
	})
}

func main() {
	processes := 4
	log.Println("Setting max processes:", processes)
	runtime.GOMAXPROCS(processes)
	go gameMap.Run()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)
	router.HandleFunc("/player-ws", playerWS)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	log.Println("Listening on port 3100")
	err := http.ListenAndServe(":3100", router)
	log.Fatal(err)
}
