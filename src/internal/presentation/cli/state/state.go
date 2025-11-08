package state

import (
	"gogue/internal/model/signal"

	"github.com/rivo/tview"
)

type State interface {
	Update(dt float64) signal.Type
	Primitive() tview.Primitive
}
