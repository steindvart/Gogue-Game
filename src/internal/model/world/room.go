package world

import (
	"errors"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

type Room struct {
	Shape   primitives.Box
	Type    RoomType
	Doors   []primitives.Point2D[int]
	Foods   []items.Food
	Elixirs []items.Elixir
	Scrolls []items.Scroll
	Weapons []items.Weapon
	Enemies []entities.Enemy

	OccupiedPositions map[primitives.Point2D[int]]bool
}

func NewRoom(roomType RoomType, shape primitives.Box) *Room {
	if shape.Size.Width < RoomMinWidth || shape.Size.Height < RoomMinHeight {
		shape.Size.Width = RoomMinWidth
		shape.Size.Height = RoomMinHeight
	}
	return &Room{
		Shape:             shape,
		Type:              roomType,
		Foods:             []items.Food{},
		Elixirs:           []items.Elixir{},
		Scrolls:           []items.Scroll{},
		Weapons:           []items.Weapon{},
		Enemies:           []entities.Enemy{},
		OccupiedPositions: make(map[primitives.Point2D[int]]bool),
	}
}

func (r *Room) GetRandomFreePosition(rand utils.Randomizer) (*primitives.Point2D[int], error) {
	// Получаем все возможные точки внутри комнаты без границ
	minX := r.Shape.Point.X + 1
	minY := r.Shape.Point.Y + 1
	width := r.Shape.Size.Width - 1
	height := r.Shape.Size.Height - 1

	totalPossiblePoints := int(width) * int(height)

	for attempt := 0; attempt < totalPossiblePoints; attempt++ {
		pos := primitives.Point2D[int]{
			X: minX + rand.Intn(int(width)),
			Y: minY + rand.Intn(int(height)),
		}

		if !r.OccupiedPositions[pos] {
			r.OccupiedPositions[pos] = true
			return &pos, nil
		}
	}

	return nil, errors.New("no available positions in Rooms")
}

func (r *Room) GetCountFreePosition() int {
	return int(r.Shape.Size.Width*r.Shape.Size.Height) - len(r.OccupiedPositions)
}

func (r *Room) createFood(rand utils.Randomizer, pos primitives.Point2D[int], allFoodType []items.FoodType) {
	foodType := getRandomElement(rand, allFoodType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newFood := items.NewFoodBuiltin(rand, itemBox, foodType)
	r.Foods = append(r.Foods, *newFood)
}

func (r *Room) createElixir(rand utils.Randomizer, pos primitives.Point2D[int], allElixirType []items.ElixirType) bool {
	elixirType := getRandomElement(rand, allElixirType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newElixir := items.NewElixirBuiltin(rand, itemBox, elixirType)
	r.Elixirs = append(r.Elixirs, *newElixir)
	return true
}

func (r *Room) createScroll(rand utils.Randomizer, pos primitives.Point2D[int], allScrollType []items.ScrollType) bool {
	scrollType := getRandomElement(rand, allScrollType)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newScroll := items.NewScrollBuiltin(rand, itemBox, scrollType)
	r.Scrolls = append(r.Scrolls, *newScroll)
	return true
}

func (r *Room) createWeapon(rand utils.Randomizer, pos primitives.Point2D[int], allWeaponTypes []items.WeaponType) bool {
	weaponType := getRandomElement(rand, allWeaponTypes)
	itemBox := primitives.Box{
		Point: pos,
		Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
	}

	newWeapon := items.NewWeaponBuiltin(rand, itemBox, weaponType)
	r.Weapons = append(r.Weapons, *newWeapon)
	return true
}

func getRandomElement[T any](random utils.Randomizer, slice []T) T {
	idx := random.Intn(len(slice))
	return slice[idx]
}
