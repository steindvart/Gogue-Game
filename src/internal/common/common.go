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

type RenderedField struct {
	EnvironmentLayer [][]GameEntityType
	ObjectLayer      [][]GameEntityType
	Width            int
	Height           int
}

func NewRenderedField(width, height int) *RenderedField {
	env := make([][]GameEntityType, height)
	obj := make([][]GameEntityType, height)
	for y := 0; y < height; y++ {
		env[y] = make([]GameEntityType, width)
		obj[y] = make([]GameEntityType, width)
	}
	return &RenderedField{
		EnvironmentLayer: env,
		ObjectLayer:      obj,
		Width:            width,
		Height:           height,
	}
}

// Compose объединяет два слоя в итоговое однослойное представление.
// Объекты имеют приоритет над окружением (если объект есть - берём его, иначе - окружение).
func (rf *RenderedField) Compose() [][]GameEntityType {
	result := make([][]GameEntityType, rf.Height)
	for y := 0; y < rf.Height; y++ {
		result[y] = make([]GameEntityType, rf.Width)
		for x := 0; x < rf.Width; x++ {
			if rf.ObjectLayer[y][x] != EntityTypeNone {
				result[y][x] = rf.ObjectLayer[y][x]
			} else {
				result[y][x] = rf.EnvironmentLayer[y][x]
			}
		}
	}
	return result
}
