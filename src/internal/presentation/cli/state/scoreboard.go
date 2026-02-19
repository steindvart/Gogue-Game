package state

import (
	"fmt"
	"gogue/internal/model/signals"
	"gogue/internal/presentation/action"
	"gogue/internal/presentation/save"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ScoreWidth = 20 // Отступ чтобы строка начиналась с пробелов, а score был в конце (как в exel)

type Scoreboard struct {
	view   *viewcli.Scoreboard
	signal signals.Type
}

func NewScoreboard(scores []int32) *Scoreboard {
	strScores := make([]string, len(scores))

	for i, score := range scores {
		strScores[i] = fmt.Sprintf("%*d", ScoreWidth, score)
	}

	view := viewcli.NewScoreboard(strScores)

	sb := &Scoreboard{
		view:   view,
		signal: signals.NoSignal,
	}

	sb.initInput()

	return sb
}

func GetScoreboard() (*Scoreboard, error) {
	dto, err := save.LoadScore(ScoreFileName)
	if err != nil {
		return nil, err
	}

	return NewScoreboard(dto.Scores), nil
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
