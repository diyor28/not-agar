package gamengine

import (
	"errors"
	"github.com/diyor28/not-agar/randomname"
	"github.com/frankenbeanies/uuid4"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const maxXY = 5000
const foodWeight = 25
const surfaceArea = maxXY * maxXY
const minXY = 0
const windowSize = 1200
const minSpeed = 2
const maxSpeed = 8
const maxNumFood = surfaceArea / 50000
const minWeight = 40
const maxWeight = minWeight * 15
const minZoom = 0.8
const maxZoom = 1.0
const maxPlayers = 20
const maxSpikes = maxPlayers
const statsNumber = 10
const speedWeightLimit = 300

type ServerResponse struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ServerRequest struct {
	Event string      `json:"Event"` // move
	Data  interface{} `json:"Data"`
}

type MoveEvent struct {
	Uuid string  `json:"uuid"`
	NewX float32 `json:"newX"`
	NewY float32 `json:"newY"`
}

type AccelerateEvent struct {
	Uuid string `json:"uuid"`
}

type MovedEvent struct {
	SelfPlayer SelfPlayer `json:"selfPlayer"`
	Players    []Player   `json:"players"`
	Foods      []Food     `json:"foods"`
	Spikes     []Spike    `json:"spikes"`
}

func getRandomColor() [3]int {
	return colors[rand.Intn(len(colors))]
}

func getRandomCoordinate() (float32, float32) {
	return float32(rand.Intn(maxXY)), float32(rand.Intn(maxXY))
}

func newWeight(weight1 float32, weight2 float32) float32 {
	nWeight := math.Sqrt(float64(weight1*weight1 + weight2*weight2))
	return float32(math.Min(nWeight, maxWeight))
}

func getSpeedFromWeight(weight float32) float32 {
	normalized := 1 - (weight-minWeight)/(speedWeightLimit-minWeight)
	newSpeed := float64(normalized*(maxSpeed-minSpeed) + minSpeed)
	return float32(math.Max(newSpeed, minSpeed))
}

var colors = [][3]int{{255, 21, 21}, {255, 243, 21}, {21, 87, 255}, {21, 255, 208}, {255, 21, 224}}

type Connection struct {
	ConnectionId string
	Socket       *websocket.Conn
	lock         *sync.Mutex
}

func (conn *Connection) WriteJSON(v interface{}) error {
	conn.lock.Lock()
	err := conn.Socket.WriteJSON(v)
	conn.lock.Unlock()
	return err
}

func (conn *Connection) Emit(event string, data interface{}) error {
	return conn.WriteJSON(ServerResponse{Event: event, Data: data})
}

type GameMap struct {
	GameId      string
	Players     []Player `json:"players"`
	Foods       []Food   `json:"foods"`
	Spikes      []Spike  `json:"spikes"`
	connections []Connection
}

func (gMap *GameMap) handleMoveEvent(data interface{}, conn *Connection) {
	var moveData MoveEvent
	if err := mapstructure.Decode(data, &moveData); err != nil {
		log.Fatal(err)
		return
	}
	player, err := gMap.UpdatePlayer(moveData.Uuid, moveData.NewX, moveData.NewY)
	if err != nil {
		log.Fatal("could not update player ", err)
		return
	}
	foods := gMap.nearByFood(player)
	players := gMap.nearByPlayers(player)
	response := MovedEvent{
		SelfPlayer: player.getSelfPlayer(),
		Foods:      foods,
		Players:    players,
	}
	if err := conn.Emit("moved", response); err != nil {
		log.Fatal("Socket emit", err)
	}
}

func (gMap *GameMap) handleAccelerate(data interface{}, conn *Connection) {
	var accelerateData AccelerateEvent
	if err := mapstructure.Decode(data, &accelerateData); err != nil {
		log.Fatal(err)
		return
	}
	player, err := gMap.GetPlayer(accelerateData.Uuid)
	if err != nil {
		log.Fatal("could not find player", err)
		return
	}
	player.Accelerating = true
}

func (gMap *GameMap) sendPong(data interface{}, conn *Connection) error {
	return conn.Emit("pong", data)
}

func (gMap *GameMap) HandleEvent(request ServerRequest, conn *Connection) {
	switch request.Event {
	case "move":
		gMap.handleMoveEvent(request.Data, conn)
	case "accelerate":
		gMap.handleAccelerate(request.Data, conn)
	case "ping":
		_ = gMap.sendPong(request.Data, conn)
	}
}

func (gMap *GameMap) AddConnection(conn *websocket.Conn) *Connection {
	connection := Connection{
		ConnectionId: uuid4.New().String(),
		Socket:       conn,
		lock:         &sync.Mutex{},
	}
	gMap.connections = append(gMap.connections, connection)
	return &connection
}

func (gMap *GameMap) populateBots() {
	maxBots := (maxPlayers - len(gMap.Players)) / 2
	//maxBots := 0
	currentBots := 0
	for _, player := range gMap.Players {
		if player.IsBot {
			currentBots++
		}
	}
	if currentBots < maxBots {
		for i := 0; i < maxBots; i++ {
			gMap.CreatePlayer(randomname.GenerateNickname(), true)
		}
	}
}

func (gMap *GameMap) populateSpikes() {
	for i := 0; i < maxSpikes; i++ {
		gMap.Spikes = append(gMap.Spikes, Spike{
			Uuid:  uuid4.New().String(),
			Color: getRandomColor(),
		})
	}
}

func (gMap *GameMap) populateFood() {
	var totalWeight float32 = 0
	maxFood := maxNumFood
	for _, p := range gMap.Players {
		totalWeight += p.Weight * p.Weight * math.Pi
	}
	for _, f := range gMap.Foods {
		totalWeight += f.Weight * f.Weight * math.Pi
	}
	//density := int(20 * totalWeight / surfaceArea)
	//fmt.Println("DENSITY", 20 * totalWeight / surfaceArea)
	//maxFood -= density
	if len(gMap.Foods) < maxFood {
		for i := 0; i < maxFood-len(gMap.Foods); i++ {
			gMap.createFood()
		}
	}
}

func (gMap *GameMap) publishStats() {
	for {
		time.Sleep(2 * time.Second)
		stats := gMap.GetStats()
		var validConnections []Connection
		for _, conn := range gMap.connections {
			if err := conn.WriteJSON(stats); err != nil {
				log.Println(err)
			} else {
				validConnections = append(validConnections, conn)
			}
		}
		gMap.connections = validConnections
	}
}

func (gMap *GameMap) Run() {
	counter := 0
	go gMap.publishStats()
	for {
		counter++
		var wg sync.WaitGroup
		if counter > 30 {
			counter = 0
			gMap.populateFood()
			gMap.populateBots()
		}
		for i := range gMap.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				player := &gMap.Players[i]
				if player.IsBot {
					player.makeMove(gMap)
				}
				player.updatePosition()
				player.passiveWeightLoss()
			}(i)
		}
		wg.Wait()
		gMap.removeEatableFood()
		gMap.removeEatablePlayers()
		time.Sleep(15 * time.Millisecond)
	}
}

