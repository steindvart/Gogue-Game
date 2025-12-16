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
	FieldMarginY      = 3
	FieldMarginX      = 31
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
	fieldPanel *tview.TextView
}

func NewGame() *Game {
	game := &Game{
		container:  tview.NewFlex(),
		infoPanel:  tview.NewTextView(),
		fieldPanel: tview.NewTextView(),
	}

	// Настройка информационной панели (левая 1/3)
	game.infoPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Player Info").
		SetBackgroundColor(tcell.ColorBlack)

	// Настройка игрового поля (правые 2/3)
	game.fieldPanel.
		SetDynamicColors(true). // Включаем обработку цветных тегов
		SetWordWrap(false).     // Отключаем перенос строк
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

func (g *Game) GetContainer() *tview.Flex {
	return g.container
}

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

func (g *Game) UpdateGameField(field [][]common.GameEntityType) {
	g.fieldPanel.Clear()

	// Отступ по Y
	for i := 0; i < FieldMarginY; i++ {
		fmt.Fprintln(g.fieldPanel)
	}

	for y := range field {
		// Отступ по X
		for i := 0; i < FieldMarginX; i++ {
			fmt.Fprint(g.fieldPanel, " ")
		}

		// Отрисовка строки поля
		for x := range field[0] {
			ch := ' '
			color := "white"

			switch field[y][x] {
			case common.EntityTypePlayer:
				ch = '☿' // 🦸
				color = "yellow"
			case common.WorldTypeWall:
				ch = '█'
				color = "gray"
			case common.WorldTypePortal:
				ch = '0'
				color = "magenta"
			case common.WorldTypePassage,
				common.WorldTypeDoor:
				ch = '█'
				color = "silver"
			case common.EntityTypeZombie:
				ch = 'Z'
				color = "green"
			case common.EntityTypeVampire:
				ch = 'V'
				color = "red"
			case common.EntityTypeGhost:
				ch = 'G'
				color = "white"
			case common.EntityTypeOgre:
				ch = 'O'
				color = "brown"
			case common.EntityTypeSnakeMage:
				ch = 'S'
				color = "lime"
			case common.FoodTypePotatoes,
				common.FoodTypeBread,
				common.FoodTypeMeat,
				common.FoodTypeMistery,
				common.FoodTypeBeer:
				ch = 'ꕔ'
				color = "orange"
			case common.ElixirTypeStrength,
				common.ElixirTypeAgility,
				common.ElixirTypeDwarfism,
				common.ElixirTypeGiantism,
				common.ElixirTypeMystery:
				ch = 'ᗨ'
				color = "purple"
			case common.ScrollTypeStrength,
				common.ScrollTypeAgility,
				common.ScrollTypeUltimate,
				common.ScrollTypeMaxHealth,
				common.ScrollTypeMystery:
				ch = '⎕'
				color = "cyan"
			case common.Weapon:
				ch = 'T'
				color = "blue"
			}

			// Правильный формат символа с тегом: [tag]char без закрывающего тега между символами
			fmt.Fprintf(g.fieldPanel, "[%s]%c", color, ch)
		}

		// Сброс цвета в конце строки и перевод строки
		fmt.Fprintf(g.fieldPanel, "[-]\n")
	}
}

func (g *Game) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.container.SetInputCapture(capture)
}
