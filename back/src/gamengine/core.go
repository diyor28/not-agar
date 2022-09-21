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
		Hub:        hub,
		Map:        gameMap,
		PlayersMap: make(map[*sockethub.Client]string),
	}
	hub.OnMessage(func(data []byte, client *sockethub.Client) {
		event := &GenericEvent{}
		if err := GenericSchema.Decode(data, event); err != nil {
			log.Println("GenericSchema.Decode(): ", err)
			return
		}
		switch event.Event {
		case "move":
			moveEvent := &MoveEvent{}
			if err := MoveSchema.Decode(data, moveEvent); err != nil {
				log.Println("MoveSchema.Decode(): ", err)
				return
			}
			engine.HandleMoveEvent(moveEvent, client)
		case "start":
			var startEvent = make(map[string]string)
			if err := StartSchema.Decode(data, &startEvent); err != nil {
				log.Println("StartSchema.Decode(): ", err)
				return
			}
			player, spikes := engine.Map.CreatePlayer(startEvent["nickname"], false)
			engine.PlayersMap[client] = player.Uuid
			startedEvent := &StartedEvent{
				Event: "started",
				Player: &StartedEventPlayer{
					player.X,
					player.Y,
					player.Weight,
					player.Color,
					player.Shell.Points,
				},
				Spikes: spikes,
			}
			if data, err := StartedSchema.Encode(startedEvent); err != nil {
				log.Println("StartedSchema.Encode(): ", err)
				return
			} else {
				client.Emit(data.Bytes())
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

func (eng *GameEngine) HandleMoveEvent(event *MoveEvent, client *sockethub.Client) {
	pl, err := eng.Map.Players.Update(eng.PlayersMap[client], event.NewX, event.NewY)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := MovedSchema.Encode(&MovedEvent{
		"moved",
		pl.X,
		pl.Y,
		pl.VelocityX,
		pl.VelocityY,
		pl.Weight,
		pl.Zoom,
		pl.Shell.Points,
	})
	if err != nil {
		log.Println(err)
		return
	}
	client.Emit(data.Bytes())
}

func (eng *GameEngine) SendPong(data []byte, client *sockethub.Client) {
	pingData := make(map[string]interface{})
	if err := PingPongSchema.Decode(data, &pingData); err != nil {
		log.Println("PingPongSchema.Decode(): ", err)
		return
	}
	pongEvent := map[string]interface{}{"event": "pong", "timestamp": pingData["timestamp"]}
	if data, err := PingPongSchema.Encode(&pongEvent); err != nil {
		log.Println("PingPongSchema.Encode(): ", err)
	} else {
		client.Emit(data.Bytes())
	}
}

func (eng *GameEngine) publishStats() {
	for range time.Tick(time.Duration(2000) * time.Millisecond) {
		stats := eng.Map.GetStats()
		statsEvent := &PlayerStatsEvent{
			Event:      "stats",
			TopPlayers: stats,
		}
		if data, err := PlayerStatsSchema.Encode(statsEvent); err != nil {
			log.Println(err)
		} else {
			eng.notifyAllPlayers(data.Bytes())
		}
	}
}

func (eng *GameEngine) notifyAllPlayers(data []byte) {
	for client, _ := range eng.PlayersMap {
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
	movedEvent := &MovedEvent{
		"moved",
		pl.X,
		pl.Y,
		pl.VelocityX,
		pl.VelocityY,
		pl.Weight,
		pl.Zoom,
		pl.Shell.Points,
	}
	if data, err := MovedSchema.Encode(movedEvent); err != nil {
		log.Println(err)
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	plrs := eng.Map.Players.Closest(pl, constants.NumPlayersResponse)
	plUpdateEvent := &PlayersUpdatedEvent{"players", plrs}
	if data, err := PlayersUpdatedSchema.Encode(plUpdateEvent); err != nil {
		log.Println(err)
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	food := eng.Map.Foods.Closest(pl, constants.NumFoodResponse)
	if data, err := FoodUpdateSchema.Encode(&FoodUpdatedEvent{Event: "food", Food: food}); err != nil {
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	return nil
}

func (eng *GameEngine) MakeMove(pl *players.Player) {
	food := eng.Map.Foods.Closest(pl, 1)
	closestFood := food[0]
	pl.UpdateDirection(closestFood.X, closestFood.Y)
}

func (eng *GameEngine) Run(framerate int) {
	delta := time.Duration(1000/framerate) * time.Millisecond
	go eng.publishStats()
	go eng.Hub.Run()
	go eng.Map.PopulateSpikes()
	go eng.publishAdminStats()
	for range time.Tick(delta) {
		eng.Map.PopulateFood()
		eng.Map.PopulateBots()
		var wg sync.WaitGroup
		for i := range eng.Map.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				pl := eng.Map.Players[i]
				if pl.IsBot {
					eng.MakeMove(pl)
				}
				pl.UpdatePosition(delta)
				pl.PassiveWeightLoss()
				if pl.IsBot {
					return
				}
				err := eng.notifyPlayer(pl)
				if err != nil {
					log.Println("error when calling eng.notifyPlayer(): ", err)
				}
			}(i)
		}
		wg.Wait()
		eng.Map.RemoveEatableFood()
		eng.removeDeadPlayers()
	}
}

func (eng *GameEngine) removeDeadPlayers() {
	ripEvent := map[string]interface{}{"event": "rip"}
	deadPlayers := eng.Map.RemoveDeadPlayers()
	for _, pl := range deadPlayers {
		client, err := eng.PlayerReverseLookUp(pl.Uuid)
		if err != nil {
			log.Println(err)
			continue
		}
		if data, err := GenericSchema.Encode(&ripEvent); err != nil {
			log.Println(err)
			continue
		} else {
			client.Emit(data.Bytes())
		}
	}
}
