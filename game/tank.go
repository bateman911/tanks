package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tank struct {
	X, Y          float64
	Angle         float64
	Speed         float64
	TurnSpeed     float64
	Health        int
	MaxHealth     int
	IsPlayer      bool
	LastShot      time.Time
	ShotCooldown  time.Duration
}

func NewTank(x, y float64, isPlayer bool) *Tank {
	return &Tank{
		X:            x,
		Y:            y,
		Angle:        0,
		Speed:        2.0,
		TurnSpeed:    0.05,
		Health:       100,
		MaxHealth:    100,
		IsPlayer:     isPlayer,
		ShotCooldown: time.Millisecond * 500,
	}
}

func (t *Tank) Update() {
	// Tank updates are handled in the main game loop
}

func (t *Tank) MoveForward() {
	t.X += math.Cos(t.Angle) * t.Speed
	t.Y += math.Sin(t.Angle) * t.Speed
	
	// Keep tank within map bounds
	if t.X < 0 {
		t.X = 0
	}
	if t.X > mapWidth {
		t.X = mapWidth
	}
	if t.Y < 0 {
		t.Y = 0
	}
	if t.Y > mapHeight {
		t.Y = mapHeight
	}
}

func (t *Tank) MoveBackward() {
	t.X -= math.Cos(t.Angle) * t.Speed * 0.5
	t.Y -= math.Sin(t.Angle) * t.Speed * 0.5
	
	// Keep tank within map bounds
	if t.X < 0 {
		t.X = 0
	}
	if t.X > mapWidth {
		t.X = mapWidth
	}
	if t.Y < 0 {
		t.Y = 0
	}
	if t.Y > mapHeight {
		t.Y = mapHeight
	}
}

func (t *Tank) TurnLeft() {
	t.Angle -= t.TurnSpeed
}

func (t *Tank) TurnRight() {
	t.Angle += t.TurnSpeed
}

func (t *Tank) Shoot() *Bullet {
	now := time.Now()
	if now.Sub(t.LastShot) < t.ShotCooldown {
		return nil
	}
	
	t.LastShot = now
	
	// Calculate bullet spawn position (at the front of the tank)
	bulletX := t.X + math.Cos(t.Angle)*25
	bulletY := t.Y + math.Sin(t.Angle)*25
	
	return NewBullet(bulletX, bulletY, t.Angle, t.IsPlayer)
}

func (t *Tank) TakeDamage(damage int) {
	t.Health -= damage
	if t.Health < 0 {
		t.Health = 0
	}
}

func (t *Tank) Draw(screen *ebiten.Image, camera *Camera) {
	if t.Health <= 0 {
		return
	}

	// Convert world coordinates to screen coordinates
	screenX := float32(t.X - camera.X)
	screenY := float32(t.Y - camera.Y)
	
	// Tank body (rectangle)
	tankWidth := float32(30)
	tankHeight := float32(20)
	
	// Tank color
	tankColor := color.RGBA{100, 100, 100, 255}
	if t.IsPlayer {
		tankColor = color.RGBA{0, 100, 200, 255}
	} else {
		tankColor = color.RGBA{200, 100, 0, 255}
	}
	
	// Draw tank body (simplified as rectangle for now)
	// We'll rotate around the center
	centerX := screenX
	centerY := screenY
	
	// Calculate rotated corners
	cos := float32(math.Cos(t.Angle))
	sin := float32(math.Sin(t.Angle))
	
	// Tank body corners relative to center
	corners := []struct{ x, y float32 }{
		{-tankWidth / 2, -tankHeight / 2},
		{tankWidth / 2, -tankHeight / 2},
		{tankWidth / 2, tankHeight / 2},
		{-tankWidth / 2, tankHeight / 2},
	}
	
	// Rotate and translate corners
	var rotatedCorners []struct{ x, y float32 }
	for _, corner := range corners {
		rotX := corner.x*cos - corner.y*sin + centerX
		rotY := corner.x*sin + corner.y*cos + centerY
		rotatedCorners = append(rotatedCorners, struct{ x, y float32 }{rotX, rotY})
	}
	
	// Draw tank body as filled polygon
	var path vector.Path
	path.MoveTo(rotatedCorners[0].x, rotatedCorners[0].y)
	for i := 1; i < len(rotatedCorners); i++ {
		path.LineTo(rotatedCorners[i].x, rotatedCorners[i].y)
	}
	path.Close()
	
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(tankColor.R) / 255
		vs[i].ColorG = float32(tankColor.G) / 255
		vs[i].ColorB = float32(tankColor.B) / 255
		vs[i].ColorA = float32(tankColor.A) / 255
	}
	screen.DrawTriangles(vs, is, nil, nil)
	
	// Draw tank cannon
	cannonLength := float32(35)
	cannonEndX := centerX + cannonLength*cos
	cannonEndY := centerY + cannonLength*sin
	
	vector.StrokeLine(screen, centerX, centerY, cannonEndX, cannonEndY, 3, color.RGBA{50, 50, 50, 255}, false)
	
	// Draw health bar above tank
	if !t.IsPlayer {
		healthBarWidth := float32(30)
		healthBarHeight := float32(4)
		healthPercentage := float32(t.Health) / float32(t.MaxHealth)
		
		barX := centerX - healthBarWidth/2
		barY := centerY - 25
		
		// Background
		vector.DrawFilledRect(screen, barX, barY, healthBarWidth, healthBarHeight, color.RGBA{100, 100, 100, 255}, false)
		
		// Health
		healthColor := color.RGBA{255, 0, 0, 255}
		if healthPercentage > 0.5 {
			healthColor = color.RGBA{0, 255, 0, 255}
		} else if healthPercentage > 0.25 {
			healthColor = color.RGBA{255, 255, 0, 255}
		}
		
		vector.DrawFilledRect(screen, barX, barY, healthBarWidth*healthPercentage, healthBarHeight, healthColor, false)
	}
}