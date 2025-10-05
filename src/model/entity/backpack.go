package entity

type Backpack struct {
	Capacity       uint
	ConsumablesNum uint
	Foods          []Food
	Elixirs        []Elixir
	Scrolls        []Scroll
	Weapon         []Weapon
	Treasures      []Treasure
}