func (gMap *GameMap) nearByFood(player *Player) []Food {
	var foods []Food

	for _, f := range gMap.Foods {
		distX := math.Abs(float64(player.X - f.X))
		distY := math.Abs(float64(player.Y - f.Y))
		if distX < windowSize && distY < windowSize {
			foods = append(foods, f)
		}
	}

	return foods
}

func (gMap GameMap) GetStats() ServerResponse {
	type StatsResponse struct {
		Nickname string `json:"nickname"`
		Weight   int    `json:"weight"`
	}
	players := gMap.Players
	numOfPlayers := len(players)
	playersInStats := statsNumber
	if numOfPlayers < statsNumber {
		playersInStats = numOfPlayers
	}
	var topPlayers []StatsResponse
	for i := 0; i < playersInStats; i++ {
		var maxIdx = i
		for j := i; j < numOfPlayers; j++ {
			if players[j].Weight > players[maxIdx].Weight {
				maxIdx = j
			}
		}
		players[i], players[maxIdx] = players[maxIdx], players[i]
	}
	for i := 0; i < playersInStats; i++ {
		topPlayers = append(topPlayers, StatsResponse{Nickname: players[i].Nickname, Weight: int(players[i].Weight)})
	}
	return ServerResponse{
		Event: "stats",
		Data:  topPlayers,
	}
}

