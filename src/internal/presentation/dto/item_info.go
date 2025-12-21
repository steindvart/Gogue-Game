package dto

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type EffectInfo struct {
	MaxHealthModify float64
	HealthModify    float64
	StrengthModify  float64
	AgilityModify   float64
	DurationSteps   int
}

type ItemInfo struct {
	*EffectInfo
	Name          string
	Type          string
	CanUse        bool
	CanTake       bool
	TreasureValue int32
}

func GetItemInfo(item primitives.Positional2D[int]) *ItemInfo {
	if item == nil {
		return nil
	}

	info := &ItemInfo{}

	switch typedItem := item.(type) {
	case *items.Food:
		info.Name = typedItem.Name
		info.Type = "Food"
		info.EffectInfo = convertEffectToEffectInfo(typedItem.Effect)

	case *items.Elixir:
		info.Name = typedItem.Name
		info.Type = "Elixir"
		info.EffectInfo = convertEffectToEffectInfo(typedItem.Effect)

	case *items.Scroll:
		info.Name = typedItem.Name
		info.Type = "Scroll"
		info.EffectInfo = convertEffectToEffectInfo(typedItem.Effect)

	case *items.Weapon:
		info.Name = typedItem.Name
		info.Type = "Weapon"
		info.EffectInfo = convertEffectToEffectInfo(typedItem.Effect)

	case *items.Treasure:
		info.Name = typedItem.Name
		info.Type = "Treasure"
		info.EffectInfo = nil
	default:
		return nil
	}

	_, info.CanUse = item.(items.Usable)
	_, info.CanTake = item.(items.Takable)

	return info
}

func convertEffectToEffectInfo(effect *primitives.Effect) *EffectInfo {
	if effect == nil {
		return nil
	}

	return &EffectInfo{
		MaxHealthModify: effect.Attributes.MaxHealth,
		HealthModify:    effect.Attributes.Health,
		StrengthModify:  effect.Attributes.Strength,
		AgilityModify:   effect.Attributes.Agility,
		DurationSteps:   int(effect.Duration.Steps),
	}
}
