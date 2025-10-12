package cli

import (
	"github.com/rivo/tview"
)

type RunnerBar struct {
	view    *tview.TextView
	length  int
	emojis  []string // список эмодзи (герой + враги)
	pos     float64  // позиция левого символа
	speed   float64  // скорость движения (символов в сек)
	spacing int      // отступ между эмодзи
}

// NewRunnerBar создаёт полоску с произвольным набором эмодзи, скоростью и отступом
func NewRunnerBar(length int, emojis []string, speed float64, spacing int) *RunnerBar {
	bar := &RunnerBar{
		view:    tview.NewTextView().SetDynamicColors(true),
		length:  length,
		emojis:  emojis,
		speed:   speed,
		spacing: spacing,
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
	b.moveAll(dt)
	b.Update()
}

func (b *RunnerBar) moveAll(dt float64) {
	b.pos += b.speed * dt
	totalLen := b.totalEmojisLen()
	if int(b.pos) > b.length {
		b.pos = -float64(totalLen)
	}
}

func (b *RunnerBar) ResetPositions() {
	b.pos = 0
}

func (b *RunnerBar) RenderBar() string {
	bar := make([]rune, b.length)
	for i := range bar {
		bar[i] = ' '
	}

	start := int(b.pos)
	for i, emoji := range b.emojis {
		emojiRunes := []rune(emoji)
		pos := start + i*(b.spacing+len(emojiRunes))
		for j, r := range emojiRunes {
			idx := (pos + j) % b.length
			if idx >= 0 && idx < b.length {
				bar[idx] = r
			}
		}
	}
	return string(bar)
}

func (b *RunnerBar) Primitive() tview.Primitive {
	return b.view
}

// totalEmojisLen возвращает суммарную длину всех эмодзи с учётом отступов
func (b *RunnerBar) totalEmojisLen() int {
	total := 0
	for i, emoji := range b.emojis {
		total += len([]rune(emoji))
		if i < len(b.emojis)-1 {
			total += b.spacing
		}
	}
	return total
}
