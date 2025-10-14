package entity

import (
	"errors"
	"math/rand"
)

const (
	RoomsCount    = 9
	RoomMinWidth  = 3
	RoomMinHeight = 3
	RoomWalls     = 2
)

type Level struct {
	Rooms       []Room
	Passages    []Passage
	LevelNumber uint
	LevelEnd    Box
}

func (l *Level) GenerateRoomsOnLevel(sizeMapWidth int, sizeMapHeight int) error {
	l.Rooms = make([]Room, RoomsCount)

	indices := rand.Perm(RoomsCount)
	startIndex := indices[0]
	finishIndex := indices[1]
	var roomType RoomType

	sizeSectionWidth := sizeMapWidth / 3
	sizeSectionHeight := sizeMapHeight / 3

	if sizeSectionWidth < RoomMinWidth || sizeSectionHeight < RoomMinHeight {
		return errors.New("map size is too small: each room must be at least 3x3")
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch i*3 + j {
			case startIndex:
				roomType = RoomTypeStart
			case finishIndex:
				roomType = RoomTypeFinish
			default:
				roomType = RoomTypeOrdinary
			}

			width := uint(rand.Intn(sizeSectionWidth-RoomMinWidth+1) + RoomMinWidth)
			height := uint(rand.Intn(sizeSectionHeight-RoomMinHeight+1) + RoomMinHeight)

			x := j*sizeSectionWidth + rand.Intn(int(width))
			y := i*sizeSectionHeight + rand.Intn(int(height))

			if roomType == RoomTypeFinish {
				l.LevelEnd = Box{
					Point: Point2D[int]{X: x + 1 + rand.Intn(int(width-RoomWalls)), Y: y + 1 + rand.Intn(int(height-RoomWalls))},
					Size:  Size2D[uint]{Height: 1, Width: 1},
				}
			}

			l.Rooms[i*3+j] = *NewRoom(
				roomType,
				Box{
					Point: Point2D[int]{
						X: x,
						Y: y,
					},
					Size: Size2D[uint]{
						Width:  width,
						Height: height,
					},
				})
		}
	}
	return nil
}
