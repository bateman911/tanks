package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Bullet struct {
	X, Y       float64
	VelX, VelY float64
	Speed      float64
	LifeTime   int
	FromPlayer bool
}

func NewBullet(x, y, angle float64, fromPlayer bool) *Bullet {
	speed := 8.0
	return &Bullet{
		X:          x,
		Y:          y,
		VelX:       math.Cos(angle) * speed,
		VelY:       math.Sin(angle) * speed,
		Speed:      speed,
		LifeTime:   180, // 3 seconds at 60 FPS
		FromPlayer: fromPlayer,
	}
}

func (b *Bullet) Update() {
	b.X += b.VelX
	b.Y += b.VelY
	b.LifeTime--
}

func (b *Bullet) Draw(screen *ebiten.Image, camera *Camera) {
	screenX := float32(b.X - camera.X)
	screenY := float32(b.Y - camera.Y)
	
	bulletColor := color.RGBA{255, 255, 0, 255}
	if !b.FromPlayer {
		bulletColor = color.RGBA{255, 100, 100, 255}
	}
	
	vector.DrawFilledCircle(screen, screenX, screenY, 3, bulletColor, false)
}