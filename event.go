package main

type EventType int

const (
	EventNone EventType = iota
	EventNewGameStarted
	EventNewGameStartFailed
	EventJoinedGame
	EventJoinGameFailed
)

type Event interface {
	EventType() EventType
}

type EventError struct {
	eventType EventType
	err       error
}

func NewEventError(eventType EventType, err error) EventError {
	return EventError{
		eventType: eventType,
		err:       err,
	}
}

func (e EventError) EventType() EventType {
	return e.eventType
}

type EventSignal struct {
	eventType EventType
}

func NewEventSignal(eventType EventType) EventSignal {
	return EventSignal{
		eventType: eventType,
	}
}

func (e EventSignal) EventType() EventType {
	return e.eventType
}
