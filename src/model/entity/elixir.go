package entity

import (
	"time"
)

type Elixir struct {
	Shape              Box
	EffectDuration     time.Duration
	AffectedAttributes Attributes
	Increment          uint
	Name               string
}
