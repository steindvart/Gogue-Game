package dto

import "gogue/internal/model/primitives"

type FogOfWarDto struct {
	ExploredTiles []primitives.Point2D[int] `json:"explored_tiles"`
	Width         int                       `json:"width"`
	Height        int                       `json:"height"`
}
