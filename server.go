package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mymmrac/battleship/api"
)

type EventManagerServer struct {
	api.UnimplementedEventManagerServer
}

type BattleshipConnector interface {
	StartNewGame() error
	StopGame() error
	WaitForConnection() error

	JoinGame() error
	ExitGame() error

	PlayerReady() error
	PlayerNotReady() error
	Shoot(x, y int) (cellKind, error)
}

type Connector struct {
	eventManagerClient api.EventManagerClient
	grpcConn           *grpc.ClientConn
	eventClientStream  api.EventManager_EventsClient

	eventManagerServer *EventManagerServer
	grpcServer         *grpc.Server
}

func NewConnector() *Connector {
	return &Connector{}
}

const grpcPort = "42443"
const grpcAddr = "127.0.0.1"

func (c *Connector) StartNewGame() error {
	c.grpcServer = grpc.NewServer()
	api.RegisterEventManagerServer(c.grpcServer, c.eventManagerServer)

	l, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return err
	}

	go func() {
		if err = c.grpcServer.Serve(l); err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (c *Connector) StopGame() error {
	c.grpcServer.GracefulStop()
	return nil
}

func (c *Connector) WaitForConnection() error {
	return nil
}

func (c *Connector) JoinGame() error {
	var err error
	c.grpcConn, err = grpc.Dial(grpcAddr+":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.eventManagerClient = api.NewEventManagerClient(c.grpcConn)

	c.eventClientStream, err = c.eventManagerClient.Events(context.Background())
	if err != nil {
		return err
	}

	err = c.eventClientStream.Send(&api.Event{
		EventType: int32(EventJoinedGame),
		Data:      nil,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Connector) ExitGame() error {
	return c.grpcConn.Close()
}

func (c *Connector) PlayerReady() error {
	// TODO implement me
	panic("implement me")
}

func (c *Connector) PlayerNotReady() error {
	// TODO implement me
	panic("implement me")
}

func (c *Connector) Shoot(x, y int) (cellKind, error) {
	// TODO implement me
	panic("implement me")
}
