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
			game := NewGame()

			if err := ebiten.RunGame(game); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Game crashed: %s\n", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
