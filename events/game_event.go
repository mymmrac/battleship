package events

import "github.com/mymmrac/battleship/data"

type GameEventType int

const (
	_ GameEventType = iota
	GameEventFromServer
	GameEventNewGameStarted
	GameEventNewGameStartFailed
	GameEventJoinedGame
	GameEventJoinGameFailed
	GameEventPlayerReady
	GameEventPlayerNotReady
	GameEventShoot
	GameEventMiss
	GameEventHit
	GameEventDestroyed
	GameEventGameEnded
)

type GameEvent interface {
	EventType() GameEventType
}

type GameEventError struct {
	Type GameEventType
	Err  error
}

func NewGameEventError(eventType GameEventType, err error) GameEventError {
	return GameEventError{
		Type: eventType,
		Err:  err,
	}
}

func (e GameEventError) EventType() GameEventType {
	return e.Type
}

type GameEventSignal struct {
	Type GameEventType
}

func NewGameEventSignal(eventType GameEventType) GameEventSignal {
	return GameEventSignal{
		Type: eventType,
	}
}

func (e GameEventSignal) EventType() GameEventType {
	return e.Type
}

type GameEventCoord struct {
	Type GameEventType
	Pos  data.Point[int]
}

func NewGameEventCoord(pos data.Point[int]) GameEventCoord {
	return GameEventCoord{
		Type: GameEventShoot,
		Pos:  pos,
	}
}

func (e GameEventCoord) EventType() GameEventType {
	return e.Type
}
