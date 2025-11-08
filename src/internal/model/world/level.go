package world

import (
	"errors"
	"gogue/internal/model/primitive"
	"math/rand"
)

const (
	RoomMinWidth     = 3
	RoomMinHeight    = 3
	MinRoomPadding   = 1
	numberXYSections = 3
)

type Level struct {
	Rooms    []Room
	Passages []Passage
	Number   uint
	End      primitive.Box
}

func calculateRoomSectionSize(mapSize primitive.Size2D[uint]) (sectionSize primitive.Size2D[uint]) {
	totalPaddingWidth := uint(MinRoomPadding * 2)
	totalPaddingHeight := uint(MinRoomPadding * 2)

	availableWidth := mapSize.Width - totalPaddingWidth
	availableHeight := mapSize.Height - totalPaddingHeight

	sectionSize.Width = availableWidth / numberXYSections
	sectionSize.Height = availableHeight / numberXYSections

	return sectionSize
}

func (l *Level) GenerateNineRooms(sizeMap primitive.Size2D[uint]) error {
	sectionSize := calculateRoomSectionSize(sizeMap)
	if sectionSize.Width < RoomMinWidth || sectionSize.Height < RoomMinHeight {
		return errors.New("map size is too small: each room section must be at least min room size")
	}

	const roomsCount = 9
	l.Rooms = make([]Room, roomsCount)
	// rand.Perm(9) возвращает массив перемешанных чисел от 0 до 8, чтобы далее не было повторений index для Start и Finish
	indexes := rand.Perm(roomsCount)
	startRoomIndex := indexes[0]
	finishRoomIndex := indexes[1]
	for y := 0; y < numberXYSections; y++ {
		for x := 0; x < numberXYSections; x++ {
			roomIndex := y*numberXYSections + x
			var roomType RoomType
			switch roomIndex {
			case startRoomIndex:
				roomType = RoomTypeStart
			case finishRoomIndex:
				roomType = RoomTypeFinish
			default:
				roomType = RoomTypeOrdinary
			}

			cellXStart := x*int(sectionSize.Width) + MinRoomPadding
			cellYStart := y*int(sectionSize.Height) + MinRoomPadding
			cellXEnd := (x+1)*int(sectionSize.Width) - MinRoomPadding
			cellYEnd := (y+1)*int(sectionSize.Height) - MinRoomPadding

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
				roomWall := 1
				l.End = primitive.Box{
					Point: primitive.Point2D[int]{
						X: x + roomWall + rand.Intn(width-(roomWall*2)),
						Y: y + roomWall + rand.Intn(height-(roomWall*2)),
					},
					Size: primitive.Size2D[uint]{Height: 1, Width: 1},
				}
			}

			roomBox := primitive.Box{
				Point: primitive.Point2D[int]{X: x, Y: y},
				Size:  primitive.Size2D[uint]{Width: uint(width), Height: uint(height)},
			}

			l.Rooms[roomIndex] = *NewRoom(roomType, roomBox)
		}
	}

	return nil
}
