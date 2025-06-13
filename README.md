# Tanks - Cross-Platform Tank Game

A cross-platform tank battle game inspired by World of Tanks Blitz, built with Go and Ebitengine.

## Features

- **Cross-platform**: Runs on Windows, macOS, and Linux
- **Real-time tank battles**: Control your tank with WASD or arrow keys
- **AI enemies**: Fight against computer-controlled tanks
- **Physics-based movement**: Realistic tank movement and rotation
- **Health system**: Damage and health mechanics
- **Dynamic camera**: Camera follows the player tank
- **Procedural map**: Randomly generated obstacles and terrain

## Controls

- **W/↑**: Move forward
- **S/↓**: Move backward  
- **A/←**: Turn left
- **D/→**: Turn right
- **Space**: Shoot

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd tanks
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
GOOS=windows GOARCH=amd64 go build -o tanks.exe main.go
```

#### macOS
```bash
GOOS=darwin GOARCH=amd64 go build -o tanks-mac main.go
```

#### Linux
```bash
GOOS=linux GOARCH=amd64 go build -o tanks-linux main.go
```

## Game Architecture

The game is structured with the following components:

- **Game**: Main game loop and state management
- **Tank**: Player and enemy tank entities with movement and combat
- **Bullet**: Projectile system with collision detection
- **Camera**: Viewport management that follows the player
- **GameMap**: Procedural obstacle generation and rendering

## Future Enhancements

- **Mobile Support**: Android and iOS versions using Ebitengine's mobile support
- **Multiplayer**: Network-based multiplayer battles
- **More Tank Types**: Different tank classes with unique abilities
- **Power-ups**: Collectible items for temporary advantages
- **Sound Effects**: Audio feedback for actions and events
- **Better Graphics**: Sprite-based rendering instead of geometric shapes
- **Collision Detection**: Tank-obstacle collision system
- **Game Modes**: Different battle scenarios and objectives

## Technical Details

Built with:
- **Go**: Primary programming language
- **Ebitengine**: 2D game engine for cross-platform development
- **Vector Graphics**: Simple geometric rendering for prototyping

The game uses a component-based architecture where each game entity (tanks, bullets, obstacles) is managed independently with its own update and draw methods.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test on multiple platforms
5. Submit a pull request

## License

This project is open source and available under the MIT License.