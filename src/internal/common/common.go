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
	FoodTypePotatoes
	FoodTypeBread
	FoodTypeMeat
	FoodTypeMistery
	FoodTypeBeer
	ElixirTypeStrength
	ElixirTypeAgility
	ElixirTypeDwarfism
	ElixirTypeGiantism
	ElixirTypeMystery
	ScrollTypeStrength
	ScrollTypeAgility
	ScrollTypeUltimate
	ScrollTypeMaxHealth
	ScrollTypeMystery
	Weapon
)
