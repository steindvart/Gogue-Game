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

func NewWeapon(box primitives.Box, t WeaponType, e *primitives.Effect) *Weapon {
	return &Weapon{
		Item:   &Item{Shape: box, Name: string(t)},
		Effect: e,
		Type:   t,
	}
}

func NewWeaponBuiltin(rnd utils.Randomizer, box primitives.Box, t WeaponType) *Weapon {
	w, _ := NewWeaponByConfig(rnd, box, GetWeaponConfig(t))
	return w
}

func NewWeaponByConfig(rnd utils.Randomizer, box primitives.Box, cfg WeaponConfig) (*Weapon, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewWeapon(box, cfg.Type, &primitives.Effect{Attributes: cfg.GenerateAttributes(rnd)}), nil
}

func (w *Weapon) Use() *primitives.Effect {
	return w.Effect
}
