package events

import (
	"github.com/google/uuid"

	"github.com/mymmrac/battleship/server/api"
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
