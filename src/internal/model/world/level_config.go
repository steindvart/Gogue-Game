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
	ItemCounts           ItemSpawnConfig
}

type ItemSpawnConfig struct {
	FoodsQuntity   uint
	ElixirsQuntity uint
	ScrollsQuntity uint
	WeaponsQuntity uint
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
		ItemCounts: ItemSpawnConfig{
			FoodsQuntity:   7,
			ElixirsQuntity: 5,
			ScrollsQuntity: 3,
			WeaponsQuntity: 1,
		},
	}
}
