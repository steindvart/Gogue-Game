package entity

const (
	BackpackDefaultCapacity uint = 9
)

type Backpack struct {
	Capacity  uint
	ItemsNum  uint
	Foods     []Food
	Elixirs   []Elixir
	Scrolls   []Scroll
	Weapon    []Weapon
	Treasures []Treasure
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:  BackpackDefaultCapacity,
		ItemsNum:  0,
		Foods:     []Food{},
		Elixirs:   []Elixir{},
		Scrolls:   []Scroll{},
		Weapon:    []Weapon{},
		Treasures: []Treasure{},
	}
}
