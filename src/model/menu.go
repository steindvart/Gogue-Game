package model

import (
	"gogue/model/signal"
)

type MenuOption struct {
	Label  string
	Signal signal.Type
}

type Menu struct {
	options []MenuOption
	active  int
}

func NewMenu(options []MenuOption) *Menu {
	return &Menu{
		options: options,
		active:  0,
	}
}

func (m *Menu) Next() {
	if len(m.options) == 0 {
		return
	}
	m.active = (m.active + 1) % len(m.options)
}

func (m *Menu) Previous() {
	if len(m.options) == 0 {
		return
	}
	m.active = (m.active - 1 + len(m.options)) % len(m.options)
}

func (m *Menu) Select() MenuOption {
	if len(m.options) == 0 {
		return MenuOption{Label: "", Signal: -1}
	}
	return m.options[m.active]
}

func (m *Menu) GetActive() int {
	return m.active
}

func (m *Menu) GetOptionsLabels() []string {
	labels := make([]string, len(m.options))
	for i, opt := range m.options {
		labels[i] = opt.Label
	}
	return labels
}
