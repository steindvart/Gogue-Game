package dto

// GameOverStats содержит статистику игровой сессии для отображения на экране завершения игры.
type GameOverStats struct {
	Treasures       int32
	EnemiesKilled   uint
	LevelReached    uint
	ConsumablesUsed uint
}
