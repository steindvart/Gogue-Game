package dto

import (
	"gogue/internal/model/primitives"
)

type EffectInfo struct {
	MaxHealthModify float64
	HealthModify    float64
	StrengthModify  float64
	AgilityModify   float64
	DurationSteps   int
}

func ConvertEffectToDto(effect *primitives.Effect) *EffectInfo {
	if effect == nil {
		return nil
	}

	return &EffectInfo{
		DurationSteps:  int(effect.Duration.Steps),
		HealthModify:   effect.Attributes.Health,
		StrengthModify: effect.Attributes.Strength,
		AgilityModify:  effect.Attributes.Agility,
	}
}

func ConvertEffectsToDto(effects []*primitives.Effect) []EffectInfo {
	dtoEffects := make([]EffectInfo, 0, len(effects))
	for _, effect := range effects {
		dtoEffects = append(dtoEffects, *ConvertEffectToDto(effect))
	}
	return dtoEffects
}
