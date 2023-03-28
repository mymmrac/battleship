package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"

	"github.com/mymmrac/battleship/cmd"
	"github.com/mymmrac/battleship/cmd/server"
)

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
	rootCmd.Flags().StringVarP(&serverPort, "port", "p", server.DefaultGRPCPort, "Battleship server port used to connect")

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Battleship game server",
		RunE:  server.BattleshipServerRunE,
	}

	server.BattleshipServerFlags(serverCmd)

	rootCmd.AddCommand(serverCmd)

	cmd.WalkCmd(rootCmd, cmd.UpdateHelp)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
