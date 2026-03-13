package signals

type Type int

const (
	NoSignal Type = iota
	Stop
	NewGame
	LoadGame
	ShowScoreboard
	GameWon
	GameOver
	ReturnToMenu
)
