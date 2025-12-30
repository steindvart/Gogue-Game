package dto

import (
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
)

type GameSaveDto struct {
	LevelNumber  uint                    `json:"level_number"`
	Size         primitives.Size2D[uint] `json:"size"`
	Rooms        []RoomSaveDto           `json:"rooms"`
	Passages     []world.Passage         `json:"passages"`
	FinishPortal primitives.Box          `json:"finish_portal"`
	FogOfWar     FogOfWarDto             `json:"fog_of_war"`
	Player       PlayerSaveDto           `json:"player"`
}

type RoomSaveDto struct {
	Box   primitives.Box            `json:"box"`
	Doors []primitives.Point2D[int] `json:"doors"`
}

type FogOfWarDto struct {
	ExploredTiles []primitives.Point2D[int] `json:"explored_tiles"`
	Width         int                       `json:"width"`
	Height        int                       `json:"height"`
}

type PlayerSaveDto struct {
	Position   primitives.Point2D[int] `json:"position"`
	MaxHealth  float64                 `json:"max_health"`
	Health     float64                 `json:"health"`
	Agility    float64                 `json:"agility"`
	Strength   float64                 `json:"strength"`
	Experience uint                    `json:"experience"`
	Level      uint                    `json:"level"`
	ViewRadius int                     `json:"view_radius"`
	// Временно убираем Backpack, Weapon — добавим позже
}

//type PlayerSaveDto struct {
//	// Character данные
//	Position         primitives.Point2D[int] `json:"position"` // из Box
//	MaxHealth        float64                 `json:"max_health"`
//	Health           float64                 `json:"health"`
//	Agility          float64                 `json:"agility"`
//	Strength         float64                 `json:"strength"`
//	TemporaryEffects []EffectSaveDto         `json:"temporary_effects"`
//
//	// Player данные
//	Experience uint `json:"experience"`
//	Level      uint `json:"level"`
//	ViewRadius int  `json:"view_radius"`
//
//	// Backpack данные
//	BackpackCapacity uint            `json:"backpack_capacity"`
//	Elixirs          []ItemSaveDto   `json:"elixirs"`
//	Scrolls          []ItemSaveDto   `json:"scrolls"`
//	Foods            []ItemSaveDto   `json:"foods"`
//	Weapons          []WeaponSaveDto `json:"weapons"`
//	Treasures        int32           `json:"treasures"`
//
//	// Weapon данные
//	EquippedWeapon *WeaponSaveDto `json:"equipped_weapon,omitempty"`
//}
//
//type EffectSaveDto struct {
//	Duration   primitives.EffectDuration `json:"duration"`
//	Attributes AttributesSaveDto         `json:"attributes"`
//}
//
//type AttributesSaveDto struct {
//	MaxHealth float64 `json:"max_health"`
//	Health    float64 `json:"health"`
//	Agility   float64 `json:"agility"`
//	Strength  float64 `json:"strength"`
//}
//
//type ItemSaveDto struct {
//	Name string         `json:"name"`
//	Box  primitives.Box `json:"box"` // предполагается, что Box — публичная структура
//}
//
//type WeaponSaveDto struct {
//	ItemSaveDto                  // встраиваем данные Item
//	Effect      EffectSaveDto    `json:"effect"`
//	Type        items.WeaponType `json:"type"`
//}
