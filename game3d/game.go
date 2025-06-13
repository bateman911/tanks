package game3d

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MapSize = 100.0
)

type Game struct {
	camera         rl.Camera3D
	player         *Tank
	enemies        []*Tank
	bullets        []*Bullet
	terrain        *Terrain
	gameTime       int
	mouseAiming    bool
	aimingCircle   AimingCircle
}

type AimingCircle struct {
	CurrentRadius float32
	MinRadius     float32
	MaxRadius     float32
	ShrinkSpeed   float32
	ExpandSpeed   float32
	IsAiming      bool
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

	// Initialize aiming circle
	aimingCircle := AimingCircle{
		CurrentRadius: 50.0,
		MinRadius:     20.0,
		MaxRadius:     80.0,
		ShrinkSpeed:   1.5,
		ExpandSpeed:   2.0,
		IsAiming:      false,
	}

	return &Game{
		camera:       camera,
		player:       player,
		enemies:      enemies,
		bullets:      make([]*Bullet, 0),
		terrain:      terrain,
		mouseAiming:  true,
		aimingCircle: aimingCircle,
	}
}

func (g *Game) Update() {
	g.gameTime++

	// Update player
	g.player.Update()
	g.handleInput()

	// Update aiming system
	g.updateAiming()

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
		g.aimingCircle.IsAiming = false // Движение ухудшает точность
	}
	if rl.IsKeyDown(rl.KeyS) {
		g.player.MoveBackward()
		g.aimingCircle.IsAiming = false
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.player.TurnLeft()
		g.aimingCircle.IsAiming = false
	}
	if rl.IsKeyDown(rl.KeyD) {
		g.player.TurnRight()
		g.aimingCircle.IsAiming = false
	}

	// Mouse aiming
	if g.mouseAiming {
		g.handleMouseAiming()
	} else {
		// Keyboard turret rotation (fallback)
		if rl.IsKeyDown(rl.KeyLeft) {
			g.player.TurretLeft()
			g.aimingCircle.IsAiming = false
		}
		if rl.IsKeyDown(rl.KeyRight) {
			g.player.TurretRight()
			g.aimingCircle.IsAiming = false
		}
	}

	// Toggle aiming mode
	if rl.IsKeyPressed(rl.KeyTab) {
		g.mouseAiming = !g.mouseAiming
		if g.mouseAiming {
			rl.DisableCursor()
		} else {
			rl.EnableCursor()
		}
	}

	// Shooting
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) || rl.IsKeyPressed(rl.KeySpace) {
		if bullet := g.player.ShootWithAccuracy(g.aimingCircle.CurrentRadius); bullet != nil {
			g.bullets = append(g.bullets, bullet)
			g.aimingCircle.IsAiming = false // После выстрела точность сбрасывается
		}
	}

	// Right mouse button for aiming
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		g.aimingCircle.IsAiming = true
	} else if rl.IsMouseButtonReleased(rl.MouseRightButton) {
		g.aimingCircle.IsAiming = false
	}
}

func (g *Game) handleMouseAiming() {
	// Получаем позицию мыши
	mousePos := rl.GetMousePosition()
	screenCenter := rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2)
	
	// Вычисляем смещение от центра экрана
	deltaX := mousePos.X - screenCenter.X
	deltaY := mousePos.Y - screenCenter.Y
	
	// Вычисляем угол поворота башни относительно корпуса танка
	mouseAngle := math.Atan2(float64(deltaX), float64(-deltaY)) // -deltaY потому что Y инвертирован
	
	// Устанавливаем поворот башни
	g.player.SetTurretRotation(float32(mouseAngle))
	
	// Возвращаем курсор в центр экрана для непрерывного управления
	rl.SetMousePosition(int(screenCenter.X), int(screenCenter.Y))
}

func (g *Game) updateAiming() {
	if g.aimingCircle.IsAiming {
		// Сведение - уменьшаем круг точности
		g.aimingCircle.CurrentRadius -= g.aimingCircle.ShrinkSpeed
		if g.aimingCircle.CurrentRadius < g.aimingCircle.MinRadius {
			g.aimingCircle.CurrentRadius = g.aimingCircle.MinRadius
		}
	} else {
		// Разведение - увеличиваем круг точности
		g.aimingCircle.CurrentRadius += g.aimingCircle.ExpandSpeed
		if g.aimingCircle.CurrentRadius > g.aimingCircle.MaxRadius {
			g.aimingCircle.CurrentRadius = g.aimingCircle.MaxRadius
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

	// Aim turret at player
	turretAngle := math.Atan2(float64(dx), float64(dz)) - float64(enemy.Rotation)
	enemy.SetTurretRotation(float32(turretAngle))

	// Shoot occasionally
	if g.gameTime%180 == 0 && distance < 30 {
		if bullet := enemy.ShootWithAccuracy(30.0); bullet != nil {
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

	// Aiming circle (crosshair)
	if g.mouseAiming {
		g.drawAimingCircle()
	}

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
	enemyText := fmt.Sprintf("Enemies: %d", aliveEnemies)
	rl.DrawText(enemyText, 10, 60, 20, rl.Black)

	// Controls
	controlsText := "WASD - Move, Mouse - Aim, LMB/Space - Shoot, RMB - Precise Aim, Tab - Toggle Mouse"
	rl.DrawText(controlsText, 10, 720, 16, rl.DarkGray)
	
	// Aiming mode indicator
	if g.mouseAiming {
		rl.DrawText("Mouse Aiming: ON", 10, 85, 20, rl.Green)
	} else {
		rl.DrawText("Mouse Aiming: OFF", 10, 85, 20, rl.Red)
	}
}

func (g *Game) drawAimingCircle() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	// Цвет круга зависит от точности
	var circleColor rl.Color
	accuracy := (g.aimingCircle.MaxRadius - g.aimingCircle.CurrentRadius) / (g.aimingCircle.MaxRadius - g.aimingCircle.MinRadius)
	
	if accuracy > 0.8 {
		circleColor = rl.Green
	} else if accuracy > 0.5 {
		circleColor = rl.Yellow
	} else {
		circleColor = rl.Red
	}

	// Рисуем круг точности
	rl.DrawCircleLines(int32(centerX), int32(centerY), g.aimingCircle.CurrentRadius, circleColor)
	
	// Рисуем крестик в центре
	crossSize := float32(10)
	rl.DrawLine(int32(centerX-crossSize), int32(centerY), int32(centerX+crossSize), int32(centerY), rl.White)
	rl.DrawLine(int32(centerX), int32(centerY-crossSize), int32(centerX), int32(centerY+crossSize), rl.White)
	
	// Показываем статус сведения
	if g.aimingCircle.IsAiming {
		rl.DrawText("AIMING...", int32(centerX-40), int32(centerY+g.aimingCircle.CurrentRadius+20), 20, circleColor)
	}
	
	// Показываем процент точности
	accuracyPercent := int32(accuracy * 100)
	accuracyText := fmt.Sprintf("Accuracy: %d%%", accuracyPercent)
	rl.DrawText(accuracyText, int32(centerX-60), int32(centerY-g.aimingCircle.CurrentRadius-30), 20, circleColor)
}