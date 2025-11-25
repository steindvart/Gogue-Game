package entities

import "gogue/internal/model/primitives"

type EnemyType float64

const (
	EnemyTypeZombie EnemyType = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnakeMage
)

type Direction float64

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

type AttributeRate float64

const (
	AttributeRateLow      AttributeRate = 25
	AttributeRateAverage  AttributeRate = 50
	AttributeRateHigh     AttributeRate = 75
	AttributeRateVeryHigh AttributeRate = 100
)

type HostilityRadius float64

const (
	HostilityRadiusLow     HostilityRadius = 2
	HostilityRadiusAverage HostilityRadius = 4
	HostilityRadiusHigh    HostilityRadius = 6
)

type Enemy struct {
	*Character
	Type            EnemyType
	HostilityRadius HostilityRadius
	IsChasing       bool
	Direction       Direction
}

type Zombie struct {
	Enemy Enemy
}

type Vampire struct {
	Enemy            Enemy
	AbsoluteEvasions float64
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

func NewZombie(box *primitives.Box) *Zombie {
	return &Zombie{
		Enemy: Enemy{
			Character: &Character{
				Shape: box,
				Attributes: primitives.Attributes{
					Agility:   float64(AttributeRateLow),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			},
			Type:            EnemyTypeZombie,
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
	}
}

func NewVampire(box *primitives.Box) *Vampire {
	return &Vampire{
		Enemy: Enemy{
			Character: &Character{
				Shape: box,
				Attributes: primitives.Attributes{
					Agility:   float64(AttributeRateHigh),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			},
			Type:            EnemyTypeVampire,
			HostilityRadius: HostilityRadiusHigh,
			Direction:       DirectionStop,
		},
		AbsoluteEvasions: 0,
	}
}

func NewGhost(box *primitives.Box) *Ghost {
	return &Ghost{
		Enemy: Enemy{
			Character: &Character{
				Shape: box,
				Attributes: primitives.Attributes{
					Agility:   float64(AttributeRateHigh),
					Strength:  float64(AttributeRateLow),
					Health:    float64(AttributeRateLow),
					MaxHealth: float64(AttributeRateLow),
				},
			},
			Type:            EnemyTypeGhost,
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
		IsVisible: true,
	}
}

func NewOgre(box *primitives.Box) *Ogre {
	return &Ogre{
		Enemy: Enemy{
			Character: &Character{
				Shape: box,
				Attributes: primitives.Attributes{
					Agility:   float64(AttributeRateLow),
					Strength:  float64(AttributeRateVeryHigh),
					Health:    float64(AttributeRateVeryHigh),
					MaxHealth: float64(AttributeRateVeryHigh),
				},
			},
			Type:            EnemyTypeOgre,
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
		IsResting: false,
	}
}

func NewSnakeMage(box *primitives.Box) *SnakeMage {
	return &SnakeMage{
		Enemy: Enemy{
			Character: &Character{
				Shape: box,
				Attributes: primitives.Attributes{
					Agility:   float64(AttributeRateVeryHigh),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			},
			Type:            EnemyTypeSnakeMage,
			HostilityRadius: HostilityRadiusHigh,
			Direction:       DirectionStop,
		},
	}
}
