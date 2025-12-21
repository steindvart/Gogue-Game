package dto

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type ItemInfo struct {
	*EffectInfo
	Name          string
	Type          string
	CanUse        bool
	CanTake       bool
	TreasureValue int32
}

func ConvertPositionalItemToDto(item primitives.Positional2D[int]) *ItemInfo {
	if item == nil {
		return nil
	}

	info := &ItemInfo{}

	switch typedItem := item.(type) {
	case *items.Food:
		info.Name = typedItem.Name
		info.Type = "Food"
		info.EffectInfo = ConvertEffectToDto(typedItem.Effect)

	case *items.Elixir:
		info.Name = typedItem.Name
		info.Type = "Elixir"
		info.EffectInfo = ConvertEffectToDto(typedItem.Effect)

	case *items.Scroll:
		info.Name = typedItem.Name
		info.Type = "Scroll"
		info.EffectInfo = ConvertEffectToDto(typedItem.Effect)

	case *items.Weapon:
		info.Name = typedItem.Name
		info.Type = "Weapon"
		info.EffectInfo = ConvertEffectToDto(typedItem.Effect)

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
