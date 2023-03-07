package main

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/mymmrac/battleship/server"
	"github.com/mymmrac/battleship/server/api"
)

const defaultGRPCPort = "42284"

func main() {
	var serverPort string

	rootCmd := &cobra.Command{
		Use:   "battleship-server",
		Short: "Battleship game server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting...")

			em := server.NewEventManagerServer()

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

	rootCmd.Flags().StringVarP(&serverPort, "port", "p", defaultGRPCPort, "Battleship server port used to start server")

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
