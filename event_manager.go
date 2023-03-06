package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/mymmrac/battleship/api"
)

type ServerEventType int

const (
	_ ServerEventType = iota
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

func (e ServerEvent) EventType() GameEventType {
	return GameEventFromServer
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
		fmt.Printf("Event: %d, from %s, data: %v\n", event.Type, event.From, event.Data)

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
			games := make([]uuid.UUID, 0, len(e.games))
			for id, g := range e.games {
				if id == g.playerA.ID {
					games = append(games, g.playerA.ID)
				}
			}

			var data []byte
			data, err = json.Marshal(games)
			if err != nil {
				return err
			}

			err = stream.Send(ServerEvent{
				Type: ServerEventListGames,
				From: uuid.Nil,
				Data: data,
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

			game := e.games[gameID]
			e.games[event.From] = game

			player := Player{
				ID:     event.From,
				Events: make(chan ServerEvent),
			}
			game.playerB = player
			go player.HandleEvents(stream)

			gameEvent := GameEventSignal{Type: GameEventJoinedGame}
			var data []byte
			data, err = json.Marshal(gameEvent)
			if err != nil {
				return err
			}

			game.playerA.Events <- ServerEvent{
				Type: ServerEventGameEvent,
				From: uuid.Nil,
				Data: data,
			}
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
	playerID uuid.UUID
	stream   api.EventManager_EventsClient
}

func NewEventManagerClient(eventManager api.EventManagerClient) (*EventManagerClient, error) {
	stream, err := eventManager.Events(context.Background())
	if err != nil {
		return nil, err
	}

	return &EventManagerClient{
		playerID: uuid.New(),
		stream:   stream,
	}, nil
}

func (c *EventManagerClient) NewGame() error {
	return c.stream.Send(ServerEvent{
		Type: ServerEventNewGame,
		From: c.playerID,
		Data: nil,
	}.ToGRPC())
}

func (c *EventManagerClient) ListGames() ([]uuid.UUID, error) {
	err := c.stream.Send(ServerEvent{
		Type: ServerEventListGames,
		From: c.playerID,
		Data: nil,
	}.ToGRPC())
	if err != nil {
		return nil, err
	}

	grpcEvent, err := c.stream.Recv()
	if err != nil {
		return nil, err
	}

	event := ServerEventFromGRPC(grpcEvent)
	if event.Type != ServerEventListGames {
		return nil, errors.New("unexpected response event: " + strconv.Itoa(int(event.Type)))
	}

	var games []uuid.UUID
	if err = json.Unmarshal(event.Data, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (c *EventManagerClient) JoinGame(gameID uuid.UUID) error {
	return c.stream.Send(ServerEvent{
		Type: ServerEventJoinGame,
		From: c.playerID,
		Data: gameID[:],
	}.ToGRPC())
}

func (c *EventManagerClient) HandleGameEvents(events chan<- GameEvent) error {
	for {
		grpcEvent, err := c.stream.Recv()
		if err != nil {
			return err
		}

		serverEvent := ServerEventFromGRPC(grpcEvent)

		if serverEvent.Type != ServerEventGameEvent {
			return errors.New("unexpected event type: " + strconv.Itoa(int(serverEvent.Type)))
		}

		events <- serverEvent
	}
}
