package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/mymmrac/battleship/api"
)

type ServerEventType int

const (
	ServerEventNone ServerEventType = iota
	ServerEventNewGame
	ServerEventListGames
	ServerEventJoinGame
	ServerEventGameEvent
)

type ServerEvent struct {
	Type ServerEventType
	From uuid.UUID
	Data []byte
}

func ServerEventFromGRPC(grpcEvent *api.Event) ServerEvent {
	return ServerEvent{
		Type: ServerEventType(grpcEvent.Type),
		From: uuid.Must(uuid.FromBytes(grpcEvent.From.Value)),
		Data: grpcEvent.Data,
	}
}

func (e ServerEvent) ToGRPC() *api.Event {
	return &api.Event{
		Type: int32(e.Type),
		From: &api.UUID{Value: e.From[:]},
		Data: e.Data,
	}
}

type Player struct {
	ID     uuid.UUID
	Events chan ServerEvent
}

func (p Player) HandleEvents(stream api.EventManager_EventsServer) {
	for event := range p.Events {
		err := stream.Send(event.ToGRPC())
		if err != nil {
			fmt.Println(err)
		}
	}
}

type MultiplayerGame struct {
	playerA Player
	playerB Player
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

		event := ServerEventFromGRPC(grpcEvent)

		switch event.Type {
		case ServerEventNewGame:
			player := Player{
				ID:     event.From,
				Events: make(chan ServerEvent),
			}
			e.games[event.From] = &MultiplayerGame{
				playerA: player,
			}
			go player.HandleEvents(stream)
		case ServerEventListGames:
			buf := &bytes.Buffer{}

			games := make([]uuid.UUID, 0, len(e.games))
			for id, g := range e.games {
				if id == g.playerA.ID {
					games = append(games, g.playerA.ID)
				}
			}

			err = binary.Write(buf, binary.BigEndian, games)
			if err != nil {
				return err
			}

			err = stream.Send(ServerEvent{
				Type: ServerEventListGames,
				From: uuid.Nil,
				Data: buf.Bytes(),
			}.ToGRPC())
			if err != nil {
				return err
			}
		case ServerEventJoinGame:
			var gameID uuid.UUID
			gameID, err = uuid.FromBytes(event.Data)
			if err != nil {
				return err
			}

			e.games[event.From] = e.games[gameID]

			player := Player{
				ID:     event.From,
				Events: make(chan ServerEvent),
			}
			e.games[gameID].playerB = player
			go player.HandleEvents(stream)
		case ServerEventGameEvent:
			game := e.games[event.From]

			if event.From == game.playerA.ID {
				game.playerB.Events <- event
			} else {
				game.playerA.Events <- event
			}
		}
	}
}

type EventManagerClient struct {
	grpcConn     *grpc.ClientConn
	eventManager api.EventManagerClient
}
