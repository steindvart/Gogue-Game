package entity

import (
	"time"
)

type Elixir struct {
	Geometry           Object
	EffectDuration     time.Duration
	AffectedAttributes Attributes[uint]
	Increment          uint
	Name               string
}
