package game

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{X: x, Y: y}
}

func (c *Camera) Follow(targetX, targetY float64, screenWidth, screenHeight, mapWidth, mapHeight int) {
	// Center camera on target
	c.X = targetX - float64(screenWidth)/2
	c.Y = targetY - float64(screenHeight)/2
	
	// Keep camera within map bounds
	if c.X < 0 {
		c.X = 0
	}
	if c.X > float64(mapWidth-screenWidth) {
		c.X = float64(mapWidth - screenWidth)
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if c.Y > float64(mapHeight-screenHeight) {
		c.Y = float64(mapHeight - screenHeight)
	}
}