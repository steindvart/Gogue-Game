package common

type GameEntityType int

const (
	EntityTypeNone GameEntityType = iota
	EntityTypePlayer
	WorldTypeWall
	WorldTypeRoomFloor
	WorldTypePortal
	WorldTypePassage
	WorldTypeDoor
	EntityTypeZombie
	EntityTypeVampire
	EntityTypeGhost
	EntityTypeOgre
	EntityTypeSnakeMage
	EntityTypeMimic
	Food
	Elixir
	Scroll
	Weapon
)
