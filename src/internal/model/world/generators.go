package world

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type RoomGenerator interface {
	GenerateRooms(config LevelConfig, random utils.Randomizer) ([]Room, primitives.Box, error)
}

type PassageGenerator interface {
	GeneratePassages(rooms []Room, config LevelConfig, random utils.Randomizer) ([]Passage, error)
}

type EntitySpawner interface {
	SpawnEntities(
		rooms []Room,
		itemsConfig ItemSpawnConfig,
		enemiesConfig EnemySpawnConfig,
		random utils.Randomizer,
	) (*SpawnedEntities, error)
}

type SpawnedEntities struct {
	Items   []primitives.Positional2D[int]
	Enemies []primitives.Positional2D[int]
}

type PlayerSpawner interface {
	SpawnPlayer(rooms []Room, random utils.Randomizer) (*entities.Player, error)
	GetStartPosition(rooms []Room, random utils.Randomizer) (*primitives.Point2D[int], error)
}
