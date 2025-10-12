package state

import (
	"gogue/model/signal"

	"github.com/rivo/tview"
)

type State interface {
	Update() signal.Type
	Primitive() tview.Primitive
}
