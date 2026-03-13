package dto

import (
	"encoding/json"
	"fmt"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
)

type SerializableItems []primitives.Positional2D[int]

type SerializableEnemies []primitives.Positional2D[int]

type GameSaveDto struct {
	LevelNumber  uint                    `json:"level_number"`
	Size         primitives.Size2D[uint] `json:"size"`
	Rooms        []RoomSaveDto           `json:"rooms"`
	Passages     []world.Passage         `json:"passages"`
	FinishPortal primitives.Box          `json:"finish_portal"`
	FogOfWar     FogOfWarDto             `json:"fog_of_war"`
	Player       PlayerSaveDto           `json:"player"`
	Items        SerializableItems       `json:"items"`
	Enemies      SerializableEnemies     `json:"enemies"`
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
	Position         primitives.Point2D[int] `json:"position"`
	MaxHealth        float64                 `json:"max_health"`
	Health           float64                 `json:"health"`
	Agility          float64                 `json:"agility"`
	Strength         float64                 `json:"strength"`
	TemporaryEffects []*primitives.Effect    `json:"temporary_effects"`
	ViewRadius       int                     `json:"view_radius"`
	Weapon           *items.Weapon           `json:"weapon"`
	Backpack         BackpackSaveDto         `json:"backpack"`
	EnemiesKilled    uint                    `json:"enemies_killed"`
	ConsumablesUsed  uint                    `json:"consumables_used"`
}

type EffectSaveDto struct {
	Duration   EffectDurationSaveDto `json:"duration"`
	Attributes AttributesSaveDto     `json:"attributes"`
}

type ItemSaveDto struct {
	Name string         `json:"name"`
	Box  primitives.Box `json:"box"`
}

type WeaponSaveDto struct {
	ItemSaveDto
	Type items.WeaponType `json:"type"`
}
type EffectDurationSaveDto struct {
	Type  primitives.EffectDurationType `json:"type"` // ← предполагается, что это int/string
	Turns uint32                        `json:"turns"`
}

type AttributesSaveDto struct {
	MaxHealth float64 `json:"max_health"`
	Health    float64 `json:"health"`
	Agility   float64 `json:"agility"`
	Strength  float64 `json:"strength"`
}

type BackpackSaveDto struct {
	Capacity  uint            `json:"capacity"`
	ItemsNum  uint            `json:"items_num"`
	Elixirs   []*items.Elixir `json:"elixirs"`
	Scrolls   []*items.Scroll `json:"scrolls"`
	Foods     []*items.Food   `json:"foods"`
	Weapons   []*items.Weapon `json:"weapons"`
	Treasures int32           `json:"treasures"`
}

func (si SerializableItems) MarshalJSON() ([]byte, error) {
	result := make([]map[string]interface{}, len(si))

	for i, item := range si {
		var typ string
		var data interface{}

		switch v := item.(type) {
		case *items.Food:
			typ = "food"
			data = v
		case *items.Elixir:
			typ = "elixir"
			data = v
		case *items.Scroll:
			typ = "scroll"
			data = v
		case *items.Weapon:
			typ = "weapon"
			data = v
		default:
			return nil, fmt.Errorf("unsupported item type at index %d %s", i, v)
		}

		result[i] = map[string]interface{}{
			"type": typ,
			"data": data,
		}
	}

	return json.Marshal(result)
}

