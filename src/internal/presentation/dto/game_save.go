package dto

import (
	"gogue/internal/model/primitives"
)

type GameSaveDto struct {
	SaveVersion  string         `json:"save_version"`
	LevelNumber  uint           `json:"level_number"`
	Rooms        []RoomDTO      `json:"rooms"`
	Passages     []PassageDTO   `json:"passages"`
	FinishPortal primitives.Box `json:"finish_portal"`
	Player       PlayerSaveDTO  `json:"player"`
	FogOfWar     FogOfWarDto    `json:"fog_of_war"`
}
