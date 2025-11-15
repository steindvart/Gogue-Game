package cli

import (
	"gogue/internal/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	Box *tview.Box
}

func NewGameBox() *Game {
	return &Game {
		Box: tview.NewBox().SetBorder(true).SetTitle("Game"),
	}
}

func (g *Game) SetFieldToScreen(screen tcell.Screen, field [][]utils.EntityType, ox, oy int) {
	// field := g.makeField(w, h)
	for y := range len(field) {
		for x := range len(field[0]) {
			ch := ' '
			switch field[y][x] {
			case utils.EntityTypePlayer:
				ch = '🦸'
			case utils.EntityTypeHorizontalWall:
				ch = '—'
			case utils.EntityTypeVerticalWall:
				ch = '|'
			case utils.EntityTypePortal:
				ch = 'O'
			case utils.EntityTypePassage:
				ch = '*'
			case utils.EntityTypeDoorOne:
				ch = '['
			case utils.EntityTypeDoorTwo:
				ch = ']'
			case utils.EntityTypeZombie:
				ch = '🧟'
			case utils.EntityTypeVampire:
				ch = '🧛'
			case utils.EntityTypeGhost:
				ch = '👻'
			case utils.EntityTypeOgre:
				ch = '👹'
			case utils.EntityTypeSnakeMage:
				ch = '🐍'
			}
			screen.SetContent(ox+x, oy+y, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
}