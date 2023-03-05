package main

import "fmt"

type BattleshipServer interface {
	StartNewGame() error
	JoinGame() error
	PlayerReady() error
	PlayerNotReady() error
	Shoot(x, y int) (cellKind, error)
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) StartNewGame() error {
	fmt.Println("Start new game")
	return nil
}

func (s *Server) JoinGame() error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) PlayerReady() error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) PlayerNotReady() error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) Shoot(x, y int) (cellKind, error) {
	// TODO implement me
	panic("implement me")
}
