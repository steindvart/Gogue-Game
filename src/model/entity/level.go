package entity

import (
	"errors"
	"math/rand"
)

const (
	RoomsCount     = 9
	RoomMinWidth   = 3
	RoomMinHeight  = 3
	RoomWalls      = 2
	MinRoomPadding = 1
)

type Level struct {
	Rooms       []Room
	Passages    []Passage
	LevelNumber uint
	LevelEnd    Box
}

func calculateRoomSectionSize(mapWidth, mapHeight int) (sectionW, sectionH int, err error) {
	totalPaddingW := (3 - 1) * MinRoomPadding * 2
	totalPaddingH := (3 - 1) * MinRoomPadding * 2

	availableW := mapWidth - totalPaddingW
	availableH := mapHeight - totalPaddingH

	sectionW = availableW / 3
	sectionH = availableH / 3

	if sectionW < RoomMinWidth || sectionH < RoomMinHeight {
		return 0, 0, errors.New("map size is too small: available space per section is smaller than minimum room size")
	}

	return sectionW, sectionH, nil
}

func (l *Level) GenerateRoomsOnLevel(sizeMapWidth int, sizeMapHeight int) error {
	l.Rooms = make([]Room, RoomsCount)

	indices := rand.Perm(RoomsCount)
	startIndex := indices[0]
	finishIndex := indices[1]

	sectionW, sectionH, err := calculateRoomSectionSize(sizeMapWidth, sizeMapHeight)
	if err != nil {
		return err
	}

	if sectionW < RoomMinWidth || sectionH < RoomMinHeight {
		return errors.New("map size is too small: each room section must be at least min room size")
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			roomIndex := i*3 + j
			var roomType RoomType

			switch roomIndex {
			case startIndex:
				roomType = RoomTypeStart
			case finishIndex:
				roomType = RoomTypeFinish
			default:
				roomType = RoomTypeOrdinary
			}

			cellXStart := j*sectionW + MinRoomPadding
			cellYStart := i*sectionH + MinRoomPadding
			cellXEnd := (j+1)*sectionW - MinRoomPadding
			cellYEnd := (i+1)*sectionH - MinRoomPadding

			maxRoomWidth := cellXEnd - cellXStart
			maxRoomHeight := cellYEnd - cellYStart

			if maxRoomWidth < RoomMinWidth || maxRoomHeight < RoomMinHeight {
				return errors.New("internal error: available space in section is smaller than minimum room size")
			}

			width := rand.Intn(maxRoomWidth-RoomMinWidth+1) + RoomMinWidth
			height := rand.Intn(maxRoomHeight-RoomMinHeight+1) + RoomMinHeight

			x := cellXStart + rand.Intn(maxRoomWidth-width+1)
			y := cellYStart + rand.Intn(maxRoomHeight-height+1)

			if roomType == RoomTypeFinish {
				if width < RoomWalls || height < RoomWalls {
					return errors.New("finish room is too small to place LevelEnd")
				}
				l.LevelEnd = Box{
					Point: Point2D[int]{
						X: x + 1 + rand.Intn(width-RoomWalls),
						Y: y + 1 + rand.Intn(height-RoomWalls),
					},
					Size: Size2D[uint]{Height: 1, Width: 1},
				}
			}

			roomBox := Box{
				Point: Point2D[int]{X: x, Y: y},
				Size:  Size2D[uint]{Width: uint(width), Height: uint(height)},
			}

			l.Rooms[roomIndex] = *NewRoom(roomType, roomBox)
		}
	}

	return nil
}
