package save

import (
	"container/list"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
	"gogue/internal/presentation/dto"
	"math/rand"
	"time"
)

func ConvertLevelToDto(level *world.Level, size primitives.Size2D[uint]) dto.GameSaveDto {
	roomsDto := make([]dto.RoomSaveDto, len(level.Rooms))
	for i, room := range level.Rooms {
		roomsDto[i] = dto.RoomSaveDto{
			Box:   *room.Box,
			Doors: room.Doors,
		}
	}

	fogOfWarList := make([]primitives.Point2D[int], 0, len(level.FogOfWar.ExploredTiles))
	for point := range level.FogOfWar.ExploredTiles {
		fogOfWarList = append(fogOfWarList, point)
	}
	fogOfWarDto := dto.FogOfWarDto{
		ExploredTiles: fogOfWarList,
		Width:         level.FogOfWar.Width,
		Height:        level.FogOfWar.Height,
	}

	playerDto := dto.PlayerSaveDto{
		Position:         level.Player.GetPosition(),
		MaxHealth:        level.Player.MaxHealth,
		Health:           level.Player.Health,
		Agility:          level.Player.Agility,
		Strength:         level.Player.Strength,
		TemporaryEffects: level.Player.TemporaryEffects,
		ViewRadius:       level.Player.ViewRadius,
		Weapon:           level.Player.Weapon,
		Backpack: dto.BackpackSaveDto{
			Capacity:  level.Player.Backpack.Capacity,
			ItemsNum:  level.Player.Backpack.ItemsNum,
			Elixirs:   MapListToSlice[*items.Elixir](level.Player.Backpack.Elixirs),
			Scrolls:   MapListToSlice[*items.Scroll](level.Player.Backpack.Scrolls),
			Foods:     MapListToSlice[*items.Food](level.Player.Backpack.Foods),
			Weapons:   MapListToSlice[*items.Weapon](level.Player.Backpack.Weapons),
			Treasures: level.Player.Backpack.Treasures,
		},
		EnemiesKilled: level.Player.EnemiesKilled,
		FoodEaten:     level.Player.FoodEaten,
		ElixirsDrunk:  level.Player.ElixirsDrunk,
		ScrollsRead:   level.Player.ScrollsRead,
		HitsDealt:     level.Player.HitsDealt,
		HitsMissed:    level.Player.HitsMissed,
		CellsMoved:    level.Player.CellsMoved,
	}

	return dto.GameSaveDto{
		LevelNumber:  level.Number,
		Size:         size,
		Rooms:        roomsDto,
		Passages:     level.Passages,
		FinishPortal: level.FinishPortal,
		FogOfWar:     fogOfWarDto,
		Player:       playerDto,
		Items:        level.Items,
		Enemies:      level.Enemies,
	}
}

func RestoreLevelFromDto(dto dto.GameSaveDto) (*world.Level, error) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := world.NewLevelWithDefaults(source, primitives.Size2D[uint]{
		Width:  uint(dto.FogOfWar.Width),
		Height: uint(dto.FogOfWar.Height),
	})

	level.Number = dto.LevelNumber

	level.Rooms = make([]world.Room, len(dto.Rooms))
	for i, roomDto := range dto.Rooms {
		level.Rooms[i] = world.Room{
			Box:   &roomDto.Box,
			Doors: roomDto.Doors,
		}
	}

	level.Passages = dto.Passages

	level.FinishPortal = dto.FinishPortal

	fogMap := make(map[primitives.Point2D[int]]bool, len(dto.FogOfWar.ExploredTiles))
	for _, p := range dto.FogOfWar.ExploredTiles {
		fogMap[p] = true
	}
	level.FogOfWar = &world.FogOfWar{
		ExploredTiles: fogMap,
		Width:         dto.FogOfWar.Width,
		Height:        dto.FogOfWar.Height,
	}

	level.Player = restorePlayer(dto.Player)
	level.Items = dto.Items
	level.Enemies = dto.Enemies

	return level, nil
}

func restorePlayer(dto dto.PlayerSaveDto) *entities.Player {
	character := &entities.Character{
		Box: &primitives.Box{
			Point: dto.Position,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		Attributes: &primitives.Attributes{
			MaxHealth: dto.MaxHealth,
			Health:    dto.Health,
			Agility:   dto.Agility,
			Strength:  dto.Strength,
		},
		TemporaryEffects: dto.TemporaryEffects,
	}

	return &entities.Player{
		Character: character,
		Backpack: &items.Backpack{
			Capacity:  dto.Backpack.Capacity,
			ItemsNum:  dto.Backpack.ItemsNum,
			Elixirs:   SliceToList(dto.Backpack.Elixirs),
			Scrolls:   SliceToList(dto.Backpack.Scrolls),
			Foods:     SliceToList(dto.Backpack.Foods),
			Weapons:   SliceToList(dto.Backpack.Weapons),
			Treasures: dto.Backpack.Treasures,
		},
		Weapon:        dto.Weapon,
		ViewRadius:    dto.ViewRadius,
		EnemiesKilled: dto.EnemiesKilled,
		FoodEaten:     dto.FoodEaten,
		ElixirsDrunk:  dto.ElixirsDrunk,
		ScrollsRead:   dto.ScrollsRead,
		HitsDealt:     dto.HitsDealt,
		HitsMissed:    dto.HitsMissed,
		CellsMoved:    dto.CellsMoved,
	}
}

func MapListToSlice[T any](l *list.List) []T {
	if l == nil {
		return nil
	}

	result := make([]T, 0, l.Len())

	for e := l.Front(); e != nil; e = e.Next() {
		if val, ok := e.Value.(T); ok {
			result = append(result, val)
		}
	}

	return result
}

func SliceToList[T any](slice []T) *list.List {
	l := list.New()
	for _, item := range slice {
		l.PushBack(item)
	}
	return l
}
