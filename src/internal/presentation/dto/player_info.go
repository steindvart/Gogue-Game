package dto

import "gogue/internal/model/entities"

type PlayerInfo struct {
	Health           float64
	MaxHealth        float64
	Strength         float64
	Agility          float64
	TemporaryEffects []EffectInfo
	*WeaponInfo
}

func ConvertPlayerToDto(player *entities.Player) *PlayerInfo {
	if player == nil {
		return nil
	}

	return &PlayerInfo{
		Health:           player.Attributes.Health,
		MaxHealth:        player.Attributes.MaxHealth,
		Strength:         player.Attributes.Strength,
		Agility:          player.Attributes.Agility,
		TemporaryEffects: ConvertEffectsToDto(player.TemporaryEffects),
		WeaponInfo:       ConvertWeaponToDto(player.Weapon),
	}
}
