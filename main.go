package main

import (
	"tanks3d/game3d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize window
	rl.InitWindow(1024, 768, "3D Tanks - World of Tanks Style")
	defer rl.CloseWindow()
	
	rl.SetTargetFPS(60)
	
	// Initialize game
	game := game3d.NewGame()
	
	// Game loop
	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}