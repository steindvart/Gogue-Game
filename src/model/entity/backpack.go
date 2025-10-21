package entity

const (
	BackpackDefaultCapacity uint = 9
)

type Backpack struct {
	Capacity    uint
	ItemsNum    uint
	Consumables []ConsumableLike
	// Foods     []Food
	// Elixirs   []Elixir
	// Scrolls   []Scroll
	// Weapon    []Weapon
	Treasures   uint
}

func NewBackpack() *Backpack {
	return &Backpack{
		Capacity:    BackpackDefaultCapacity,
		ItemsNum:    0,
		Consumables: []ConsumableLike{},
		// Foods:     []Food{},
		// Elixirs:   []Elixir{},
		// Scrolls:   []Scroll{},
		// Weapon:    []Weapon{},
		Treasures: 0,
	}
}
