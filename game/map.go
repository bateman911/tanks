package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Obstacle struct {
	X, Y   float64
	Width  float64
	Height float64
	Type   string
}

type GameMap struct {
	Obstacles []Obstacle
}

func NewGameMap() *GameMap {
	obstacles := make([]Obstacle, 0)
	
	// Generate random obstacles (trees, rocks, buildings)
	for i := 0; i < 50; i++ {
		obstacle := Obstacle{
			X:      rand.Float64() * mapWidth,
			Y:      rand.Float64() * mapHeight,
			Width:  20 + rand.Float64()*30,
			Height: 20 + rand.Float64()*30,
			Type:   "tree",
		}
		obstacles = append(obstacles, obstacle)
	}
	
	// Add some larger buildings
	for i := 0; i < 10; i++ {
		obstacle := Obstacle{
			X:      rand.Float64() * mapWidth,
			Y:      rand.Float64() * mapHeight,
			Width:  50 + rand.Float64()*50,
			Height: 50 + rand.Float64()*50,
			Type:   "building",
		}
		obstacles = append(obstacles, obstacle)
	}
	
	return &GameMap{
		Obstacles: obstacles,
	}
}

func (m *GameMap) Draw(screen *ebiten.Image, camera *Camera) {
	for _, obstacle := range m.Obstacles {
		screenX := float32(obstacle.X - camera.X)
		screenY := float32(obstacle.Y - camera.Y)
		
		// Only draw if visible on screen
		if screenX > -obstacle.Width && screenX < screenWidth+obstacle.Width &&
			screenY > -obstacle.Height && screenY < screenHeight+obstacle.Height {
			
			var obstacleColor color.RGBA
			switch obstacle.Type {
			case "tree":
				obstacleColor = color.RGBA{34, 100, 34, 255}
			case "building":
				obstacleColor = color.RGBA{100, 100, 100, 255}
			default:
				obstacleColor = color.RGBA{139, 69, 19, 255}
			}
			
			vector.DrawFilledRect(screen, screenX, screenY, float32(obstacle.Width), float32(obstacle.Height), obstacleColor, false)
		}
	}
}