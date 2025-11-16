package items

import (
	"errors"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type WeaponType string

const (
	WeaponTypeDagger  WeaponType = "Dagger"  // Strength: 3...10, Agility: 5...12
	WeaponTypeSpear   WeaponType = "Spear"   // Strength: 7...18, Agility: 2...8
	WeaponTypeSword   WeaponType = "Sword"   // Strength: 15...20, Agility: -2...5
	WeaponTypeAxe     WeaponType = "Axe"     // Strength: 17...35, Agility: -9...-4
	WeaponTypeMaul    WeaponType = "Maul"    // Strength: 20...45, Agility: -15...-5
	WeaponTypeMystery WeaponType = "Mystery" // Strength: 1...50, Agility: -20...20
)

type WeaponConfig struct {
	Type          WeaponType
	StrengthRange primitives.AttributeRange
	AgilityRange  primitives.AttributeRange
	Description   string
}

var WeaponRegistry = map[WeaponType]WeaponConfig{
	WeaponTypeDagger: {
		Type:          WeaponTypeDagger,
		StrengthRange: primitives.AttributeRange{Min: 3, Max: 10},
		AgilityRange:  primitives.AttributeRange{Min: 5, Max: 12},
		Description:   "A light, quick weapon with high agility bonus but moderate damage",
	},
	WeaponTypeSpear: {
		Type:          WeaponTypeSpear,
		StrengthRange: primitives.AttributeRange{Min: 7, Max: 18},
		AgilityRange:  primitives.AttributeRange{Min: 2, Max: 8},
		Description:   "A reach weapon offering good damage with some agility bonus",
	},
	WeaponTypeSword: {
		Type:          WeaponTypeSword,
		StrengthRange: primitives.AttributeRange{Min: 15, Max: 20},
		AgilityRange:  primitives.AttributeRange{Min: -2, Max: 5},
		Description:   "A balanced weapon combining decent damage with slight agility adjustment",
	},
	WeaponTypeAxe: {
		Type:          WeaponTypeAxe,
		StrengthRange: primitives.AttributeRange{Min: 17, Max: 35},
		AgilityRange:  primitives.AttributeRange{Min: -9, Max: -4},
		Description:   "A heavy weapon with high damage but reduced agility",
	},
	WeaponTypeMaul: {
		Type:          WeaponTypeMaul,
		StrengthRange: primitives.AttributeRange{Min: 20, Max: 45},
		AgilityRange:  primitives.AttributeRange{Min: -15, Max: -5},
		Description:   "The heaviest weapon with devastating damage but significant agility penalty",
	},
	WeaponTypeMystery: {
		Type:          WeaponTypeMystery,
		StrengthRange: primitives.AttributeRange{Min: 1, Max: 50},
		AgilityRange:  primitives.AttributeRange{Min: -20, Max: 20},
		Description:   "A mysterious weapon with unpredictable effects",
	},
}

func GetWeaponConfig(t WeaponType) WeaponConfig {
	if cfg, exists := WeaponRegistry[t]; exists {
		return cfg
	}
	return WeaponRegistry[WeaponTypeMystery]
}

func (cfg WeaponConfig) GenerateAttributes(rnd *utils.RandomGenerator) primitives.Attributes {
	return primitives.Attributes{
		Strength: float64(utils.RandomRoundedFloatInRange(rnd, cfg.StrengthRange.Min, cfg.StrengthRange.Max)),
		Agility:  float64(utils.RandomRoundedFloatInRange(rnd, cfg.AgilityRange.Min, cfg.AgilityRange.Max)),
	}
}

func (cfg WeaponConfig) Validate() error {
	if cfg.StrengthRange.Min > cfg.StrengthRange.Max {
		return errors.New("invalid weapon config: min strength cannot be greater than max strength")
	}
	if cfg.AgilityRange.Min > cfg.AgilityRange.Max {
		return errors.New("invalid weapon config: min agility cannot be greater than max agility")
	}
	return nil
}
