package entity

type Backpack struct {
	Capacity  uint
	ItemsNum  uint
	Foods     []Food
	Elixirs   []Elixir
	Scrolls   []Scroll
	Weapon    []Weapon
	Treasures []Treasure
}
