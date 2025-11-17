/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/ABC10946/minesweeper/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

var (
	width  int
	height int
	cell   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "minesweeper",
	Short: "minesweeper command",
	Long:  `minesweeper command`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		game_ := game.Game{}
		game_.CellSize = cell
		game_.GameMode = game.MenuWindow

		game_.MS.Init(width, height)
		game_.MS.SummonBomb()
		game_.MS.CountBomb()
		ebiten.SetWindowSize(game_.CellSize*width, game_.CellSize*height)
		ebiten.SetWindowTitle("MINESWEEPER")
		if err := ebiten.RunGame(&game_); err != nil {
			log.Fatal(err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.minesweeper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntVar(&width, "width", 10, "width of minesweeper field")
	rootCmd.Flags().IntVar(&height, "height", 10, "height of minesweeper field")
	rootCmd.Flags().IntVar(&cell, "size", 70, "size of cell")

}
