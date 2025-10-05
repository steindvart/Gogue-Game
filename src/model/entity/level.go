package entity

type Level struct {
	Rooms       []Room
	Passages    []Passage
	LevelNumber uint
	LevelEnd    Box
}
