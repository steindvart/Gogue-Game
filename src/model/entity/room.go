package entity

type RoomType uint

const (
	RoomTypeStart RoomType = iota
	RoomTypeOrdinary
	RoomTypeFinish
)

type Room struct {
	Shape   Box
	Type    RoomType
	Foods   []Food
	Elixirs []Elixir
	Scrolls []Scroll
	Weapons []Weapon
	Enemies []Enemy
	Portal  *Box
}

func NewRoom(roomType RoomType, shape Box, portal *Box) Room {
	return Room{
		Shape:   shape,
		Type:    roomType,
		Foods:   []Food{},
		Elixirs: []Elixir{},
		Scrolls: []Scroll{},
		Weapons: []Weapon{},
		Enemies: []Enemy{},
		Portal:  portal,
	}
}
