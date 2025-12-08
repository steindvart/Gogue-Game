package common

type GameEntityType int

const (
	EntityTypePlayer GameEntityType = iota + 1
	WorldTypeWall
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
