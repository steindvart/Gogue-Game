package dto

import "gogue/internal/model/primitives"

type RoomDTO struct {
	Box   primitives.Box            `json:"box"`
	Doors []primitives.Point2D[int] `json:"doors"`
}

type PassageDTO struct {
	DoorOne primitives.Point2D[int]   `json:"door_one"`
	DoorTwo primitives.Point2D[int]   `json:"door_two"`
	Way     []primitives.Point2D[int] `json:"way"`
}
