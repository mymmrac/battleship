package main

type GameEventType int

const (
	GameEventNone GameEventType = iota
	GameEventNewGameStarted
	GameEventNewGameStartFailed
	GameEventJoinedGame
	GameEventJoinGameFailed
)

type GameEvent interface {
	EventType() GameEventType
}

type GameEventError struct {
	eventType GameEventType
	err       error
}

func NewGameEventError(eventType GameEventType, err error) GameEventError {
	return GameEventError{
		eventType: eventType,
		err:       err,
	}
}

func (e GameEventError) EventType() GameEventType {
	return e.eventType
}

type GameEventSignal struct {
	eventType GameEventType
}

func NewGameEventSignal(eventType GameEventType) GameEventSignal {
	return GameEventSignal{
		eventType: eventType,
	}
}

func (e GameEventSignal) EventType() GameEventType {
	return e.eventType
}
