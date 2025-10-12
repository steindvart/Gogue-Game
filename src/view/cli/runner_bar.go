package cli

import (
	"github.com/rivo/tview"
)

type RunnerBar struct {
	view       *tview.TextView
	length     int
	heroPos    float64
	heroSpeed  float64
	enemyPos   []float64
	enemyType  []string
	enemySpeed []float64
}

func NewRunnerBar(length int) *RunnerBar {
	bar := &RunnerBar{
		view:       tview.NewTextView().SetDynamicColors(true),
		length:     length,
		heroPos:    2,
		heroSpeed:  3, // символов в секунду
		enemyPos:   []float64{float64(length - 4), float64(length - 8), float64(length - 12)},
		enemyType:  []string{"🦇", "👻", "🧟‍♂️"},
		enemySpeed: []float64{2, 2, 2}, // скорость врагов (символов в сек)
	}
	bar.view.SetTextAlign(tview.AlignLeft)
	bar.view.SetBorder(false)
	bar.Update()
	return bar
}

func (b *RunnerBar) Update() {
	bar := make([]rune, b.length)
	for i := range bar {
		bar[i] = ' '
	}
	// Враги
	for i, pos := range b.enemyPos {
		ipos := int(pos)
		if ipos >= 0 && ipos < b.length {
			et := []rune(b.enemyType[i])
			for j, r := range et {
				if ipos+j < b.length {
					bar[ipos+j] = r
				}
			}
		}
	}
	// Персонаж
	hero := []rune("🦸")
	hpos := int(b.heroPos)
	for j, r := range hero {
		if hpos+j < b.length {
			bar[hpos+j] = r
		}
	}
	b.view.SetText(string(bar))
}

// UpdatePositions обновляет позиции объектов с учётом dt (секунды)
func (b *RunnerBar) UpdatePositions(dt float64) {
	b.heroPos += b.heroSpeed * dt
	if int(b.heroPos) >= b.length-2 {
		b.heroPos = 2
	}

	for i := range b.enemyPos {
		b.enemyPos[i] -= b.enemySpeed[i] * dt
		if int(b.enemyPos[i]) < 0 {
			b.enemyPos[i] = float64(b.length - 4 - i*4)
		}
	}
	b.Update()
}

func (b *RunnerBar) Primitive() tview.Primitive {
	return b.view
}
