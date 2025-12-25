package dto

import "gogue/internal/model/primitives"

type PlayerSaveDTO struct {
	Position   primitives.Point2D[int] `json:"position"`
	Health     float64                 `json:"health"`
	MaxHealth  float64                 `json:"max_health"`
	Strength   float64                 `json:"strength"`
	Agility    float64                 `json:"agility"`
	Level      uint                    `json:"level"`
	Experience uint                    `json:"experience"`
	ViewRadius int                     `json:"view_radius"`

	EquippedWeapon *ItemSaveDTO `json:"equipped_weapon,omitempty"`

	// Backpack []ItemSaveDTO `json:"backpack"`
}
