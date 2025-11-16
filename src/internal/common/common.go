package common

type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
	EntityTypeHorizontalWall
	EntityTypeVerticalWall
	EntityTypePortal
	EntityTypePassage
	EntityTypeDoorOne
	EntityTypeDoorTwo
	EntityTypeZombie
	EntityTypeVampire
	EntityTypeGhost
	EntityTypeOgre
	EntityTypeSnakeMage
)