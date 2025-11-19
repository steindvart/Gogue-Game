package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Weapon struct {
	*Item
	*primitives.Effect
	Type WeaponType
}

func NewWeapon(rnd *utils.RandomGenerator, box primitives.Box, t WeaponType) *Weapon {
	cfg := GetWeaponConfig(t)

	return &Weapon{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Effect: &primitives.Effect{
			Attributes: cfg.GenerateAttributes(rnd),
			Duration:   0,
		},
		Type: t,
	}
}

func NewWeaponByConfig(rnd *utils.RandomGenerator, box primitives.Box, cfg WeaponConfig) (*Weapon, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Weapon{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Effect: &primitives.Effect{
			Attributes: cfg.GenerateAttributes(rnd),
			Duration:   0,
		},
		Type: WeaponTypeCustom,
	}, nil
}

func (w *Weapon) Use() *primitives.Effect {
	return w.Effect
}

func AsWeapon(item any) *Weapon {
	w, ok := item.(*Weapon)
	if ok {
		return w
	}
	return nil
}
