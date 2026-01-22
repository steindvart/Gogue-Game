package dto

import (
	"gogue/internal/model/entities"
)

type PlayerInfo struct {
	Treasures        int32
	Health           float64
	MaxHealth        float64
	Strength         float64
	Agility          float64
	TemporaryEffects []EffectInfo
	LevelNumber      uint
	*WeaponInfo
}

func ConvertPlayerToDto(player *entities.Player, levelNumber uint) *PlayerInfo {
	if player == nil {
		return nil
	}

	return &PlayerInfo{
		Treasures:        player.Backpack.Treasures,
		Health:           player.Attributes.Health,
		MaxHealth:        player.Attributes.MaxHealth,
		Strength:         player.Attributes.Strength,
		Agility:          player.Attributes.Agility,
		TemporaryEffects: ConvertEffectsToDto(player.TemporaryEffects),
		LevelNumber:      levelNumber,
		WeaponInfo:       ConvertWeaponToDto(player.Weapon),
	}
}
