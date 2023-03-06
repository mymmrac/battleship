//go:generate protoc --go_out=. --go-grpc_out=. --experimental_allow_proto3_optional event_manager.proto

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/mymmrac/battleship/api"
)

const grpcPort = "42443"
const grpcAddr = "127.0.0.1"

func main() {
	rootCmd := &cobra.Command{
		Use:   "battleship",
		Short: "Battleship is two players sea battle game",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting...")

			game, err := NewGame()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Load game failed: %s\n", err)
				os.Exit(1)
			}

			if err = ebiten.RunGame(game); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Game crashed: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Bye!")
		},
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Battleship game server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting...")

			em := NewEventManagerServer()

			grpcServer := grpc.NewServer()
			api.RegisterEventManagerServer(grpcServer, em)

			listener, err := net.Listen("tcp", ":"+grpcPort)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Server crashed: %s\n", err)
				os.Exit(1)
			}

			// TODO: Graceful shutdown
			if err = grpcServer.Serve(listener); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Server crashed: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Bye!")
		},
	}

	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
