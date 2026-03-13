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
	AttributeRateLow      AttributeRate = 5
	AttributeRateAverage  AttributeRate = 10
	AttributeRateHigh     AttributeRate = 15
	AttributeRateVeryHigh AttributeRate = 20
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
	IsChasing      bool
	StepsRemaining int // Сколько шагов осталось в текущем Direction (idle-патрулирование)
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

// Zombie — медленный, туповатый, но живучий. Низкий урон, много HP.
// Роль: «мешок с хитпоинтами», тренировочный враг для начала игры.
func NewZombie(box *primitives.Box) *Zombie {
	return &Zombie{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   3,
					Strength:  8,
					Health:    35,
					MaxHealth: 35,
				},
			},
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
	}
}

// Vampire — быстрый, ловкий, средний урон. Часто уклоняется.
// Роль: «ловкач», сложно попасть, но и бьёт не так больно.
func NewVampire(box *primitives.Box) *Vampire {
	return &Vampire{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   14,
					Strength:  10,
					Health:    30,
					MaxHealth: 30,
				},
			},
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
		AbsoluteEvasions: 0,
	}
}

// Ghost — хрупкий, слабый, но невидимый до агрессии. Лёгкий враг.
// Роль: «неожиданность», пугает, но быстро убивается.
func NewGhost(box *primitives.Box) *Ghost {
	return &Ghost{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   10,
					Strength:  6,
					Health:    15,
					MaxHealth: 15,
				},
			},
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
		IsVisible: true,
	}
}

// Ogre — танк. Очень много HP, сильно бьёт, но медленный и неповоротливый.
// Роль: «мини-босс», требует подготовки (оружие/эликсиры).
func NewOgre(box *primitives.Box) *Ogre {
	return &Ogre{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   2,
					Strength:  18,
					Health:    60,
					MaxHealth: 60,
				},
			},
			HostilityRadius: HostilityRadiusLow,
			Direction:       DirectionStop,
		},
		IsResting: false,
	}
}

// SnakeMage — стеклянная пушка. Высокий урон, ловкий, но хрупкий.
// Роль: «приоритетная цель», нужно убивать быстро или избегать.
func NewSnakeMage(box *primitives.Box) *SnakeMage {
	return &SnakeMage{
		Enemy: &Enemy{
			Character: &Character{
				Box: box,
				Attributes: &primitives.Attributes{
					Agility:   12,
					Strength:  14,
					Health:    22,
					MaxHealth: 22,
				},
			},
			HostilityRadius: HostilityRadiusAverage,
			Direction:       DirectionStop,
		},
	}
}

func (e *Enemy) GetEnemy() *Enemy {
	return e
}

// ScaleAttributes масштабирует атрибуты врага на указанный множитель.
// Используется для повышения сложности на поздних уровнях подземелья.
// multiplier = 1.0 означает базовые значения, 1.5 = +50% ко всем атрибутам.
func (e *Enemy) ScaleAttributes(multiplier float64) {
	if multiplier <= 0 {
		return
	}

	attrs := e.Character.Attributes
	attrs.Health *= multiplier
	attrs.MaxHealth *= multiplier
	attrs.Strength *= multiplier
	// Agility масштабируется с половинным коэффициентом,
	// чтобы не делать врагов неуязвимыми за счёт уклонения.
	attrs.Agility *= 1.0 + (multiplier-1.0)*0.5
}
