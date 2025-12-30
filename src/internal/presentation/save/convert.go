package save

import (
	"container/list"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
	"gogue/internal/presentation/dto"
	"gogue/internal/utils"
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
		Position:   level.Player.GetPosition(),
		MaxHealth:  level.Player.MaxHealth,
		Health:     level.Player.Health,
		Agility:    level.Player.Agility,
		Strength:   level.Player.Strength,
		ViewRadius: level.Player.ViewRadius,
	}

	return dto.GameSaveDto{
		LevelNumber:  level.Number,
		Size:         size,
		Rooms:        roomsDto,
		Passages:     level.Passages,
		FinishPortal: level.FinishPortal,
		FogOfWar:     fogOfWarDto,
		Player:       playerDto,
	}
}

func RestoreLevelFromDto(dto dto.GameSaveDto, random utils.Randomizer) (*world.Level, error) {
	level := world.NewLevelWithDefaults(random, primitives.Size2D[uint]{
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
		TemporaryEffects: nil,
	}

	backpack := &items.Backpack{
		Capacity:  10,
		ItemsNum:  0,
		Elixirs:   list.New(),
		Scrolls:   list.New(),
		Foods:     list.New(),
		Weapons:   list.New(),
		Treasures: 0,
	}

	var weapon *items.Weapon = nil

	return &entities.Player{
		Character:  character,
		Backpack:   backpack,
		Weapon:     weapon,
		Experience: dto.Experience,
		Level:      dto.Level,
		ViewRadius: dto.ViewRadius,
	}
}
