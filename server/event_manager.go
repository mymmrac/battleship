//go:generate protoc --go_out=. --go-grpc_out=. --experimental_allow_proto3_optional event_manager.proto

package server

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/mymmrac/battleship/events"
	"github.com/mymmrac/battleship/server/api"
)

type Player struct {
	ID     uuid.UUID
	Events chan events.ServerEvent
}

func (p *Player) HandleEvents(stream api.EventManager_EventsServer) {
	for event := range p.Events {
		err := stream.Send(event.ToGRPC())
		if err != nil {
			fmt.Println(err) // FIXME
		}
	}
}

type MultiplayerGame struct {
	playerA *Player
	playerB *Player
}

type EventManagerServer struct {
	api.UnimplementedEventManagerServer

	games map[uuid.UUID]*MultiplayerGame
}

func NewEventManagerServer() *EventManagerServer {
	return &EventManagerServer{
		games: map[uuid.UUID]*MultiplayerGame{},
	}
}

func (e *EventManagerServer) Events(stream api.EventManager_EventsServer) error {
	for {
		grpcEvent, err := stream.Recv()
		if err != nil {
			return err
		}

		event := events.ServerEventFromGRPC(grpcEvent)
		fmt.Printf("Event: %d, from %s, data: %v\n", event.Type, event.From, event.Data)

		switch event.Type {
		case events.ServerEventNewGame:
			player := &Player{
				ID:     event.From,
				Events: make(chan events.ServerEvent),
			}
			e.games[event.From] = &MultiplayerGame{
				playerA: player,
			}
			go player.HandleEvents(stream)
		case events.ServerEventListGames:
			games := make([]uuid.UUID, 0, len(e.games))
			for id, g := range e.games {
				if g.playerB == nil && id == g.playerA.ID {
					games = append(games, g.playerA.ID)
				}
			}

			var data []byte
			data, err = json.Marshal(games)
			if err != nil {
				return err
			}

			err = stream.Send(events.ServerEvent{
				Type: events.ServerEventListGames,
				From: uuid.Nil,
				Data: data,
			}.ToGRPC())
			if err != nil {
				return err
			}
		case events.ServerEventJoinGame:
			var gameID uuid.UUID
			gameID, err = uuid.FromBytes(event.Data)
			if err != nil {
				return err
			}

			game := e.games[gameID]
			e.games[event.From] = game

			player := &Player{
				ID:     event.From,
				Events: make(chan events.ServerEvent),
			}
			game.playerB = player
			go player.HandleEvents(stream)

			gameEvent := events.GameEventSignal{Type: events.GameEventJoinedGame}
			var data []byte
			data, err = json.Marshal(gameEvent)
			if err != nil {
				return err
			}

			game.playerA.Events <- events.ServerEvent{
				Type: events.ServerEventGameEvent,
				From: uuid.Nil,
				Data: data,
			}
		case events.ServerEventGameEvent:
			game := e.games[event.From]

			if event.From == game.playerA.ID {
				game.playerB.Events <- event
			} else {
				game.playerA.Events <- event
			}
		}
	}
}
