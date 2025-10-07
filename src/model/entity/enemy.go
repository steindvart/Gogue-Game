package entity

type EnemyType uint

const (
	EnemyTypeZombie EnemyType = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnakeMage
)

type Direction uint

const (
	DirectionForward Direction = iota
	DirectionBack
	DirectionLeft
	DirectionRight
	DirectionDiagonallyForwardLeft
	DirectionDiagonallyForwardRight
	DirectionDiagonallyBackLeft
	DirectionDiagonallyBackRight
	DirectionStop
)

type AttributeRate uint

const (
	AttributeRateLow      AttributeRate = 25
	AttributeRateAverage  AttributeRate = 50
	AttributeRateHigh     AttributeRate = 75
	AttributeRateVeryHigh AttributeRate = 100
)

type HostilityRadius uint

const (
	HostilityRadiusLow     HostilityRadius = 2
	HostilityRadiusAverage HostilityRadius = 4
	HostilityRadiusHigh    HostilityRadius = 6
)

type Enemy struct {
	Character       Character
	Type            EnemyType
	HostilityRadius HostilityRadius
	IsChasing       bool
	Direction       Direction
}

type Zombie struct {
	Enemy Enemy
}

type Vampire struct {
	Enemy          Enemy
	HadFirstDamage bool
}

type Ghost struct {
	Enemy     Enemy
	IsVisible bool
}

type Ogre struct {
	Enemy     Enemy
	IsResting bool
}

type SnakeMage struct {
	Enemy Enemy
}

func NewZombie(box Box) *Zombie {
	enemy := Zombie{
		Enemy: Enemy{
			Character: Character{
				Shape:    box,
				Agility:  uint(AttributeRateLow),
				Strength: uint(AttributeRateAverage),
				Health:   float64(AttributeRateHigh),
			},
			Type:            EnemyTypeZombie,
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
	}

	return &enemy
}

func NewVampire(box Box) *Vampire {
	enemy := Vampire{
		Enemy: Enemy{
			Character: Character{
				Shape:    box,
				Agility:  uint(AttributeRateHigh),
				Strength: uint(AttributeRateAverage),
				Health:   float64(AttributeRateHigh),
			},
			Type:            EnemyTypeVampire,
			HostilityRadius: HostilityRadiusHigh,
			Direction:       DirectionStop,
		},
		HadFirstDamage: false,
	}

	return &enemy
}

func NewGhost(box Box) *Ghost {
	enemy := Ghost{
		Enemy: Enemy{
			Character: Character{
				Shape:    box,
				Agility:  uint(AttributeRateHigh),
				Strength: uint(AttributeRateLow),
				Health:   float64(AttributeRateLow),
			},
			Type:            EnemyTypeGhost,
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
		IsVisible: true,
	}

	return &enemy
}

func NewOgre(box Box) *Ogre {
	enemy := Ogre{
		Enemy: Enemy{
			Character: Character{
				Shape:    box,
				Agility:  uint(AttributeRateLow),
				Strength: uint(AttributeRateVeryHigh),
				Health:   float64(AttributeRateVeryHigh),
			},
			Type:            EnemyTypeOgre,
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
		IsResting: false,
	}

	return &enemy
}

func NewSnakeMage(box Box) *SnakeMage {
	enemy := SnakeMage{
		Enemy: Enemy{
			Character: Character{
				Shape:    box,
				Agility:  uint(AttributeRateVeryHigh),
				Strength: uint(AttributeRateAverage),
				Health:   float64(AttributeRateHigh),
			},
			Type:            EnemyTypeSnakeMage,
			HostilityRadius: HostilityRadiusHigh,
			Direction:       DirectionStop,
		},
	}

	return &enemy
}
