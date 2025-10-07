package model

type Signal int

type MenuOption struct {
	Label  string
	Signal Signal
}

type Menu struct {
	Options []MenuOption
	Active  int
}

func (m *Menu) Next() {
	if len(m.Options) == 0 {
		return
	}
	m.Active = (m.Active + 1) % len(m.Options)
}

func (m *Menu) Previous() {
	if len(m.Options) == 0 {
		return
	}
	m.Active = (m.Active - 1 + len(m.Options)) % len(m.Options)
}

func (m *Menu) Select() MenuOption {
	if len(m.Options) == 0 {
		return MenuOption{Label: "", Signal: -1}
	}
	return m.Options[m.Active]
}
