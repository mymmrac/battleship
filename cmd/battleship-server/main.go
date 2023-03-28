package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mymmrac/battleship/cmd"
	"github.com/mymmrac/battleship/cmd/server"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "battleship-server",
		Short: "Battleship game server",
		RunE:  server.BattleshipServerRunE,
	}

	server.BattleshipServerFlags(rootCmd)
	cmd.WalkCmd(rootCmd, cmd.UpdateHelp)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
