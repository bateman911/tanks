package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"tanks/game"
)

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Tanks - World of Tanks Style Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}