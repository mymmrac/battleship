package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/mymmrac/battleship/server"
	"github.com/mymmrac/battleship/server/api"
)

const (
	DefaultGRPCPort    = "42284"
	defaultStopTimeout = 4 * time.Second
)

func BattleshipServerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("port", "p", DefaultGRPCPort, "Battleship server port used to start server")
	cmd.Flags().DurationP("timeout", "t", defaultStopTimeout, "Battleship server timeout duration")
}

func BattleshipServerRunE(cmd *cobra.Command, _ []string) error {
	fmt.Println("Starting...")

	serverPort, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}

	stopTimeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	em := server.NewEventManagerServer()

	grpcServer := grpc.NewServer()
	api.RegisterEventManagerServer(grpcServer, em)

	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		return fmt.Errorf("server crashed: %w", err)
	}

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Server crashed: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	fmt.Println("Listening on port", serverPort)
	<-quit
	fmt.Println("Stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()

	done := make(chan struct{}, 1)
	go func() {
		grpcServer.GracefulStop()

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Printf("Stopping failed: %s\n", ctx.Err())
	case <-done:
		fmt.Println("Bye!")
	}

	return nil
}
