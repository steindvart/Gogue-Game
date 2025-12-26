package world

import "gogue/internal/model/primitives"

type LevelConfig struct {
	MapSize primitives.Size2D[uint]

	// Параметры комнат
	RoomMinWidth   int
	RoomMinHeight  int
	MapMaxWidth    uint
	MapMaxHeight   uint
	MinRoomPadding int

	// Параметры генерации мира
	MaxExtraPassageCount int
	NumberXYSections     int
	RoomsCount           int
	ItemsSpawnConfig     ItemSpawnConfig
	EnemiesSpawnConfig   EnemySpawnConfig
}

type ItemSpawnConfig struct {
	FoodsQuantity   uint
	ElixirsQuantity uint
	ScrollsQuantity uint
	WeaponsQuantity uint
}

type EnemySpawnConfig struct {
	Quantity uint
}

func DefaultLevelConfig(mapSize primitives.Size2D[uint]) LevelConfig {
	return LevelConfig{
		MapSize:              mapSize,
		RoomMinWidth:         3,
		RoomMinHeight:        3,
		MapMaxWidth:          200,
		MapMaxHeight:         150,
		MinRoomPadding:       1,
		MaxExtraPassageCount: 2,
		NumberXYSections:     3,
		RoomsCount:           9,
		ItemsSpawnConfig: ItemSpawnConfig{
			FoodsQuantity:   7,
			ElixirsQuantity: 5,
			ScrollsQuantity: 3,
			WeaponsQuantity: 1,
		},
		EnemiesSpawnConfig: EnemySpawnConfig{
			Quantity: 10,
		},
	}
}
