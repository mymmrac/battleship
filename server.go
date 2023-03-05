package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"

	"github.com/mymmrac/battleship/api"
)

type EventManagerServer struct {
	api.UnimplementedEventManagerServer
}

type BattleshipConnector interface {
	StartNewGame() error
	StopGame() error

	JoinGame() error
	ExitGame() error

	PlayerReady() error
	PlayerNotReady() error
	Shoot(x, y int) (cellKind, error)
}

type Connector struct {
	eventManagerClient api.EventManagerClient
	grpcConn           *grpc.ClientConn

	eventManagerServer *EventManagerServer
	grpcServer         *grpc.Server
}

func NewConnector() *Connector {
	return &Connector{}
}

const grpcPort = "42443"
const grpcAddr = "127.0.0.1"

func (c *Connector) StartNewGame() error {
	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
	c.grpcServer = grpc.NewServer(grpc.Creds(altsTC))

	c.eventManagerServer = &EventManagerServer{}
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

func (c *Connector) JoinGame() error {
	var err error
	altsTC := alts.NewClientCreds(alts.DefaultClientOptions())
	c.grpcConn, err = grpc.Dial(grpcAddr+":"+grpcPort, grpc.WithTransportCredentials(altsTC))
	if err != nil {
		return err
	}

	c.eventManagerClient = api.NewEventManagerClient(c.grpcConn)
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
