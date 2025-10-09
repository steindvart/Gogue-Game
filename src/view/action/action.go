package action

type Type int

const (
	NoAction Type = iota
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	Select
)
