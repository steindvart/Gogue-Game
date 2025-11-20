package common

type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
	EntityTypeWall
	EntityTypePortal
	EntityTypePassage
	EntityTypeDoor
	EntityTypeZombie
	EntityTypeVampire
	EntityTypeGhost
	EntityTypeOgre
	EntityTypeSnakeMage
)
