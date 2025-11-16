package items

import (
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Weapon struct {
	*Item
	Type               WeaponType
	AffectedAttributes primitives.Attributes
}

func NewWeapon(rnd *utils.RandomGenerator, box primitives.Box, t WeaponType) *Weapon {
	cfg := GetWeaponConfig(t)

	return &Weapon{
		Item: &Item{
			Shape: box,
			Name:  string(cfg.Type),
		},
		Type:               t,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
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
		Type:               WeaponTypeCustom,
		AffectedAttributes: cfg.GenerateAttributes(rnd),
	}, nil
}

func (w *Weapon) Use() primitives.Attributes {
	return w.AffectedAttributes
}

func AsWeapon(item any) *Weapon {
	w, ok := item.(*Weapon)
	if ok {
		return w
	}
	return nil
}
