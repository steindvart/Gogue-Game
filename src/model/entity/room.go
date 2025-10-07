package entity

import (
	"errors"
	"math/rand"
)

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

const (
	RoomMinHeight = 4
	RoomMinWidth  = 4
	RoomWalls     = 2
)

type TileType uint

type Room struct {
	Shape   Box
	Type    RoomType
	Foods   []Food
	Elixirs []Elixir
	Scrolls []Scroll
	Weapons []Weapon
	Enemies []Enemy
	Portal  *Portal
}

func NewRoom(roomType RoomType, minX int, minY int, maxWidth int, maxHeight int) (*Room, error) {
	if maxWidth < RoomMinHeight || maxHeight < RoomMinHeight {
		return nil, errors.New("room cannot be built: map is too small to fit a room of minimum required size")
	}

	width, height, err := calculateRoomSize(maxWidth, maxHeight)
	if err != nil {
		return nil, err
	}

	shape := Box{
		Point: Point2D[int]{X: minX + rand.Intn(int(width)), Y: minY + rand.Intn(int(height))},
		Size:  Size2D[uint]{Height: height, Width: width},
	}

	var portal *Portal = nil
	if roomType == RoomTypeFinish || roomType == RoomTypeStart {
		portal = generatePortal(shape.Point.X, shape.Point.Y, int(width), int(height))
	}

	return &Room{
		Shape:   shape,
		Type:    roomType,
		Foods:   []Food{},
		Elixirs: []Elixir{},
		Scrolls: []Scroll{},
		Weapons: []Weapon{},
		Enemies: []Enemy{},
		Portal:  portal,
	}, nil
}

func calculateRoomSize(maxWidth int, maxHeight int) (uint, uint, error) {
	width := uint(rand.Intn(maxWidth-RoomMinWidth+1) + RoomMinWidth)
	height := uint(rand.Intn(maxHeight-RoomMinHeight+1) + RoomMinHeight)

	return width, height, nil
}

func generatePortal(x int, y int, width int, height int) *Portal {
	portal := &Portal{
		Shape: Box{
			Point: Point2D[int]{X: x + 1 + rand.Intn(width-RoomWalls), Y: y + 1 + rand.Intn(height-RoomWalls)},
			Size:  Size2D[uint]{Height: 1, Width: 1},
		},
	}
	return portal
}
