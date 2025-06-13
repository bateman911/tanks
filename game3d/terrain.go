package game3d

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Obstacle struct {
	Position rl.Vector3
	Size     rl.Vector3
	Color    rl.Color
	Type     string
}

type Terrain struct {
	Obstacles []Obstacle
}

func NewTerrain() *Terrain {
	obstacles := make([]Obstacle, 0)

	// Generate random obstacles (buildings, rocks, trees)
	for i := 0; i < 30; i++ {
		obstacle := Obstacle{
			Position: rl.NewVector3(
				(rand.Float32()-0.5)*MapSize*1.5,
				0,
				(rand.Float32()-0.5)*MapSize*1.5,
			),
			Size: rl.NewVector3(
				2+rand.Float32()*4,
				1+rand.Float32()*3,
				2+rand.Float32()*4,
			),
			Color: rl.Brown,
			Type:  "building",
		}
		obstacles = append(obstacles, obstacle)
	}

	// Add some trees
	for i := 0; i < 50; i++ {
		obstacle := Obstacle{
			Position: rl.NewVector3(
				(rand.Float32()-0.5)*MapSize*1.8,
				0,
				(rand.Float32()-0.5)*MapSize*1.8,
			),
			Size: rl.NewVector3(
				0.5+rand.Float32(),
				3+rand.Float32()*2,
				0.5+rand.Float32(),
			),
			Color: rl.DarkGreen,
			Type:  "tree",
		}
		obstacles = append(obstacles, obstacle)
	}

	return &Terrain{
		Obstacles: obstacles,
	}
}

func (t *Terrain) Draw() {
	// Draw ground plane
	rl.DrawPlane(rl.NewVector3(0, -0.5, 0), rl.NewVector2(MapSize*2, MapSize*2), rl.Green)

	// Draw obstacles
	for _, obstacle := range t.Obstacles {
		switch obstacle.Type {
		case "building":
			rl.DrawCubeV(
				rl.NewVector3(obstacle.Position.X, obstacle.Position.Y+obstacle.Size.Y/2, obstacle.Position.Z),
				obstacle.Size,
				obstacle.Color,
			)
		case "tree":
			// Tree trunk
			rl.DrawCubeV(
				rl.NewVector3(obstacle.Position.X, obstacle.Position.Y+obstacle.Size.Y/2, obstacle.Position.Z),
				rl.NewVector3(obstacle.Size.X, obstacle.Size.Y, obstacle.Size.Z),
				rl.Brown,
			)
			// Tree crown
			rl.DrawSphere(
				rl.NewVector3(obstacle.Position.X, obstacle.Position.Y+obstacle.Size.Y+1, obstacle.Position.Z),
				1.5,
				rl.DarkGreen,
			)
		}
	}
}