package dto

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
)

type ItemInfo struct {
	Name           string
	Type           string
	HealthModify   float64
	StrengthModify float64
	AgilityModify  float64
	DurationSteps  int
	CanUse         bool
	CanTake        bool
	TreasureValue  int32
}

func GetItemInfo(item primitives.Positional2D[int]) *ItemInfo {
	if item == nil {
		return nil
	}

	info := &ItemInfo{}

	if baseItem, ok := item.(items.Item); ok {
		info.Name = baseItem.Name
	}

	// Определяем тип предмета и заполняем соответствующую информацию
	switch typedItem := item.(type) {
	case *items.Food:
		info.Name = typedItem.Name
		info.Type = "Food"
		if typedItem.Effect != nil {
			info.HealthModify = typedItem.Effect.Attributes.Health
			info.StrengthModify = typedItem.Effect.Attributes.Strength
			info.AgilityModify = typedItem.Effect.Attributes.Agility
			info.DurationSteps = int(typedItem.Effect.Duration.Steps)
		}
		info.CanUse = true
		info.CanTake = false

	case *items.Elixir:
		info.Name = typedItem.Name
		info.Type = "Elixir"
		if typedItem.Effect != nil {
			info.HealthModify = typedItem.Effect.Attributes.Health
			info.StrengthModify = typedItem.Effect.Attributes.Strength
			info.AgilityModify = typedItem.Effect.Attributes.Agility
			info.DurationSteps = int(typedItem.Effect.Duration.Steps)
		}
		info.CanUse = true
		info.CanTake = false

	case *items.Scroll:
		info.Name = typedItem.Name
		info.Type = "Scroll"
		if typedItem.Effect != nil {
			info.HealthModify = typedItem.Effect.Attributes.Health
			info.StrengthModify = typedItem.Effect.Attributes.Strength
			info.AgilityModify = typedItem.Effect.Attributes.Agility
			info.DurationSteps = 0 // Свитки действуют мгновенно
		}
		info.CanUse = true
		info.CanTake = false

	case *items.Weapon:
		info.Name = typedItem.Name
		info.Type = "Weapon"
		if typedItem.Effect != nil {
			info.HealthModify = typedItem.Effect.Attributes.Health
			info.StrengthModify = typedItem.Effect.Attributes.Strength
			info.AgilityModify = typedItem.Effect.Attributes.Agility
			info.DurationSteps = 0 // Оружие действует постоянно после подбора
		}
		info.CanUse = true
		info.CanTake = true

	case *items.Treasure:
		info.Name = typedItem.Name
		info.Type = "Treasure"
		info.TreasureValue = typedItem.Value
		info.CanUse = false
		info.CanTake = true

	default:
		// Неизвестный тип предмета
		return nil
	}

	return info
}
