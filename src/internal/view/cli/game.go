package cli

import (
	"fmt"
	"gogue/internal/common"
	"gogue/internal/model/primitives"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MaxPanelHeight    = 40
	MaxInfoPanelWidth = 40
)

type PlayerInfo struct {
	Health           float64
	MaxHealth        float64
	Strength         float64
	Agility          float64
	TemporaryEffects []*primitives.Effect
}

type Game struct {
	container  *tview.Flex
	infoPanel  *tview.TextView
	fieldPanel *tview.Box
}

func NewGame() *Game {
	game := &Game{
		container:  tview.NewFlex(),
		infoPanel:  tview.NewTextView(),
		fieldPanel: tview.NewBox(),
	}

	// Настройка информационной панели (левая треть)
	game.infoPanel.
		SetBorder(true).
		SetTitle("Player Info").
		SetBackgroundColor(tcell.ColorBlack)

	// Настройка игрового поля (правые две трети)
	game.fieldPanel.
		SetBorder(true).
		SetTitle("Game").
		SetBackgroundColor(tcell.ColorBlack)

	// Горизонтальная компоновка: информационная панель (фиксированная ширина) + игровое поле
	horizontalFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	horizontalFlex.
		AddItem(game.infoPanel, MaxInfoPanelWidth, 0, false).
		AddItem(game.fieldPanel, 0, 1, true)

	horizontalFlex.SetBackgroundColor(tcell.ColorBlack)

	// Создаём пустой primitives.Box с тёмным фоном для выравнивания элементов
	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	// Вертикальная компоновка: верхний отступ + панели (ограниченная высота) + нижний отступ
	game.container.SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, MaxPanelHeight, 0, true). // Основные панели (фиксированная высота)
		AddItem(gapBox, 0, 1, false).                     // Нижний отступ (растягивается)
		SetBackgroundColor(tcell.ColorBlack)

	return game
}

// GetContainer возвращает контейнер для встраивания в tview приложение
func (g *Game) GetContainer() *tview.Flex {
	return g.container
}

// UpdatePlayerInfo обновляет информацию об игроке на информационной панели
func (g *Game) UpdatePlayerInfo(info PlayerInfo) {
	g.infoPanel.Clear()
	fmt.Fprintf(g.infoPanel, "[yellow::b]HP:[-:-:-] %.1f / %.1f\n", info.Health, info.MaxHealth)
	fmt.Fprintf(g.infoPanel, "[red::b]Strength:[-:-:-] %.1f\n", info.Strength)
	fmt.Fprintf(g.infoPanel, "[green::b]Agility:[-:-:-] %.1f\n\n", info.Agility)

	if len(info.TemporaryEffects) > 0 {
		fmt.Fprintf(g.infoPanel, "[cyan::b]Active Effects:[-:-:-]\n")
		for _, effect := range info.TemporaryEffects {
			fmt.Fprintf(g.infoPanel, "  • Duration: %d steps\n", effect.Duration.Steps)
			if effect.Attributes.Health != 0 {
				fmt.Fprintf(g.infoPanel, "    HP: %+.1f\n", effect.Attributes.Health)
			}
			if effect.Attributes.Strength != 0 {
				fmt.Fprintf(g.infoPanel, "    Str: %+.1f\n", effect.Attributes.Strength)
			}
			if effect.Attributes.Agility != 0 {
				fmt.Fprintf(g.infoPanel, "    Agi: %+.1f\n", effect.Attributes.Agility)
			}
		}
	}
}

// SetFieldDrawFunc устанавливает функцию отрисовки игрового поля
func (g *Game) SetFieldDrawFunc(draw func(screen tcell.Screen, x, y, width, height int) (int, int, int, int)) {
	g.fieldPanel.SetDrawFunc(draw)
}

// SetInputCapture устанавливает обработчик ввода для всего контейнера
func (g *Game) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.container.SetInputCapture(capture)
}

// SetFieldToScreen отрисовывает игровое поле на экране
func (g *Game) SetFieldToScreen(screen tcell.Screen, field [][]common.GameEntityType, ox, oy int) {
	const marginX, marginY = 31, 4

	for y := range len(field) - marginY {
		for x := range len(field[0]) - marginX {
			ch := ' '
			style := tcell.StyleDefault.Background(tcell.ColorBlack)
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

			screen.SetContent(ox+x+marginX, oy+y+marginY, ch, nil, style)
		}
	}
}
