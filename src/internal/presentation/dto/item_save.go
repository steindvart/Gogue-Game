package dto

import "gogue/internal/model/primitives"

type ItemSaveDTO struct {
	Type     string                  `json:"type"` // "Food", "Elixir" и т.д.
	Name     string                  `json:"name"`
	Position primitives.Point2D[int] `json:"position"`
	//Effect        *EffectSaveDTO          `json:"effect,omitempty"`
	TreasureValue int32 `json:"treasure_value,omitempty"`
}
