package gamengine

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	_map "github.com/diyor28/not-agar/src/gamengine/map"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/sockethub"
	"log"
	"sync"
	"time"
)

type GameEngine struct {
	Hub        *sockethub.Hub
	Map        *_map.Map
	PlayersMap map[*sockethub.Client]string
}

func NewGameMap() *GameEngine {
	hub := sockethub.NewHub()
	gameMap := _map.New()
	engine := GameEngine{
		Hub: hub,
		Map: gameMap,
	}
	hub.OnMessage(func(data []byte, client *sockethub.Client) {
		event := &GenericEvent{}
		if err := GenericSchema.Decode(data, event); err != nil {
			log.Println("GenericSchema.Decode(): ", err)
			return
		}
		switch event.Event {
		case "move":
			var moveEvent *MoveEvent
			if err := MovedSchema.Decode(data, moveEvent); err != nil {
				log.Println("MovedSchema.Decode(): ", err)
				return
			}
			engine.HandleMoveEvent(moveEvent, client)
		case "start":
			var startEvent map[string]string
			if err := StartSchema.Decode(data, &startEvent); err != nil {
				log.Println("StartSchema.Decode(): ", err)
				return
			}
			player, spikes := engine.Map.CreatePlayer(startEvent["event"], false)
			startedEvent := &StartedEvent{
				Event:  "started",
				Player: player,
				Spikes: spikes,
			}
			if bytes, err := StartedSchema.Encode(startedEvent); err != nil {
				log.Println("StartedSchema.Encode(): ", err)
				return
			} else {
				client.Emit(bytes)
				client.Join(fmt.Sprintf("player/%s", player.Uuid))
			}
		case "ping":
			engine.SendPong(data, client)
		}
	})
	return &engine
}

func (eng *GameEngine) PlayerReverseLookUp(uuid string) (*sockethub.Client, error) {
	for k, pUuid := range eng.PlayersMap {
		if pUuid == uuid {
			return k, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no players for uuid %s found", uuid))
}

func (eng *GameEngine) HandleMoveEvent(data *MoveEvent, client *sockethub.Client) {
	pl, err := eng.Map.Players.Update(eng.PlayersMap[client], data.NewX, data.NewY)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := MovedSchema.Encode(&MovedEvent{"moved", pl.X, pl.Y, pl.Weight, pl.Zoom})
	if err != nil {
		log.Println(err)
		return
	}
	client.Emit(bytes)
}

func (eng *GameEngine) SendPong(data []byte, client *sockethub.Client) {
	client.Emit(data)
}

func (eng *GameEngine) publishStats() {
	for {
		time.Sleep(2 * time.Second)
		stats := eng.Map.GetStats()
		statsEvent := &PlayerStatsEvent{
			Event:      "stats",
			TopPlayers: stats,
		}
		if bytes, err := PlayerStatsSchema.Encode(statsEvent); err != nil {
			log.Println(err)
		} else {
			eng.notifyAllPlayers(bytes)
		}
	}
}

func (eng *GameEngine) notifyAllPlayers(data []byte) {
	for _, p := range eng.Map.Players.Real() {
		client, err := eng.PlayerReverseLookUp(p.Uuid)
		if err != nil {
			log.Println(err)
			continue
		}
		client.Emit(data)
	}
}

func (eng *GameEngine) publishAdminStats() {
	for {
		var result = make(map[string]interface{})
		botsCount := eng.Map.Players.BotsCount()
		stats := eng.Map.GetStats()

		result["botsCount"] = botsCount
		result["playersCount"] = len(eng.Map.Players) - botsCount
		result["topPlayers"] = stats
		// TODO: fix later
		//eng.Hub.Emit("stats", result, "admin")
		time.Sleep(2 * time.Second)
	}
}

func (eng *GameEngine) notifyPlayer(pl *players.Player) error {
	client, err := eng.PlayerReverseLookUp(pl.Uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	movedEvent := &MovedEvent{"moved", pl.X, pl.Y, pl.Weight, pl.Zoom}
	if bytes, err := MovedSchema.Encode(movedEvent); err != nil {
		log.Println(err)
		return err
	} else {
		eng.Hub.Emit(bytes, pl.Uuid)
	}
	plrs := eng.Map.Players.Closest(pl, constants.NumPlayersResponse)
	plUpdateEvent := &PlayersUpdatedEvent{"pUpdated", plrs}
	if bytes, err := PlayersUpdatedSchema.Encode(plUpdateEvent); err != nil {
		log.Println(err)
		return err
	} else {
		eng.Hub.Emit(bytes, pl.Uuid)
	}
	foods := eng.Map.Foods.Closest(pl, constants.NumFoodResponse)
	if bytes, err := FoodUpdateSchema.Encode(&FoodUpdateEvent{Event: "foodUpdated", Food: foods}); err != nil {
		return err
	} else {
		client.Emit(bytes)
	}
	return nil
}

func (eng *GameEngine) MakeMove(pl *players.Player) {
	foods := eng.Map.Foods.Closest(pl, 1)
	closestFood := foods[0]
	pl.UpdateDirection(closestFood.X, closestFood.Y)
}

func (eng *GameEngine) Run() {
	counter := 0
	go eng.publishStats()
	go eng.Hub.Run()
	go eng.Map.PopulateSpikes()
	go eng.publishAdminStats()
	for {
		counter++
		var wg sync.WaitGroup
		if counter > 30 {
			counter = 0
			eng.Map.PopulateFood()
			eng.Map.PopulateBots()
		}
		for i := range eng.Map.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				pl := eng.Map.Players[i]
				if pl.IsBot {
					eng.MakeMove(pl)
				}
				pl.UpdatePosition()
				pl.PassiveWeightLoss()
				if pl.IsBot {
					return
				}
				if counter%2 == 0 {
					err := eng.notifyPlayer(pl)
					if err != nil {
						log.Println("error when calling eng.notifyPlayer(): ", err)
					}
				}
			}(i)
		}
		wg.Wait()
		eng.Map.RemoveEatableFood()
		eng.removeDeadPlayers()
		time.Sleep(15 * time.Millisecond)
		//fmt.Println("running loop")
	}
}

func (eng *GameEngine) removeDeadPlayers() {
	// take first players, compare it to every other players after it
	// get a new array of players that
	//for _, p := range eatenPlayers {
	//	eng.Hub.Emit("rip", "", p.Uuid)
	//}
	//eng.Players = newPlayers
}
