package model

type Signal int

const (
	NoSignal Signal = iota
	StopSignal
	NewGameSignal
	LoadGameSignal
	ShowScoreboardSignal
)
