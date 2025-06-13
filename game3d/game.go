package game3d

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MapSize = 100.0
)

type Game struct {
	camera    rl.Camera3D
	player    *Tank
	enemies   []*Tank
	bullets   []*Bullet
	terrain   *Terrain
	gameTime  int
}

func NewGame() *Game {
	// Initialize camera
	camera := rl.Camera3D{
		Position:   rl.NewVector3(10, 15, 10),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}

	// Create player tank
	player := NewTank(rl.NewVector3(0, 0, 0), true)

	// Create enemy tanks
	enemies := []*Tank{
		NewTank(rl.NewVector3(20, 0, 20), false),
		NewTank(rl.NewVector3(-20, 0, 20), false),
		NewTank(rl.NewVector3(30, 0, -10), false),
	}

	// Create terrain
	terrain := NewTerrain()

	return &Game{
		camera:  camera,
		player:  player,
		enemies: enemies,
		bullets: make([]*Bullet, 0),
		terrain: terrain,
	}
}

func (g *Game) Update() {
	g.gameTime++

	// Update player
	g.player.Update()
	g.handleInput()

	// Update enemies with AI
	for _, enemy := range g.enemies {
		if enemy.Health > 0 {
			g.updateEnemyAI(enemy)
			enemy.Update()
		}
	}

	// Update bullets
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		bullet.Update()

		// Remove bullets that are out of bounds or expired
		if bullet.Position.X < -MapSize || bullet.Position.X > MapSize ||
			bullet.Position.Z < -MapSize || bullet.Position.Z > MapSize ||
			bullet.LifeTime <= 0 {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			continue
		}

		// Check bullet collisions
		g.checkBulletCollisions(bullet, i)
	}

	// Update camera to follow player
	g.updateCamera()
}

func (g *Game) handleInput() {
	if g.player.Health <= 0 {
		return
	}

	// Tank movement
	if rl.IsKeyDown(rl.KeyW) {
		g.player.MoveForward()
	}
	if rl.IsKeyDown(rl.KeyS) {
		g.player.MoveBackward()
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.player.TurnLeft()
	}
	if rl.IsKeyDown(rl.KeyD) {
		g.player.TurnRight()
	}

	// Turret rotation
	if rl.IsKeyDown(rl.KeyLeft) {
		g.player.TurretLeft()
	}
	if rl.IsKeyDown(rl.KeyRight) {
		g.player.TurretRight()
	}

	// Shooting
	if rl.IsKeyPressed(rl.KeySpace) {
		if bullet := g.player.Shoot(); bullet != nil {
			g.bullets = append(g.bullets, bullet)
		}
	}
}

func (g *Game) updateEnemyAI(enemy *Tank) {
	// Simple AI: move towards player and shoot occasionally
	dx := g.player.Position.X - enemy.Position.X
	dz := g.player.Position.Z - enemy.Position.Z
	distance := math.Sqrt(float64(dx*dx + dz*dz))

	if distance > 5 {
		// Calculate target angle
		targetAngle := math.Atan2(float64(dx), float64(dz))
		angleDiff := targetAngle - float64(enemy.Rotation)

		// Normalize angle difference
		for angleDiff > math.Pi {
			angleDiff -= 2 * math.Pi
		}
		for angleDiff < -math.Pi {
			angleDiff += 2 * math.Pi
		}

		// Turn towards player
		if math.Abs(angleDiff) > 0.1 {
			if angleDiff > 0 {
				enemy.TurnRight()
			} else {
				enemy.TurnLeft()
			}
		} else {
			enemy.MoveForward()
		}
	}

	// Shoot occasionally
	if g.gameTime%180 == 0 && distance < 30 {
		if bullet := enemy.Shoot(); bullet != nil {
			g.bullets = append(g.bullets, bullet)
		}
	}
}

func (g *Game) checkBulletCollisions(bullet *Bullet, bulletIndex int) {
	// Check collision with player
	if !bullet.FromPlayer && g.player.Health > 0 {
		if g.checkCollision(bullet.Position, g.player.Position, 2.0) {
			g.player.TakeDamage(25)
			g.bullets = append(g.bullets[:bulletIndex], g.bullets[bulletIndex+1:]...)
			return
		}
	}

	// Check collision with enemies
	if bullet.FromPlayer {
		for _, enemy := range g.enemies {
			if enemy.Health > 0 && g.checkCollision(bullet.Position, enemy.Position, 2.0) {
				enemy.TakeDamage(25)
				g.bullets = append(g.bullets[:bulletIndex], g.bullets[bulletIndex+1:]...)
				return
			}
		}
	}
}

func (g *Game) checkCollision(pos1, pos2 rl.Vector3, radius float32) bool {
	dx := pos1.X - pos2.X
	dz := pos1.Z - pos2.Z
	distance := math.Sqrt(float64(dx*dx + dz*dz))
	return distance < float64(radius)
}

func (g *Game) updateCamera() {
	// Third-person camera following the player
	cameraDistance := float32(15)
	cameraHeight := float32(8)

	// Calculate camera position behind the tank
	cameraX := g.player.Position.X - float32(math.Sin(float64(g.player.Rotation)))*cameraDistance
	cameraZ := g.player.Position.Z - float32(math.Cos(float64(g.player.Rotation)))*cameraDistance

	g.camera.Position = rl.NewVector3(cameraX, cameraHeight, cameraZ)
	g.camera.Target = rl.NewVector3(g.player.Position.X, g.player.Position.Y+1, g.player.Position.Z)
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.SkyBlue)

	rl.BeginMode3D(g.camera)

	// Draw terrain
	g.terrain.Draw()

	// Draw tanks
	g.player.Draw()
	for _, enemy := range g.enemies {
		if enemy.Health > 0 {
			enemy.Draw()
		}
	}

	// Draw bullets
	for _, bullet := range g.bullets {
		bullet.Draw()
	}

	// Draw grid for reference
	rl.DrawGrid(100, 1.0)

	rl.EndMode3D()

	// Draw UI
	g.drawUI()

	rl.EndDrawing()
}

func (g *Game) drawUI() {
	// Health bar
	healthBarWidth := int32(200)
	healthBarHeight := int32(20)
	healthPercentage := float32(g.player.Health) / 100.0

	// Background
	rl.DrawRectangle(10, 10, healthBarWidth, healthBarHeight, rl.Gray)

	// Health
	healthColor := rl.Red
	if healthPercentage > 0.5 {
		healthColor = rl.Green
	} else if healthPercentage > 0.25 {
		healthColor = rl.Yellow
	}

	rl.DrawRectangle(10, 10, int32(float32(healthBarWidth)*healthPercentage), healthBarHeight, healthColor)

	// Health text
	rl.DrawText("Health", 10, 35, 20, rl.Black)

	// Game status
	if g.player.Health <= 0 {
		rl.DrawText("GAME OVER - Press ESC to exit", 300, 350, 30, rl.Red)
	}

	// Enemy count
	aliveEnemies := 0
	for _, enemy := range g.enemies {
		if enemy.Health > 0 {
			aliveEnemies++
		}
	}
	rl.DrawText(rl.TextFormat("Enemies: %d", aliveEnemies), 10, 60, 20, rl.Black)

	// Controls
	rl.DrawText("WASD - Move, Arrow Keys - Turret, Space - Shoot", 10, 720, 20, rl.DarkGray)
}