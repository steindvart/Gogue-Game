package world

import (
	"errors"
	"gogue/internal/model/primitives"
	"math/rand"
	"reflect"
	"testing"
)

func TestRoom_NewRoom(t *testing.T) {
	tests := []struct {
		name     string
		roomType RoomType
		mapSize  primitives.Box
		want     Room
	}{
		{
			name:     "New finish room",
			roomType: RoomTypeFinish,
			mapSize:  primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Box:  &primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type: RoomTypeFinish,
			},
		},
		{
			name:     "New ordinary room",
			roomType: RoomTypeOrdinary,
			mapSize:  primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Box:  &primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type: RoomTypeOrdinary,
			},
		},
		{
			name:     "New start room",
			roomType: RoomTypeStart,
			mapSize:  primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			want: Room{
				Box:  &primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
				Type: RoomTypeStart,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := NewRoom(tt.roomType, tt.mapSize)
			reflect.DeepEqual(room, tt.want)
		})
	}
}

func TestRoom_GetRandomFreePosition(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	tests := []struct {
		name              string
		roomType          RoomType
		roomSize          primitives.Box
		OccupiedPositions map[primitives.Point2D[int]]bool
		wantPos           primitives.Point2D[int]
		wantErr           error
	}{
		{
			name:              "Base Room",
			roomType:          RoomTypeFinish,
			roomSize:          primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			OccupiedPositions: map[primitives.Point2D[int]]bool{},
			wantPos:           primitives.Point2D[int]{X: 8, Y: 7},
			wantErr:           nil,
		},
		{
			name:     "Room With Occupied Positions",
			roomType: RoomTypeOrdinary,
			roomSize: primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 10, Width: 10}},
			OccupiedPositions: map[primitives.Point2D[int]]bool{
				{X: 12, Y: 7}:  true,
				{X: 11, Y: 11}: true,
				{X: 13, Y: 12}: true,
			},
			wantPos: primitives.Point2D[int]{X: 11, Y: 10},
			wantErr: nil,
		},
		{
			name:     "Absence Free Positions In Room",
			roomType: RoomTypeOrdinary,
			roomSize: primitives.Box{Point: primitives.Point2D[int]{X: 6, Y: 3}, Size: primitives.Size2D[uint]{Height: 3, Width: 3}},
			OccupiedPositions: map[primitives.Point2D[int]]bool{
				{X: 8, Y: 5}: true,
				{X: 7, Y: 4}: true,
				{X: 8, Y: 4}: true,
				{X: 7, Y: 5}: true,
			},
			wantErr: errors.New("no available positions in room"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := NewRoom(tt.roomType, tt.roomSize)
			for ocPos := range tt.OccupiedPositions {
				room.occupiedPositions[ocPos] = true
			}
			gotPoint, err := room.GetRandomFreePosition(source)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("addDoorsAtRoom() expected error '%v', but got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("addDoorsAtRoom() expected error '%v', but got '%v'", tt.wantErr, err)
					return
				}
			} else {
				if err != nil {
					t.Errorf("GetRandomFreePosition() expected no error, but got '%v'", err)
					return
				}
				if *gotPoint != tt.wantPos {
					t.Errorf("Expected free position %+v, got %+v", tt.wantPos, gotPoint)
				}
			}
		})
	}
}
