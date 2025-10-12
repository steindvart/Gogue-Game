package cli

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var titleArt = []string{
	"   _____                           _____                      ",
	"  / ____|                         / ____|                     ",
	" | |  __  ___   __ _ _   _  ___  | |  __  __ _ _ __ ___   ___ ",
	" | | |_ |/ _ \\ / _` | | | |/ _ \\ | | |_ |/ _` | '_ ` _ \\ / _ \\",
	" | |__| | (_) | (_| | |_| |  __/ | |__| | (_| | | | | | |  __/",
	"  \\_____|\\___/ \\__, |\\__,_|\\___|  \\_____|\\__,_|_| |_| |_|\\___|",
	"                __/ |                                         ",
	"               |___/                                          ",
}

type MainMenuTView struct {
	flex  *tview.Flex
	list  *tview.List
	title *tview.TextView
	hall  *tview.TextView
	frame int // для анимации радуги
}

func NewMainMenuTView(options []string, active int) *MainMenuTView {
	title := newTitle()
	list := newList(options, active)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, len(titleArt)+1, 0, false).
		AddItem(list, 0, 1, true)

	m := &MainMenuTView{
		flex:  flex,
		list:  list,
		title: title,
		hall:  nil,
		frame: 0,
	}

	m.SetRainbowTitleFrame(0)
	return m
}

func newTitle() *tview.TextView {
	// Изначально пусто, будет обновляться через SetRainbowTitleFrame
	title := tview.NewTextView().SetDynamicColors(true)
	title.SetTextAlign(tview.AlignCenter)
	title.SetBorder(false)

	return title
}

func newHall() *tview.TextView {
	var hallArt = []string{
		" _____________________________________________",
		"|.'',                                     ,''.|",
		"|.'.'',                                 ,''.'.|",
		"|.'.'.'',                             ,''.'.'.|",
		"|.'.'.'.'',                         ,''.'.'.'.|",
		"|.'.'.'.'.|                         |.'.'.'.'.|",
		"|.'.'.'.'.|===;                 ;===|.'.'.'.'.|",
		"|.'.'.'.'.|:::|',             ,'|:::|.'.'.'.'.|",
		"|.'.'.'.'.|---|'.|, _______ ,|.'|---|.'.'.'.'.|",
		"|.'.'.'.'.|:::|'.|'|???????|'|.'|:::|.'.'.'.'.|",
		"|,',',',',|---|',|'|???????|'|,'|---|,',',',',|",
		"|.'.'.'.'.|:::|'.|'|???????|'|.'|:::|.'.'.'.'.|",
		"|.'.'.'.'.|---|','   /%%%\\   ','|---|.'.'.'.'.|",
		"|.'.'.'.'.|===:'    /%%%%%\\    ':===|.'.'.'.'.|",
		"|.'.'.'.'.|%%%%%%%%%%%%%%%%%%%%%%%%%|.'.'.'.'.|",
		"|.'.'.'.','       /%%%%%%%%%\\       ','.'.'.'.|",
		"|.'.'.','        /%%%%%%%%%%%\\        ','.'.'.|",
		"|.'.','         /%%%%%%%%%%%%%\\         ','.'.|",
		"|.','          /%%%%%%%%%%%%%%%\\          ','.|",
		"|;____________/%%%%%%%%%%%%%%%%%\\____________;|",
	}

	hall := tview.NewTextView().SetDynamicColors(true)
	for _, line := range hallArt {
		hall.Write([]byte(line + "\n"))
	}
	hall.SetTextAlign(tview.AlignCenter)
	hall.SetBorder(false)
	return hall
}

func newList(options []string, active int) *tview.List {
	list := tview.NewList()
	for _, opt := range options {
		list.AddItem(opt, "", 0, nil)
	}

	list.SetCurrentItem(active)
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorYellow)
	list.SetBorder(false)

	return list
}

func (m *MainMenuTView) SetRainbowTitleFrame(frame int) {
	m.frame = frame
	m.title.SetText(DrawRainbowTitle(frame))
}

func (m *MainMenuTView) Primitive() tview.Primitive {
	return m.flex
}

func (m *MainMenuTView) SetActive(idx int) {
	m.list.SetCurrentItem(idx)
}

func (m *MainMenuTView) SetOptions(options []string) {
	m.list.Clear()
	for _, opt := range options {
		m.list.AddItem(opt, "", 0, nil)
	}
}

func (m *MainMenuTView) SetInputCapture(handler func(event *tcell.EventKey) *tcell.EventKey) {
	m.list.SetInputCapture(handler)
}

// DrawRainbowTitle формирует titleArt с анимированной радугой
func DrawRainbowTitle(frame int) string {
	var rainbowColors = []tcell.Color{
		tcell.ColorRed,
		tcell.ColorOrange,
		tcell.ColorYellow,
		tcell.ColorGreen,
		tcell.ColorBlue,
		tcell.ColorIndigo,
		tcell.ColorViolet,
	}

	var sb strings.Builder
	for row, line := range titleArt {
		for col, ch := range line {
			colorIdx := (col + row + frame) % len(rainbowColors)
			fmt.Fprintf(&sb, "[#%06x]%c", rainbowColors[colorIdx].Hex(), ch)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
