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

const (
	// Цвета игрока
	ColorPlayerFg = "#FFD700" // Gold

	// Цвета стен и структур
	ColorWall   = "#808080" // Gray
	ColorDoor   = "#C0C0C0" // Silber
	ColorPortal = "#FF00FF" // Purple

	// Цвета врагов
	ColorZombie    = "#00FF00" // Green
	ColorVampire   = "#FF0000" // Red
	ColorGhost     = "#E0E0E0" // Light gray
	ColorOgre      = "#8B4513" // Saddle brown
	ColorSnakeMage = "#32CD32" // Lime green

	// Цвета предметов
	ColorFood   = "#FFA500" // Orange
	ColorElixir = "#9370DB" // Medium purple
	ColorScroll = "#00CED1" // Dark turquoise
	ColorWeapon = "#4169E1" // Royal blue

	// Цвета текста (для информационной панели)
	ColorHP       = "#FFD700" // Gold
	ColorStrength = "#FF6347" // Tomato
	ColorAgility  = "#90EE90" // Light green
	ColorEffects  = "#87CEEB" // Sky blue
)

type PlayerInfo struct {
	Health           float64
	MaxHealth        float64
	Strength         float64
	Agility          float64
	TemporaryEffects []*primitives.Effect
}

type Game struct {
	container    *tview.Flex
	infoPanel    *tview.Flex // Контейнер для statsPanel и effectsPanel
	statsPanel   *tview.TextView
	effectsPanel *tview.TextView
	fieldPanel   *tview.TextView
}

func NewGame() *Game {
	game := &Game{
		container:    tview.NewFlex(),
		infoPanel:    tview.NewFlex(),
		statsPanel:   tview.NewTextView(),
		effectsPanel: tview.NewTextView(),
		fieldPanel:   tview.NewTextView(),
	}

	// Настройка панели статистики
	game.statsPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Stats").
		SetBackgroundColor(tcell.ColorBlack)

	// Настройка панели эффектов
	game.effectsPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Effects").
		SetBackgroundColor(tcell.ColorBlack)

	// Настройка контейнера информационной панели (вертикальное расположение)
	game.infoPanel.
		SetDirection(tview.FlexRow).
		AddItem(game.statsPanel, 0, 1, false).
		AddItem(game.effectsPanel, 0, 1, false).
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
	// Обновление панели статистики
	g.statsPanel.Clear()
	fmt.Fprintf(g.statsPanel, "\n")
	fmt.Fprintf(g.statsPanel, " [%s::b]HP:[-:-:-]       %.f/%-.f\n", ColorHP, info.Health, info.MaxHealth)
	fmt.Fprintf(g.statsPanel, " [%s::b]Strength:[-:-:-] %.f\n", ColorStrength, info.Strength)
	fmt.Fprintf(g.statsPanel, " [%s::b]Agility:[-:-:-]  %.f\n", ColorAgility, info.Agility)

	// Обновление панели эффектов
	g.effectsPanel.Clear()
	fmt.Fprintf(g.effectsPanel, " [%s::b]ACTIVE EFFECTS:[-:-:-] (%d)\n\n", ColorEffects, len(info.TemporaryEffects))

	if len(info.TemporaryEffects) > 0 {
		for i, effect := range info.TemporaryEffects {
			if i > 0 {
				fmt.Fprintln(g.effectsPanel, "───────────────────────────────────")
			}
			fmt.Fprintf(g.effectsPanel, "[%s::b]Duration:[-:-:-] %d steps\n", ColorEffects, effect.Duration.Steps)
			fmt.Fprintf(g.effectsPanel, "  HP:  %+6.1f\n", effect.Attributes.Health)
			fmt.Fprintf(g.effectsPanel, "  STR: %+6.1f\n", effect.Attributes.Strength)
			fmt.Fprintf(g.effectsPanel, "  AGI: %+6.1f\n", effect.Attributes.Agility)
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
			color := "#FFFFFF"   // Белый по умолчанию
			bgcolor := "#000000" // Чёрный фон для всех

			switch field[y][x] {
			case common.EntityTypePlayer:
				ch = '☿'
				color = ColorPlayerFg
			case common.WorldTypeWall:
				ch = '█'
				color = ColorWall
			case common.WorldTypeRoomFloor:
				ch = ' '
			case common.WorldTypePortal:
				ch = '◎'
				color = ColorPortal
			case common.WorldTypePassage,
				common.WorldTypeDoor:
				ch = '█'
				color = ColorDoor
			case common.EntityTypeZombie:
				ch = 'Z'
				color = ColorZombie
			case common.EntityTypeVampire:
				ch = 'V'
				color = ColorVampire
			case common.EntityTypeGhost:
				ch = 'G'
				color = ColorGhost
			case common.EntityTypeOgre:
				ch = 'O'
				color = ColorOgre
			case common.EntityTypeSnakeMage:
				ch = 'S'
				color = ColorSnakeMage
			case common.FoodTypePotatoes,
				common.FoodTypeBread,
				common.FoodTypeMeat,
				common.FoodTypeMistery,
				common.FoodTypeBeer:
				ch = 'ð'
				color = ColorFood
			case common.ElixirTypeStrength,
				common.ElixirTypeAgility,
				common.ElixirTypeDwarfism,
				common.ElixirTypeGiantism,
				common.ElixirTypeMystery:
				ch = '¶'
				color = ColorElixir
			case common.ScrollTypeStrength,
				common.ScrollTypeAgility,
				common.ScrollTypeUltimate,
				common.ScrollTypeMaxHealth,
				common.ScrollTypeMystery:
				ch = '!'
				color = ColorScroll
			case common.Weapon:
				ch = 'ƒ'
				color = ColorWeapon
			}

			// @todo - сделать разные фоны для разных типов поверхностей
			// @todo - сделать так, чтобы цвет фона зависел от типа поверхности под объектом (разделить поле на два слоя?)
			// Формат с фоном: [foreground:background]char
			fmt.Fprintf(g.fieldPanel, "[%s:%s]%c", color, bgcolor, ch)
		}

		// Сброс цвета в конце строки и перевод строки
		fmt.Fprintf(g.fieldPanel, "[-]\n")
	}
}

func (g *Game) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.container.SetInputCapture(capture)
}
