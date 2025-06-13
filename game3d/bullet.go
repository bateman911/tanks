package game3d

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Position   rl.Vector3
	Velocity   rl.Vector3
	Speed      float32
	LifeTime   int
	FromPlayer bool
}

func NewBullet(position rl.Vector3, angle float32, fromPlayer bool) *Bullet {
	speed := float32(0.8)
	
	velocity := rl.NewVector3(
		float32(math.Sin(float64(angle)))*speed,
		0,
		float32(math.Cos(float64(angle)))*speed,
	)

	return &Bullet{
		Position:   position,
		Velocity:   velocity,
		Speed:      speed,
		LifeTime:   300, // 5 seconds at 60 FPS
		FromPlayer: fromPlayer,
	}
}

func (b *Bullet) Update() {
	b.Position.X += b.Velocity.X
	b.Position.Y += b.Velocity.Y
	b.Position.Z += b.Velocity.Z
	b.LifeTime--
}

func (b *Bullet) Draw() {
	bulletColor := rl.Yellow
	if !b.FromPlayer {
		bulletColor = rl.Orange
	}

	rl.DrawSphere(b.Position, 0.2, bulletColor)
	
	// Draw bullet trail
	trailPos := rl.NewVector3(
		b.Position.X - b.Velocity.X*2,
		b.Position.Y - b.Velocity.Y*2,
		b.Position.Z - b.Velocity.Z*2,
	)
	
	rl.DrawLine3D(b.Position, trailPos, bulletColor)
}