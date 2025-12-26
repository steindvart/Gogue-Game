package dto

import (
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
)

type GameSaveDto struct {
	LevelNumber  uint            `json:"level_number"`
	Rooms        []world.Room    `json:"rooms"`
	Passages     []world.Passage `json:"passages"`
	FinishPortal primitives.Box  `json:"finish_portal"`
	Player       entities.Player `json:"player"`
	FogOfWar     FogOfWarDto     `json:"fog_of_war"`
}
