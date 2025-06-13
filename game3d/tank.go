package game3d

import (
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tank struct {
	Position       rl.Vector3
	Rotation       float32 // Body rotation
	TurretRotation float32 // Turret rotation relative to body
	Speed          float32
	TurnSpeed      float32
	Health         int
	MaxHealth      int
	IsPlayer       bool
	LastShot       time.Time
	ShotCooldown   time.Duration
}

func NewTank(position rl.Vector3, isPlayer bool) *Tank {
	return &Tank{
		Position:       position,
		Rotation:       0,
		TurretRotation: 0,
		Speed:          0.2,
		TurnSpeed:      0.03,
		Health:         100,
		MaxHealth:      100,
		IsPlayer:       isPlayer,
		ShotCooldown:   time.Millisecond * 800,
	}
}

func (t *Tank) Update() {
	// Keep tank within map bounds
	if t.Position.X < -MapSize {
		t.Position.X = -MapSize
	}
	if t.Position.X > MapSize {
		t.Position.X = MapSize
	}
	if t.Position.Z < -MapSize {
		t.Position.Z = -MapSize
	}
	if t.Position.Z > MapSize {
		t.Position.Z = MapSize
	}
}

func (t *Tank) MoveForward() {
	t.Position.X += float32(math.Sin(float64(t.Rotation))) * t.Speed
	t.Position.Z += float32(math.Cos(float64(t.Rotation))) * t.Speed
}

func (t *Tank) MoveBackward() {
	t.Position.X -= float32(math.Sin(float64(t.Rotation))) * t.Speed * 0.5
	t.Position.Z -= float32(math.Cos(float64(t.Rotation))) * t.Speed * 0.5
}

func (t *Tank) TurnLeft() {
	t.Rotation -= t.TurnSpeed
}

func (t *Tank) TurnRight() {
	t.Rotation += t.TurnSpeed
}

func (t *Tank) TurretLeft() {
	t.TurretRotation -= t.TurnSpeed * 0.8
}

func (t *Tank) TurretRight() {
	t.TurretRotation += t.TurnSpeed * 0.8
}

// Новый метод для установки поворота башни напрямую (для мыши)
func (t *Tank) SetTurretRotation(angle float32) {
	t.TurretRotation = angle
	
	// Ограничиваем поворот башни (например, ±180 градусов)
	maxTurretAngle := float32(math.Pi)
	if t.TurretRotation > maxTurretAngle {
		t.TurretRotation = maxTurretAngle
	}
	if t.TurretRotation < -maxTurretAngle {
		t.TurretRotation = -maxTurretAngle
	}
}

func (t *Tank) Shoot() *Bullet {
	now := time.Now()
	if now.Sub(t.LastShot) < t.ShotCooldown {
		return nil
	}

	t.LastShot = now

	// Calculate bullet spawn position (at the end of the cannon)
	totalRotation := t.Rotation + t.TurretRotation
	cannonLength := float32(3.0)
	
	bulletX := t.Position.X + float32(math.Sin(float64(totalRotation)))*cannonLength
	bulletZ := t.Position.Z + float32(math.Cos(float64(totalRotation)))*cannonLength
	bulletY := t.Position.Y + 1.0

	bulletPos := rl.NewVector3(bulletX, bulletY, bulletZ)

	return NewBullet(bulletPos, totalRotation, t.IsPlayer)
}

// Новый метод стрельбы с учетом точности
func (t *Tank) ShootWithAccuracy(accuracyRadius float32) *Bullet {
	now := time.Now()
	if now.Sub(t.LastShot) < t.ShotCooldown {
		return nil
	}

	t.LastShot = now

	// Calculate bullet spawn position (at the end of the cannon)
	totalRotation := t.Rotation + t.TurretRotation
	cannonLength := float32(3.0)
	
	bulletX := t.Position.X + float32(math.Sin(float64(totalRotation)))*cannonLength
	bulletZ := t.Position.Z + float32(math.Cos(float64(totalRotation)))*cannonLength
	bulletY := t.Position.Y + 1.0

	bulletPos := rl.NewVector3(bulletX, bulletY, bulletZ)

	// Добавляем разброс в зависимости от точности
	// Чем больше accuracyRadius, тем больше разброс
	spreadFactor := accuracyRadius / 100.0 // Нормализуем разброс
	angleSpread := (rand.Float32() - 0.5) * spreadFactor * 0.2 // ±10% от разброса
	
	finalAngle := totalRotation + angleSpread

	return NewBullet(bulletPos, finalAngle, t.IsPlayer)
}

func (t *Tank) TakeDamage(damage int) {
	t.Health -= damage
	if t.Health < 0 {
		t.Health = 0
	}
}

func (t *Tank) Draw() {
	if t.Health <= 0 {
		return
	}

	// Tank colors
	bodyColor := rl.Gray
	turretColor := rl.DarkGray
	if t.IsPlayer {
		bodyColor = rl.Blue
		turretColor = rl.DarkBlue
	} else {
		bodyColor = rl.Red
		turretColor = rl.Maroon // Заменил DarkRed на Maroon
	}

	// Draw tank body with rotation
	rl.PushMatrix()
	rl.Translatef(t.Position.X, t.Position.Y, t.Position.Z)
	rl.Rotatef(t.Rotation*rl.Rad2deg, 0, 1, 0)
	rl.DrawCube(rl.NewVector3(0, 0, 0), 3, 1, 4, bodyColor)
	rl.PopMatrix()

	// Draw turret
	turretY := t.Position.Y + 0.7
	
	rl.PushMatrix()
	rl.Translatef(t.Position.X, turretY, t.Position.Z)
	rl.Rotatef((t.Rotation+t.TurretRotation)*rl.Rad2deg, 0, 1, 0)
	rl.DrawCube(rl.NewVector3(0, 0, 0), 2, 0.8, 2.5, turretColor)
	
	// Draw cannon
	rl.DrawCube(rl.NewVector3(0, 0, 2), 0.3, 0.3, 2, rl.Black)
	rl.PopMatrix()

	// Draw health bar above tank (for enemies)
	if !t.IsPlayer {
		healthBarWidth := float32(3)
		healthBarHeight := float32(0.2)
		healthPercentage := float32(t.Health) / float32(t.MaxHealth)

		barY := t.Position.Y + 3
		barPos := rl.NewVector3(t.Position.X-healthBarWidth/2, barY, t.Position.Z)

		// Background
		rl.DrawCubeV(barPos, rl.NewVector3(healthBarWidth, healthBarHeight, 0.1), rl.Gray)

		// Health
		healthColor := rl.Red
		if healthPercentage > 0.5 {
			healthColor = rl.Green
		} else if healthPercentage > 0.25 {
			healthColor = rl.Yellow
		}

		healthPos := rl.NewVector3(t.Position.X-healthBarWidth/2, barY, t.Position.Z)
		rl.DrawCubeV(healthPos, rl.NewVector3(healthBarWidth*healthPercentage, healthBarHeight, 0.1), healthColor)
	}
}