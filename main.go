package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
