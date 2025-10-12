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
		heroSpeed:  3, // символов в секунду
		enemyType:  []string{"🦇", "👻", "🧟‍♂️"},
		enemySpeed: []float64{2, 2, 2},
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
	b.moveHero(dt)
	b.moveEnemies(dt)
	b.Update()
}

func (b *RunnerBar) moveHero(dt float64) {
	b.heroPos += b.heroSpeed * dt
	if int(b.heroPos) >= b.length-2 {
		b.heroPos = 2
	}
}

func (b *RunnerBar) moveEnemies(dt float64) {
	for i := range b.enemyPos {
		b.enemyPos[i] -= b.enemySpeed[i] * dt
		if int(b.enemyPos[i]) < 0 {
			b.enemyPos[i] = float64(b.length - 4 - i*4)
		}
	}
}

func (b *RunnerBar) ResetPositions() {
	b.heroPos = 2
	b.enemyPos = make([]float64, len(b.enemyType))
	for i := range b.enemyType {
		b.enemyPos[i] = float64(b.length - 4 - i*4)
	}
}

func (b *RunnerBar) RenderBar() string {
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
	return string(bar)
}

func (b *RunnerBar) Primitive() tview.Primitive {
	return b.view
}
