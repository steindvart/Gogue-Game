package signal

type Type int

const (
	NoSignal Type = iota
	Stop
	NewGame
	LoadGame
	ShowScoreboard
)
