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

const defaultGRPCPort = "42284"

func main() {
	var serverAddr string
	var serverPort string

	rootCmd := &cobra.Command{
		Use:   "battleship",
		Short: "Battleship is two players sea battle game",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting...")

			game, err := NewGame(serverAddr, serverPort)
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

	rootCmd.Flags().StringVarP(&serverAddr, "address", "a", "127.0.0.1", "Battleship server address used to connect")
	rootCmd.Flags().StringVarP(&serverPort, "port", "p", defaultGRPCPort, "Battleship server port used to connect")

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Battleship game server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting...")

			em := NewEventManagerServer()

			grpcServer := grpc.NewServer()
			api.RegisterEventManagerServer(grpcServer, em)

			listener, err := net.Listen("tcp", ":"+serverPort)
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

	serverCmd.Flags().StringVarP(&serverPort, "port", "p", defaultGRPCPort, "Battleship server port used to start server")

	rootCmd.AddCommand(serverCmd)

	walkCmd(rootCmd, updateHelp)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func walkCmd(cmd *cobra.Command, f func(*cobra.Command)) {
	f(cmd)
	for _, childCmd := range cmd.Commands() {
		walkCmd(childCmd, f)
	}
}

func updateHelp(cmd *cobra.Command) {
	cmd.InitDefaultHelpFlag()
	f := cmd.Flags().Lookup("help")
	if f != nil {
		if cmd.Name() != "" {
			f.Usage = "Help for " + cmd.Name()
		} else {
			f.Usage = "Help for this command"
		}
	}
}
