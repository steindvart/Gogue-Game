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
	RoomMinHeight               = 4
	RoomMinWidth                = 4
	RoomsCountEntireLevelHeight = 3
	RoomsCountEntireLevelWidth  = 3
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
	Portal  Portal
}

func NewRoom(roomType RoomType, x, y int, mapHeight, mapWidth int, rand *rand.Rand) (*Room, error) {
	width, height, err := calculateRoomSize(mapHeight, mapWidth, rand)
	if err != nil {
		return nil, err
	}

	shape := Box{
		Point: Point2D[int]{X: x, Y: y},
		Size:  Size2D[uint]{Height: height, Width: width},
	}

	return &Room{
		Shape:   shape,
		Type:    roomType,
		Foods:   []Food{},
		Elixirs: []Elixir{},
		Scrolls: []Scroll{},
		Weapons: []Weapon{},
		Enemies: []Enemy{},
		Portal:  Portal{},
	}, nil
}

func calculateRoomSize(mapHeight, mapWidth int, rand *rand.Rand) (uint, uint, error) {
	roomMaxHeight := mapHeight / RoomsCountEntireLevelHeight
	roomMaxWidth := mapWidth / RoomsCountEntireLevelWidth

	if roomMaxWidth < RoomMinWidth || roomMaxHeight < RoomMinHeight {
		return 0, 0, errors.New("room cannot be built: map is too small to fit a room of minimum required size")
	}

	width := uint(rand.Intn(roomMaxWidth-RoomMinWidth+1) + RoomMinWidth)
	height := uint(rand.Intn(roomMaxHeight-RoomMinHeight+1) + RoomMinHeight)

	return width, height, nil
}
