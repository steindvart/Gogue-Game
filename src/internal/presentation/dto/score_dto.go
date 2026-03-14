package dto

import "gogue/internal/model/entities"

type ScoreEntry struct {
	Treasures    int32 `json:"treasures"`
	LevelReached uint  `json:"level_reached"`
	entities.GameStats
}

type ScoreDto struct {
	Entries []ScoreEntry `json:"entries"`
}
