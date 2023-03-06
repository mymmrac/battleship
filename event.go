package main

type GameEventType int

const (
	_ GameEventType = iota
	GameEventFromServer
	GameEventNewGameStarted
	GameEventNewGameStartFailed
	GameEventJoinedGame
	GameEventJoinGameFailed
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
