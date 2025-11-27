package cli

import (
	"gogue/internal/common"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	*tview.Box
}

func NewGame() *Game {
	return &Game{
		Box: tview.NewBox().SetBorder(true).SetTitle("Game"),
	}
}

func (*Game) SetFieldToScreen(screen tcell.Screen, field [][]common.GameEntityType, ox, oy int) {
	for y := range len(field) {
		for x := range len(field[0]) {
			ch := ' '
			switch field[y][x] {
			case common.EntityTypePlayer:
				ch = '☿' // 🦸
			case common.WorldTypeWall:
				ch = '⚀'
			case common.WorldTypePortal:
				ch = '0'
			case common.WorldTypePassage:
				ch = '*'
			case common.WorldTypeDoor:
				ch = 'П'
			case common.EntityTypeZombie:
				ch = '🧟'
			case common.EntityTypeVampire:
				ch = '🧛'
			case common.EntityTypeGhost:
				ch = '👻'
			case common.EntityTypeOgre:
				ch = '👹'
			case common.EntityTypeSnakeMage:
				ch = '🐍'
			case common.FoodTypePotatoes:
				ch = '🍟'
			case common.FoodTypeBread:
				ch = '🥖'
			case common.FoodTypeMeat:
				ch = '🍖'
			case common.FoodTypeMistery:
				ch = '🍄'
			case common.FoodTypeBeer:
				ch = '🍺'

			}
			screen.SetContent(ox+x, oy+y, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}
}
