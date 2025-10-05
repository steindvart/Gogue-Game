package entity

type Room struct {
	Shape   Box
	Foods   []Food
	Elixirs []Elixir
	Scrolls []Scroll
	Weapons []Weapon
	Enemies []Enemy
}
