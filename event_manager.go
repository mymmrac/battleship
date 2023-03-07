package main

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"

	"github.com/mymmrac/battleship/events"
	"github.com/mymmrac/battleship/server/api"
)

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
	return c.stream.Send(events.ServerEvent{
		Type: events.ServerEventNewGame,
		From: c.playerID,
		Data: nil,
	}.ToGRPC())
}

func (c *EventManagerClient) ListGames() ([]uuid.UUID, error) {
	err := c.stream.Send(events.ServerEvent{
		Type: events.ServerEventListGames,
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

	event := events.ServerEventFromGRPC(grpcEvent)
	if event.Type != events.ServerEventListGames {
		return nil, errors.New("unexpected response event: " + strconv.Itoa(int(event.Type)))
	}

	var games []uuid.UUID
	if err = json.Unmarshal(event.Data, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (c *EventManagerClient) JoinGame(gameID uuid.UUID) error {
	return c.stream.Send(events.ServerEvent{
		Type: events.ServerEventJoinGame,
		From: c.playerID,
		Data: gameID[:],
	}.ToGRPC())
}

func (c *EventManagerClient) HandleGameEvents(gameEvents chan<- events.GameEvent) error {
	for {
		grpcEvent, err := c.stream.Recv()
		if err != nil {
			return err
		}

		serverEvent := events.ServerEventFromGRPC(grpcEvent)

		if serverEvent.Type != events.ServerEventGameEvent {
			return errors.New("unexpected event type: " + strconv.Itoa(int(serverEvent.Type)))
		}

		gameEvents <- serverEvent
	}
}

func (c *EventManagerClient) SendGameEvent(event events.GameEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return c.stream.Send(events.ServerEvent{
		Type: events.ServerEventGameEvent,
		From: c.playerID,
		Data: data,
	}.ToGRPC())
}
