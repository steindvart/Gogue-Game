package dto

type ScoreEntry struct {
	Treasures     int32 `json:"treasures"`
	LevelReached  uint  `json:"level_reached"`
	EnemiesKilled uint  `json:"enemies_killed"`
	FoodEaten     uint  `json:"food_eaten"`
	ElixirsDrunk  uint  `json:"elixirs_drunk"`
	ScrollsRead   uint  `json:"scrolls_read"`
	HitsDealt     uint  `json:"hits_dealt"`
	HitsMissed    uint  `json:"hits_missed"`
	CellsMoved    uint  `json:"cells_moved"`
}

type ScoreDto struct {
	Entries []ScoreEntry `json:"entries"`
}
