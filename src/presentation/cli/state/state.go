package state

import (
	"gogue/model/signal"
	"gogue/view/action"
)

type State interface {
	Input(a action.Type) signal.Type
	Update() signal.Type
	Render()
}

type Renderer interface {
	Render()
}
