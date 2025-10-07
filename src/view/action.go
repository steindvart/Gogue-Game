package view

type ActionType int

const (
	NoAction ActionType = iota
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	Select
)
