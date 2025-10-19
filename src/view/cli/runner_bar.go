package cli

import (
	"github.com/rivo/tview"
)

type RunnerBar struct {
	view    *tview.TextView
	length  int
	content string  // строка, которую нужно прокручивать
	pos     float64 // позиция первого символа (сдвиг)
	speed   float64 // скорость движения (символов в сек)
}

// NewRunnerBar принимает строку (например, эмодзи через пробел), которую будет циклически прокручивать
func NewRunnerBar(length int, content string, speed float64) *RunnerBar {
	bar := &RunnerBar{
		view:    tview.NewTextView().SetDynamicColors(true),
		length:  length,
		content: content,
		speed:   speed,
	}
	bar.ResetPositions()
	bar.view.SetTextAlign(tview.AlignLeft)
	bar.view.SetBorder(false)
	bar.Update()
	return bar
}

func (b *RunnerBar) Update() {
	b.view.SetText(b.RenderBar())
}

func (b *RunnerBar) UpdatePositions(dt float64) {
	b.pos += b.speed * dt
	contentLen := len([]rune(b.content))
	if contentLen == 0 {
		b.pos = 0
	} else if int(b.pos) >= contentLen {
		b.pos = 0
	}
	b.Update()
}

func (b *RunnerBar) ResetPositions() {
	b.pos = 0
}

func (b *RunnerBar) RenderBar() string {
	bar := make([]rune, b.length)
	for i := range bar {
		bar[i] = ' '
	}
	contentRunes := []rune(b.content)
	contentLen := len(contentRunes)
	if contentLen == 0 {
		return string(bar)
	}
	start := int(b.pos) % contentLen
	for i := 0; i < b.length; i++ {
		idx := (start + i) % contentLen
		bar[i] = contentRunes[idx]
	}
	return string(bar)
}

func (b *RunnerBar) Primitive() tview.Primitive {
	return b.view
}
