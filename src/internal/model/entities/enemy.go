package entities

import "gogue/internal/model/primitives"

type EnemyType int

const (
	EnemyTypeZombie EnemyType = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnakeMage
)

var EnemyTypes = []EnemyType{
	EnemyTypeZombie,
	EnemyTypeVampire,
	EnemyTypeGhost,
	EnemyTypeOgre,
	EnemyTypeSnakeMage,
}

type Direction int

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

type AttributeRate int

const (
	AttributeRateLow      AttributeRate = 25
	AttributeRateAverage  AttributeRate = 50
	AttributeRateHigh     AttributeRate = 75
	AttributeRateVeryHigh AttributeRate = 100
)

type HostilityRadius int

const (
	HostilityRadiusLow     HostilityRadius = 2
	HostilityRadiusAverage HostilityRadius = 4
	HostilityRadiusHigh    HostilityRadius = 6
)

type EnemyProvider interface {
	GetEnemy() *Enemy
}

type Enemy struct {
	*Character
	HostilityRadius
	Direction
	IsChasing bool
}

type Zombie struct {
	*Enemy
}

type Vampire struct {
	*Enemy
	AbsoluteEvasions uint
}

type Ghost struct {
	*Enemy
	IsVisible bool
}

type Ogre struct {
	*Enemy
	IsResting bool
}

type SnakeMage struct {
	*Enemy
}

func NewZombie(box *primitives.Box) *Zombie {
	return &Zombie{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   float64(AttributeRateLow),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			},
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
	}
}

func NewVampire(box *primitives.Box) *Vampire {
	return &Vampire{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   float64(AttributeRateHigh),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			},
			HostilityRadius: HostilityRadiusHigh,
			Direction:       DirectionStop,
		},
		AbsoluteEvasions: 0,
	}
}

func NewGhost(box *primitives.Box) *Ghost {
	return &Ghost{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   float64(AttributeRateHigh),
					Strength:  float64(AttributeRateLow),
					Health:    float64(AttributeRateLow),
					MaxHealth: float64(AttributeRateLow),
				},
			},
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
		IsVisible: true,
	}
}

func NewOgre(box *primitives.Box) *Ogre {
	return &Ogre{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   float64(AttributeRateLow),
					Strength:  float64(AttributeRateVeryHigh),
					Health:    float64(AttributeRateVeryHigh),
					MaxHealth: float64(AttributeRateVeryHigh),
				},
			},
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
		IsResting: false,
	}
}

func NewSnakeMage(box *primitives.Box) *SnakeMage {
	return &SnakeMage{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   float64(AttributeRateVeryHigh),
					Strength:  float64(AttributeRateAverage),
					Health:    float64(AttributeRateHigh),
					MaxHealth: float64(AttributeRateHigh),
				},
			}, HostilityRadius: HostilityRadiusHigh,
			Direction: DirectionStop,
		},
	}
}

func (e *Enemy) GetEnemy() *Enemy {
	return e
}
