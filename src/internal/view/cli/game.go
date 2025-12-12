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
		// @todo - сделать нормальную заливку фона
		Box: tview.NewBox().SetBorder(true).SetTitle("Game").SetBackgroundColor(tcell.ColorBlack),
	}
}

func (*Game) SetFieldToScreen(screen tcell.Screen, field [][]common.GameEntityType, ox, oy int) {
	for y := range len(field) {
		for x := range len(field[0]) {
			ch := ' '
			style := tcell.StyleDefault
			switch field[y][x] {
			case common.EntityTypePlayer:
				ch = '☿' // 🦸
				style = style.Foreground(tcell.ColorYellow)
			case common.WorldTypeWall:
				ch = '█'
				style = style.Foreground(tcell.ColorGray)
			case common.WorldTypePortal:
				ch = '0'
			case common.WorldTypePassage:
				fallthrough
			case common.WorldTypeDoor:
				ch = '█'
				style = style.Foreground(tcell.ColorSilver)
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
				ch = 'ꕔ' //'🍟'
			case common.FoodTypeBread:
				ch = 'ꕔ' //'🥖'
			case common.FoodTypeMeat:
				ch = 'ꕔ' //'🍖'
			case common.FoodTypeMistery:
				ch = 'ꕔ' //'🍄'
			case common.FoodTypeBeer:
				ch = 'ꕔ' //'🍺'
			case common.ElixirTypeStrength:
				ch = 'ᗨ' //'🧡'
			case common.ElixirTypeAgility:
				ch = 'ᗨ' //'💚'
			case common.ElixirTypeDwarfism:
				ch = 'ᗨ' //'🩵'
			case common.ElixirTypeGiantism:
				ch = 'ᗨ' //'💙'
			case common.ElixirTypeMystery:
				ch = 'ᗨ' //'🖤'
			case common.ScrollTypeStrength:
				ch = '⎕' //'📙'
			case common.ScrollTypeAgility:
				ch = '⎕' //'📗'
			case common.ScrollTypeUltimate:
				ch = '⎕' //'📘'
			case common.ScrollTypeMaxHealth:
				ch = '⎕' //'📕'
			case common.ScrollTypeMystery:
				ch = '⎕' //'📓'
			case common.Weapon:
				ch = 'T' //'🗡'
			}

			screen.SetContent(ox+x, oy+y, ch, nil, style)
		}
	}
}
