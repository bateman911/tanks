package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	mapWidth     = 2048
	mapHeight    = 1536
)

type Game struct {
	player   *Tank
	enemies  []*Tank
	camera   *Camera
	gameMap  *GameMap
	bullets  []*Bullet
	gameTime int
}

func NewGame() *Game {
	player := NewTank(100, 100, true)
	
	enemies := []*Tank{
		NewTank(500, 300, false),
		NewTank(800, 200, false),
		NewTank(300, 600, false),
	}

	camera := NewCamera(0, 0)
	gameMap := NewGameMap()

	return &Game{
		player:  player,
		enemies: enemies,
		camera:  camera,
		gameMap: gameMap,
		bullets: make([]*Bullet, 0),
	}
}

func (g *Game) Update() error {
	g.gameTime++
	
	// Update player
	g.player.Update()
	
	// Handle player input
	g.handleInput()
	
	// Update enemies with simple AI
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
		if bullet.X < 0 || bullet.X > mapWidth || bullet.Y < 0 || bullet.Y > mapHeight || bullet.LifeTime <= 0 {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			continue
		}
		
		// Check bullet collisions
		g.checkBulletCollisions(bullet, i)
	}
	
	// Update camera to follow player
	g.camera.Follow(g.player.X, g.player.Y, screenWidth, screenHeight, mapWidth, mapHeight)
	
	return nil
}

func (g *Game) handleInput() {
	if g.player.Health <= 0 {
		return
	}

	// Movement
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.MoveForward()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.MoveBackward()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.TurnLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.TurnRight()
	}

	// Shooting
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if bullet := g.player.Shoot(); bullet != nil {
			g.bullets = append(g.bullets, bullet)
		}
	}
}

func (g *Game) updateEnemyAI(enemy *Tank) {
	// Simple AI: move towards player and shoot occasionally
	dx := g.player.X - enemy.X
	dy := g.player.Y - enemy.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	if distance > 50 {
		// Move towards player
		targetAngle := math.Atan2(dy, dx)
		angleDiff := targetAngle - enemy.Angle
		
		// Normalize angle difference
		for angleDiff > math.Pi {
			angleDiff -= 2 * math.Pi
		}
		for angleDiff < -math.Pi {
			angleDiff += 2 * math.Pi
		}
		
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
	if g.gameTime%120 == 0 && distance < 300 {
		if bullet := enemy.Shoot(); bullet != nil {
			g.bullets = append(g.bullets, bullet)
		}
	}
}

func (g *Game) checkBulletCollisions(bullet *Bullet, bulletIndex int) {
	// Check collision with player
	if !bullet.FromPlayer && g.player.Health > 0 {
		if g.checkCollision(bullet.X, bullet.Y, g.player.X, g.player.Y, 20) {
			g.player.TakeDamage(25)
			g.bullets = append(g.bullets[:bulletIndex], g.bullets[bulletIndex+1:]...)
			return
		}
	}
	
	// Check collision with enemies
	if bullet.FromPlayer {
		for _, enemy := range g.enemies {
			if enemy.Health > 0 && g.checkCollision(bullet.X, bullet.Y, enemy.X, enemy.Y, 20) {
				enemy.TakeDamage(25)
				g.bullets = append(g.bullets[:bulletIndex], g.bullets[bulletIndex+1:]...)
				return
			}
		}
	}
}

func (g *Game) checkCollision(x1, y1, x2, y2, radius float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < radius
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{34, 139, 34, 255}) // Forest green background
	
	// Draw map elements (trees, obstacles)
	g.gameMap.Draw(screen, g.camera)
	
	// Draw tanks
	g.player.Draw(screen, g.camera)
	for _, enemy := range g.enemies {
		if enemy.Health > 0 {
			enemy.Draw(screen, g.camera)
		}
	}
	
	// Draw bullets
	for _, bullet := range g.bullets {
		bullet.Draw(screen, g.camera)
	}
	
	// Draw UI
	g.drawUI(screen)
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Health bar
	healthBarWidth := 200.0
	healthBarHeight := 20.0
	healthPercentage := float64(g.player.Health) / 100.0
	
	// Background
	vector.DrawFilledRect(screen, 10, 10, float32(healthBarWidth), float32(healthBarHeight), color.RGBA{100, 100, 100, 255}, false)
	
	// Health
	healthColor := color.RGBA{255, 0, 0, 255}
	if healthPercentage > 0.5 {
		healthColor = color.RGBA{0, 255, 0, 255}
	} else if healthPercentage > 0.25 {
		healthColor = color.RGBA{255, 255, 0, 255}
	}
	
	vector.DrawFilledRect(screen, 10, 10, float32(healthBarWidth*healthPercentage), float32(healthBarHeight), healthColor, false)
	
	// Health text
	ebitenutil.DebugPrintAt(screen, "Health", 10, 35)
	
	// Game status
	if g.player.Health <= 0 {
		ebitenutil.DebugPrintAt(screen, "GAME OVER - Press R to restart", screenWidth/2-100, screenHeight/2)
	}
	
	// Enemy count
	aliveEnemies := 0
	for _, enemy := range g.enemies {
		if enemy.Health > 0 {
			aliveEnemies++
		}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Enemies: %d", aliveEnemies), 10, 55)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}