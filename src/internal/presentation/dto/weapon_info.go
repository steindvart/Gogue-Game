package dto

import "gogue/internal/model/items"

type WeaponInfo struct {
	*EffectInfo
	Type string
}

func ConvertWeaponToDto(weapon *items.Weapon) *WeaponInfo {
	if weapon == nil {
		return nil
	}

	return &WeaponInfo{
		EffectInfo: ConvertEffectToDto(weapon.Effect),
		Type:       string(weapon.Type),
	}
}
