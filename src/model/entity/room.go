package entity

type Room struct {
	Geometry Object
	Foods    []Food
	Elixirs  []Elixir
	Scrolls  []Scroll
	Weapons  []Weapon
	Enemies  []Enemy
}
