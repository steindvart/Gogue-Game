package state

import (
	"gogue/internal/model/signals"

	"github.com/rivo/tview"
)

type State interface {
	Update(dt float64) signals.Type
	Primitive() tview.Primitive
}
