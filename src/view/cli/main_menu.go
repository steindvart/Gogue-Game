package cli

import (
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

type MainMenuTView struct {
	flex  *tview.Flex
	list  *tview.List
	title *tview.TextView
	hall  *tview.TextView
}

func NewMainMenuTView(options []string, active int) *MainMenuTView {
	title := tview.NewTextView().SetDynamicColors(true)
	for _, line := range titleArt {
		title.Write([]byte(line + "\n"))
	}
	title.SetTextAlign(tview.AlignCenter)
	title.SetBorder(false)

	hall := tview.NewTextView().SetDynamicColors(true)
	for _, line := range hallArt {
		hall.Write([]byte(line + "\n"))
	}
	hall.SetTextAlign(tview.AlignCenter)
	hall.SetBorder(false)

	list := tview.NewList()
	for _, opt := range options {
		list.AddItem(opt, "", 0, nil)
	}
	list.SetCurrentItem(active)
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorYellow)
	list.SetBorder(false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, len(titleArt)+1, 0, false).
		AddItem(list, 0, 1, true).
		AddItem(hall, len(hallArt)+1, 0, false)

	return &MainMenuTView{
		flex:  flex,
		list:  list,
		title: title,
		hall:  hall,
	}
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
