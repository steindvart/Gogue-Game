package state

import (
	"fmt"
	"gogue/internal/model/signals"
	"gogue/internal/presentation/action"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scoreboard struct {
	view   *viewcli.Scoreboard
	signal signals.Type
}

func NewScoreboard() *Scoreboard {
	var scores []string
	const width = 20

	for i := 1; i < 100; i++ {
		scores = append(scores, fmt.Sprintf("%*d", width, i))
	}

	view := viewcli.NewScoreboard(scores)

	sb := &Scoreboard{
		view:   view,
		signal: signals.NoSignal,
	}

	sb.initInput()

	return sb
}

func GetScoreboard() *Scoreboard {
	return NewScoreboard()
}

func (s *Scoreboard) Update(float64) signals.Type {
	sig := s.signal
	s.signal = signals.NoSignal // сброс сигнала после чтения
	return sig
}

func (s *Scoreboard) Primitive() tview.Primitive {
	return s.view.Primitive()
}

func (s *Scoreboard) initInput() {
	s.view.SetInputCapture(s.handleEvent)
}

func (s *Scoreboard) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	switch s.eventToAction(event) {
	case action.Exit:
		s.signal = signals.Stop

		return nil
	default:
		return event
	}
}

func (s *Scoreboard) eventToAction(event *tcell.EventKey) action.Type {
	switch event.Key() {
	case tcell.KeyEsc:
		return action.Exit
	}
	return action.NoAction
}