func (gMap *GameMap) nearByPlayers(player *Player) []Player {
	var players []Player

	for _, p := range gMap.Players {
		if p.Uuid == player.Uuid {
			continue
		}
		distX := math.Abs(float64(player.X - p.X))
		distY := math.Abs(float64(player.Y - p.Y))
		windSize := windowSize / float64(player.Zoom)
		if distX < windSize && distY < windSize {
			players = append(players, p)
		}
	}

	return players
}

func (gMap *GameMap) removePlayerIndex(index int) {
	gMap.Players = append(gMap.Players[:index], gMap.Players[index+1:]...)
}

func (gMap *GameMap) removeFoodIndex(index int) {
	gMap.Foods = append(gMap.Foods[:index], gMap.Foods[index+1:]...)
}
func (gMap *GameMap) removeFoodUuid(uuid string) {
	for i, p := range gMap.Foods {
		if p.Uuid == uuid {
			gMap.Foods = append(gMap.Foods[:i], gMap.Foods[i+1:]...)
		}
	}
}

func (gMap *GameMap) RemovePlayerUUID(uuid string) {
	for i, p := range gMap.Players {
		if p.Uuid == uuid {
			gMap.Players = append(gMap.Players[:i], gMap.Players[i+1:]...)
		}
	}
}

func (gMap *GameMap) createFood() Food {
	x, y := getRandomCoordinate()
	food := Food{
		Uuid:   uuid4.New().String(),
		X:      x,
		Y:      y,
		Color:  getRandomColor(),
		Weight: foodWeight,
	}
	gMap.Foods = append(gMap.Foods, food)
	return food
}

func (gMap *GameMap) GetPlayer(uuid string) (*Player, error) {
	for i := range gMap.Players {
		if gMap.Players[i].Uuid == uuid {
			return &gMap.Players[i], nil
		}
	}
	return nil, errors.New("no player found")
}

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) SelfPlayer {
	x, y := getRandomCoordinate()
	player := Player{
		Uuid:         uuid4.New().String(),
		X:            x,
		Y:            y,
		Color:        getRandomColor(),
		Weight:       minWeight,
		Accelerating: false,
		Speed:        maxSpeed,
		Zoom:         1,
		Nickname:     nickname,
		IsBot:        isBot,
	}
	if len(gMap.Players) >= maxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.Players = append(gMap.Players, player)
	return player.getSelfPlayer()
}

func (gMap *GameMap) UpdatePlayer(uuid string, newX float32, newY float32) (*Player, error) {
	player, err := gMap.GetPlayer(uuid)
	if err != nil {
		return nil, err
	}
	player.updateDirection(newX, newY)
	return player, nil
}

func (gMap *GameMap) removeEatableFood() {
	for i := range gMap.Players {
		player := &gMap.Players[i]
		var filteredFoods []Food
		for k := range gMap.Foods {
			food := &gMap.Foods[k]
			if player.foodEatable(food) {
				player.eatFood(food)
			} else {
				filteredFoods = append(filteredFoods, *food)
			}
		}
		gMap.Foods = filteredFoods
	}
}

func (gMap *GameMap) removeEatablePlayers() {
	// take first player, compare it to every other player after it
	// get a new array of players that
	if len(gMap.Players) < 2 {
		return
	}
	var canBeEaten = make(map[int]bool, len(gMap.Players))
	for i := 0; i < len(gMap.Players); i++ {
		canBeEaten[i] = false
	}
	for i := range gMap.Players {
		p1 := &gMap.Players[i]
		for k := range gMap.Players[i+1:] {
			p2 := &gMap.Players[k]
			if p1.playerEatable(p2) {
				p1.eatPlayer(p2)
				canBeEaten[k] = true
			} else if p2.playerEatable(p1) {
				p2.eatPlayer(p1)
				canBeEaten[i] = true
			}
		}
	}
	//fmt.Println(eatenPlayers)
	var newPlayers []Player
	var eatenPlayers []Player
	for index, value := range canBeEaten {
		if value {
			if gMap.Players[index].IsBot {
				continue
			}
			eatenPlayers = append(eatenPlayers, gMap.Players[index])
		} else {
			newPlayers = append(newPlayers, gMap.Players[index])
		}
	}
	//for _, conn := range gMap.connections {
	//	if err := conn.Emit("removed", eatenPlayers); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	gMap.Players = newPlayers
}
