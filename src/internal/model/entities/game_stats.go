package entities

type GameStats struct {
	EnemiesKilled uint `json:"enemies_killed"`
	FoodEaten     uint `json:"food_eaten"`
	ElixirsDrunk  uint `json:"elixirs_drunk"`
	ScrollsRead   uint `json:"scrolls_read"`
	HitsDealt     uint `json:"hits_dealt"`
	HitsMissed    uint `json:"hits_missed"`
	CellsMoved    uint `json:"cells_moved"`
}
