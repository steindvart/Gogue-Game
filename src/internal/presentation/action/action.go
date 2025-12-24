package action

type Type int

const (
	NoAction Type = iota
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	MoveLeftUpperCorner
	MoveRightUpperCorner
	MoveLefLowerCorner
	MoveRightLowerCorner
	Select
	Use
	Take
	Exit

	Num0
	Num1
	Num2
	Num3
	Num4
	Num5
	Num6
	Num7
	Num8
	Num9

	// Backpack actions
	ToggleBackpack
	OpenWeaponsTab
	OpenFoodTab
	OpenElixirsTab
	OpenScrollsTab
)
