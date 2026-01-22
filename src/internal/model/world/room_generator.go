package world

import (
	"errors"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type GridRoomGenerator struct{}

func NewGridRoomGenerator() *GridRoomGenerator {
	return &GridRoomGenerator{}
}

func (g *GridRoomGenerator) GenerateRooms(config LevelConfig, random utils.Randomizer) ([]Room, primitives.Box, error) {
	sectionSize, err := g.calculateSectionSize(config)
	if err != nil {
		return nil, primitives.Box{}, err
	}

	rooms := make([]Room, config.RoomsCount)
	indexes := random.Perm(config.RoomsCount)
	startRoomIndex := indexes[0]
	finishRoomIndex := indexes[1]

	var finishPortal primitives.Box

	for y := 0; y < config.NumberXYSections; y++ {
		for x := 0; x < config.NumberXYSections; x++ {
			roomIndex := y*config.NumberXYSections + x
			roomType := RoomTypeOrdinary

			switch roomIndex {
			case startRoomIndex:
				roomType = RoomTypeStart
			case finishRoomIndex:
				roomType = RoomTypeFinish
			}

			room, err := g.generateSingleRoom(x, y, roomType, sectionSize, config, random)
			if err != nil {
				return nil, primitives.Box{}, err
			}

			rooms[roomIndex] = room

			if roomType == RoomTypeFinish {
				portalPos, err := room.GetRandomFreePosition(random)
				if err != nil {
					return nil, primitives.Box{}, errors.New("no free positions in finish room to place level portal")
				}
				finishPortal = primitives.Box{
					Point: *portalPos,
					Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
				}
			}
		}
	}

	return rooms, finishPortal, nil
}

func (g *GridRoomGenerator) calculateSectionSize(config LevelConfig) (primitives.Size2D[uint], error) {
	if config.MapSize.Width > config.MapMaxWidth || config.MapSize.Height > config.MapMaxHeight {
		return primitives.Size2D[uint]{}, errors.New("game map size is too big")
	}

	sectionSize := primitives.Size2D[uint]{
		Width:  config.MapSize.Width / uint(config.NumberXYSections),
		Height: config.MapSize.Height / uint(config.NumberXYSections),
	}

	const areaForPassage = 1
	if sectionSize.Width < uint(config.RoomMinWidth+areaForPassage) ||
		sectionSize.Height < uint(config.RoomMinHeight+areaForPassage) {
		return primitives.Size2D[uint]{}, errors.New("game map size is too small")
	}

	return sectionSize, nil
}

func (g *GridRoomGenerator) generateSingleRoom(
	x, y int,
	roomType RoomType,
	sectionSize primitives.Size2D[uint],
	config LevelConfig,
	random utils.Randomizer,
) (Room, error) {
	cellXStart := x*int(sectionSize.Width) + config.MinRoomPadding
	cellYStart := y*int(sectionSize.Height) + config.MinRoomPadding
	cellXEnd := (x + 1) * int(sectionSize.Width)
	cellYEnd := (y + 1) * int(sectionSize.Height)

	maxRoomWidth := cellXEnd - cellXStart
	maxRoomHeight := cellYEnd - cellYStart

	if maxRoomWidth < config.RoomMinWidth || maxRoomHeight < config.RoomMinHeight {
		return Room{}, errors.New("available space in section is smaller than minimum room size")
	}

	width := random.Intn(maxRoomWidth-config.RoomMinWidth+1) + config.RoomMinWidth
	height := random.Intn(maxRoomHeight-config.RoomMinHeight+1) + config.RoomMinHeight

	xCell := cellXStart + random.Intn(maxRoomWidth-width+1)
	yCell := cellYStart + random.Intn(maxRoomHeight-height+1)

	roomBox := primitives.Box{
		Point: primitives.Point2D[int]{X: xCell, Y: yCell},
		Size:  primitives.Size2D[uint]{Width: uint(width), Height: uint(height)},
	}
	room, err := NewRoom(roomType, roomBox)
	if err != nil {
		return Room{}, err
	}
	return *room, nil
}