func (si *SerializableItems) UnmarshalJSON(data []byte) error {
	var rawItems []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return err
	}

	resultItems := make([]primitives.Positional2D[int], len(rawItems))

	for i, rawItem := range rawItems {
		var itemType string
		if err := json.Unmarshal(rawItem["type"], &itemType); err != nil {
			return fmt.Errorf("missing or invalid 'type' at index %d", i)
		}

		var item primitives.Positional2D[int]

		switch itemType {
		case "food":
			var food items.Food
			if err := json.Unmarshal(rawItem["data"], &food); err != nil {
				return fmt.Errorf("failed to unmarshal food at index %d: %w", i, err)
			}
			item = &food

		case "elixir":
			var elixir items.Elixir
			if err := json.Unmarshal(rawItem["data"], &elixir); err != nil {
				return fmt.Errorf("failed to unmarshal elixir at index %d: %w", i, err)
			}
			item = &elixir

		case "scroll":
			var scroll items.Scroll
			if err := json.Unmarshal(rawItem["data"], &scroll); err != nil {
				return fmt.Errorf("failed to unmarshal scroll at index %d: %w", i, err)
			}
			item = &scroll

		case "weapon":
			var weapon items.Weapon
			if err := json.Unmarshal(rawItem["data"], &weapon); err != nil {
				return fmt.Errorf("failed to unmarshal weapon at index %d: %w", i, err)
			}
			item = &weapon

		default:
			return fmt.Errorf("unknown item type %q at index %d", itemType, i)
		}

		resultItems[i] = item
	}

	*si = resultItems
	return nil
}

func (si SerializableEnemies) MarshalJSON() ([]byte, error) {
	result := make([]map[string]interface{}, len(si))

	for i, item := range si {
		var typ string
		var data interface{}

		switch v := item.(type) {
		case *entities.Zombie:
			typ = "zombie"
			data = v
		case *entities.Vampire:
			typ = "vampire"
			data = v
		case *entities.Ghost:
			typ = "ghost"
			data = v
		case *entities.Ogre:
			typ = "ogre"
			data = v
		case *entities.SnakeMage:
			typ = "snakeMage"
			data = v
		case *entities.Mimic:
			typ = "mimic"
			data = v
		default:
			return nil, fmt.Errorf("unsupported enemy type at index %d %s", i, v)
		}

		result[i] = map[string]interface{}{
			"type": typ,
			"data": data,
		}
	}

	return json.Marshal(result)
}

func (si *SerializableEnemies) UnmarshalJSON(data []byte) error {
	var rawItems []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return err
	}

	resultItems := make([]primitives.Positional2D[int], len(rawItems))

	for i, rawItem := range rawItems {
		var itemType string
		if err := json.Unmarshal(rawItem["type"], &itemType); err != nil {
			return fmt.Errorf("missing or invalid 'type' at index %d", i)
		}

		var item primitives.Positional2D[int]

		switch itemType {
		case "zombie":
			var zombie entities.Zombie
			if err := json.Unmarshal(rawItem["data"], &zombie); err != nil {
				return fmt.Errorf("failed to unmarshal zombie at index %d: %w", i, err)
			}
			item = &zombie

		case "vampire":
			var vampire entities.Vampire
			if err := json.Unmarshal(rawItem["data"], &vampire); err != nil {
				return fmt.Errorf("failed to unmarshal vampire at index %d: %w", i, err)
			}
			item = &vampire

		case "ghost":
			var ghost entities.Ghost
			if err := json.Unmarshal(rawItem["data"], &ghost); err != nil {
				return fmt.Errorf("failed to unmarshal ghost at index %d: %w", i, err)
			}
			item = &ghost

		case "ogre":
			var ogre entities.Ogre
			if err := json.Unmarshal(rawItem["data"], &ogre); err != nil {
				return fmt.Errorf("failed to unmarshal ogre at index %d: %w", i, err)
			}
			item = &ogre

		case "snakeMage":
			var snakeMage entities.SnakeMage
			if err := json.Unmarshal(rawItem["data"], &snakeMage); err != nil {
				return fmt.Errorf("failed to unmarshal snakeMage at index %d: %w", i, err)
			}
			item = &snakeMage

		case "mimic":
			var mimic entities.Mimic
			if err := json.Unmarshal(rawItem["data"], &mimic); err != nil {
				return fmt.Errorf("failed to unmarshal mimic at index %d: %w", i, err)
			}
			item = &mimic

		default:
			return fmt.Errorf("unknown enemy type %q at index %d", itemType, i)
		}

		resultItems[i] = item
	}

	*si = resultItems
	return nil
}
