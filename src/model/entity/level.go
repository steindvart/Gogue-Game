package entity

import (
	"math/rand"
)

const (
	RoomsCount    = 9
	RoomMinWidth  = 3
	RoomMinHeight = 3
	RoomWalls     = 2
)

const (
	LEVEL_HEIGHT = 100
	LEVEL_WIDTH  = 150
)

type Level struct {
	Rooms       []Room
	Passages    []Passage
	LevelNumber uint
	LevelEnd    Box
}

func (l *Level) GenerateRoomsOnLevel() {
	l.Rooms = make([]Room, RoomsCount)

	indices := rand.Perm(RoomsCount)
	startIndex := indices[0]
	finishIndex := indices[1]
	var roomType RoomType

	sizeSectionWidth := LEVEL_WIDTH / 3
	sizeSectionHeight := LEVEL_HEIGHT / 3

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
}
