package model

type Signal int

type option struct {
	Label  string
	Signal Signal
}

type Menu struct {
	Options []option
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

func (m *Menu) Select() Signal {
	if len(m.Options) == 0 {
		return -1
	}
	return m.Options[m.Active].Signal
}
