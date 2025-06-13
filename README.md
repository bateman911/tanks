# 3D Tanks - Cross-Platform Tank Battle Game

A 3D tank battle game inspired by World of Tanks, built with Go and Raylib-go for cross-platform compatibility.

## Features

- **Full 3D Graphics**: Real 3D tanks, terrain, and environment
- **Cross-platform**: Runs on Windows, macOS, and Linux
- **Third-person Camera**: Dynamic camera that follows the player tank
- **Separate Turret Control**: Independent turret rotation from tank body
- **AI Enemies**: Computer-controlled tanks with basic AI
- **3D Physics**: Realistic 3D movement and bullet trajectories
- **Health System**: Damage mechanics with visual health bars
- **Procedural Terrain**: Randomly generated obstacles, buildings, and trees
- **3D Audio Ready**: Structure prepared for 3D positional audio

## Controls

- **W**: Move forward
- **S**: Move backward  
- **A**: Turn tank left
- **D**: Turn tank right
- **←**: Rotate turret left
- **→**: Rotate turret right
- **Space**: Shoot
- **ESC**: Exit game

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- OpenGL support (usually built into modern systems)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd tanks3d
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the game:
```bash
go run main.go
```

### Building for Different Platforms

#### Windows
```bash
GOOS=windows GOARCH=amd64 go build -o tanks3d.exe main.go
```

#### macOS
```bash
GOOS=darwin GOARCH=amd64 go build -o tanks3d-mac main.go
```

#### Linux
```bash
GOOS=linux GOARCH=amd64 go build -o tanks3d-linux main.go
```

## 3D Game Architecture

The game uses a modern 3D architecture:

- **Game3D**: Main game loop with 3D scene management
- **Tank**: 3D tank entities with separate body and turret rotation
- **Bullet**: 3D projectiles with realistic trajectories
- **Camera3D**: Third-person 3D camera system
- **Terrain**: 3D procedural world generation with obstacles
- **3D Rendering**: OpenGL-based rendering through Raylib

## Technical Features

### 3D Graphics
- Real-time 3D rendering with OpenGL
- Perspective projection and 3D transformations
- 3D models represented as geometric primitives
- Dynamic lighting and shadows (ready for implementation)

### 3D Physics
- 3D collision detection
- Realistic bullet trajectories
- Tank movement in 3D space
- Boundary checking for 3D world

### Camera System
- Third-person camera following the player
- Smooth camera transitions
- Configurable camera distance and height
- Target tracking system

## Future Enhancements

### Graphics
- **3D Models**: Load actual tank models (.obj, .fbx)
- **Textures**: PBR materials and texture mapping
- **Lighting**: Dynamic lighting system with shadows
- **Particle Effects**: Explosions, smoke, muzzle flashes
- **Skybox**: 3D environment backgrounds

### Gameplay
- **Terrain Collision**: Tank-obstacle collision detection
- **Destructible Environment**: Breakable obstacles
- **Multiple Tank Types**: Different tank classes
- **Power-ups**: Collectible items in 3D space
- **Multiplayer**: Network-based 3D battles

### Audio
- **3D Positional Audio**: Spatial sound effects
- **Engine Sounds**: Tank movement audio
- **Weapon Audio**: Shooting and explosion sounds
- **Ambient Audio**: Environmental sounds

### Mobile Support
- **Touch Controls**: Mobile-friendly input system
- **Performance Optimization**: Mobile GPU optimization
- **UI Scaling**: Responsive UI for different screen sizes

## Performance Considerations

- Efficient 3D rendering with frustum culling
- Level-of-detail (LOD) system for distant objects
- Optimized collision detection
- Memory management for 3D assets

## Dependencies

- **Raylib-go**: 3D graphics and game engine
- **OpenGL**: Hardware-accelerated 3D rendering
- **Go**: Cross-platform runtime

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement 3D features
4. Test on multiple platforms
5. Submit a pull request

## License

This project is open source and available under the MIT License.