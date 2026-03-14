package state

import (
	"gogue/internal/model/signals"
	"gogue/internal/presentation/action"
	"gogue/internal/presentation/dto"
	"gogue/internal/presentation/save"
	viewcli "gogue/internal/view/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scoreboard struct {
	view   *viewcli.Scoreboard
	signal signals.Type
}

func NewScoreboard(entries []dto.ScoreEntry) *Scoreboard {
	view := viewcli.NewScoreboard(entries)

	sb := &Scoreboard{
		view:   view,
		signal: signals.NoSignal,
	}

	sb.initInput()

	return sb
}

func GetScoreboard() (*Scoreboard, error) {
	scoreDto, err := save.LoadScore(ScoreFileName)
	if err != nil {
		return nil, err
	}

	return NewScoreboard(scoreDto.Entries), nil
}

func (s *Scoreboard) Update(dt float64) signals.Type {
	s.view.Update(dt)

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
