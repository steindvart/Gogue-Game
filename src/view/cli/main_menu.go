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

type MainMenu struct {
	flex  *tview.Flex
	list  *tview.List
	title *tview.TextView
	hall  *tview.TextView
	frame int // для анимации радуги
}

func NewMainMenu(options []string, active int) *MainMenu {
	title := newTitle()
	list := newList(options, active)

	// Создаём пустой Box с тёмным фоном для выравнивания элементов
	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	// Центрируем список по горизонтали с помощью Flex
	hFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(gapBox, 0, 1, false). // пустое пространство слева
		AddItem(list, 10, 0, true).   // ширина списка (можно скорректировать)
		AddItem(gapBox, 0, 1, false)  // пустое пространство справа

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gapBox, 10, 0, false).
		AddItem(title, len(titleArt)+3, 0, false).
		AddItem(hFlex, 0, 1, true)

	m := &MainMenu{
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

func newList(options []string, active int) *tview.List {
	list := tview.NewList()
	for _, opt := range options {
		list.AddItem(opt, "", 0, nil)
	}

	list.SetCurrentItem(active)
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorWhite)
	list.ShowSecondaryText(false)
	list.SetBorder(false)

	return list
}

func (m *MainMenu) SetRainbowTitleFrame(frame int) {
	m.frame = frame
	m.title.SetText(DrawRainbowTitle(frame))
}

func (m *MainMenu) Primitive() tview.Primitive {
	return m.flex
}

func (m *MainMenu) SetActive(idx int) {
	m.list.SetCurrentItem(idx)
}

func (m *MainMenu) SetOptions(options []string) {
	m.list.Clear()
	for _, opt := range options {
		m.list.AddItem(opt, "", 0, nil)
	}
}

func (m *MainMenu) SetInputCapture(handler func(event *tcell.EventKey) *tcell.EventKey) {
	m.list.SetInputCapture(handler)
}

// hsvToRGB возвращает цвет tcell.Color по HSV (hue [0..360), s,v [0..1])
func hsvToRGB(h, s, v float64) tcell.Color {
	var r, g, b float64
	i := int(h/60.0) % 6
	f := h/60.0 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch i {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return tcell.NewRGBColor(int32(r*255), int32(g*255), int32(b*255))
}

func DrawRainbowTitle(frame int) string {
	const paletteSize = 64 // Чем больше, тем плавнее

	var palette [paletteSize]tcell.Color
	for i := 0; i < paletteSize; i++ {
		h := float64(i) * 360.0 / float64(paletteSize)
		palette[i] = hsvToRGB(h, 1.0, 1.0)
	}

	var sb strings.Builder
	for row, line := range titleArt {
		for col, ch := range line {
			colorIdx := (col + row + frame) % paletteSize
			fmt.Fprintf(&sb, "[#%06x]%c", palette[colorIdx].Hex(), ch)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
