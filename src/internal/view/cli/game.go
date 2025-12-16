package cli

import (
	"fmt"
	"gogue/internal/common"

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
	// Общие цвета
	ColorBlack = "#000000" // Black
	ColorWhite = "#FFFFFF" // White

	// Цвета игрока
	ColorPlayer = "#FFD700" // Gold

	// Цвета стен и структур
	ColorWall    = "#808080" // Gray
	ColorFloor   = "#808080" // Gray
	ColorPassage = "#C0C0C0" // Silber
	ColorPortal  = "#FF00FF" // Purple

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

type EffectInfo struct {
	DurationSteps  int
	HealthModify   float64
	StrengthModify float64
	AgilityModify  float64
}

type PlayerInfo struct {
	Health           float64
	MaxHealth        float64
	Strength         float64
	Agility          float64
	TemporaryEffects []EffectInfo
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

	game.setupPanels()
	game.setupLayout()

	return game
}

func (g *Game) setupPanels() {
	g.statsPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Stats").
		SetBackgroundColor(tcell.ColorBlack)

	g.effectsPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Effects").
		SetBackgroundColor(tcell.ColorBlack)

	g.infoPanel.
		SetDirection(tview.FlexRow).
		AddItem(g.statsPanel, 0, 1, false).
		AddItem(g.effectsPanel, 0, 1, false).
		SetBackgroundColor(tcell.ColorBlack)

	g.fieldPanel.
		SetDynamicColors(true).
		SetWordWrap(false).
		SetBorder(true).
		SetTitle("Game").
		SetBackgroundColor(tcell.ColorBlack)
}

func (g *Game) setupLayout() {
	// Горизонтальная компоновка: информационная панель + игровое поле
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(g.infoPanel, MaxInfoPanelWidth, 0, false).
		AddItem(g.fieldPanel, 0, 1, true)
	horizontalFlex.SetBackgroundColor(tcell.ColorBlack)

	// Пустой бокс для нижнего отступа
	gapBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)

	// Вертикальная компоновка: панели + нижний отступ
	g.container.
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, MaxPanelHeight, 0, true).
		AddItem(gapBox, 0, 1, false).
		SetBackgroundColor(tcell.ColorBlack)
}

func (g *Game) GetContainer() *tview.Flex {
	return g.container
}

func (g *Game) UpdatePlayerInfo(info *PlayerInfo) {
	g.updateStatsPanel(info)
	g.updateEffectsPanel(info.TemporaryEffects)
}

func (g *Game) updateStatsPanel(info *PlayerInfo) {
	g.statsPanel.Clear()
	fmt.Fprintln(g.statsPanel)
	fmt.Fprintf(g.statsPanel, " [%s::b]HP:[-:-:-]       %.f/%-.f\n", ColorHP, info.Health, info.MaxHealth)
	fmt.Fprintf(g.statsPanel, " [%s::b]Strength:[-:-:-] %.f\n", ColorStrength, info.Strength)
	fmt.Fprintf(g.statsPanel, " [%s::b]Agility:[-:-:-]  %.f\n", ColorAgility, info.Agility)
}

func (g *Game) updateEffectsPanel(effects []EffectInfo) {
	g.effectsPanel.Clear()
	fmt.Fprintf(g.effectsPanel, " [%s::b]ACTIVE EFFECTS:[-:-:-] (%d)\n\n", ColorEffects, len(effects))

	for i, effect := range effects {
		if i > 0 {
			fmt.Fprintln(g.effectsPanel, "───────────────────────────────────")
		}
		g.renderEffect(&effect)
	}
}

func (g *Game) renderEffect(effect *EffectInfo) {
	fmt.Fprintf(g.effectsPanel, "[%s::b]Duration:[-:-:-] %d steps\n", ColorEffects, effect.DurationSteps)
	fmt.Fprintf(g.effectsPanel, "  HP:  %+6.1f\n", effect.HealthModify)
	fmt.Fprintf(g.effectsPanel, "  STR: %+6.1f\n", effect.StrengthModify)
	fmt.Fprintf(g.effectsPanel, "  AGI: %+6.1f\n", effect.AgilityModify)
}

func (g *Game) UpdateGameField(field [][]common.GameEntityType) {
	g.fieldPanel.Clear()
	g.renderField(field)
}

func (g *Game) renderMarginsY(margin int) {
	for i := 0; i < margin; i++ {
		fmt.Fprintln(g.fieldPanel)
	}
}

func (g *Game) renderField(field [][]common.GameEntityType) {
	g.renderMarginsY(FieldMarginY)

	for y := range field {
		g.renderMarginX(FieldMarginX)
		g.renderRow(field, y)
		g.resetColorAndNewLine()
	}
}

func (g *Game) renderMarginX(margin int) {
	for i := 0; i < margin; i++ {
		fmt.Fprint(g.fieldPanel, " ")
	}
}

func (g *Game) renderRow(field [][]common.GameEntityType, y int) {
	for x := range field[0] {
		ch, colorFg, colorBg := g.getCellAppearance(field[y][x])

		// @todo - сделать разные фоны для разных типов поверхностей
		// @todo - сделать так, чтобы цвет фона зависел от типа поверхности под объектом (разделить поле на два слоя?)
		// Формат с фоном: [foreground:background]char
		fmt.Fprintf(g.fieldPanel, "[%s:%s]%c", colorFg, colorBg, ch)
	}
}

func (g *Game) resetColorAndNewLine() {
	fmt.Fprintf(g.fieldPanel, "[-]\n")
}

func (g *Game) getCellAppearance(entityType common.GameEntityType) (rune, string, string) {
	switch entityType {
	case common.EntityTypePlayer:
		return '☿', ColorPlayer, ColorBlack
	case common.WorldTypeWall:
		return ' ', ColorWall, ColorWall
	case common.WorldTypeRoomFloor:
		return '.', ColorFloor, ColorBlack
	case common.WorldTypePortal:
		return '◎', ColorPortal, ColorBlack
	case common.WorldTypePassage, common.WorldTypeDoor:
		return ' ', ColorPassage, ColorPassage
	case common.EntityTypeZombie:
		return 'Z', ColorZombie, ColorBlack
	case common.EntityTypeVampire:
		return 'V', ColorVampire, ColorBlack
	case common.EntityTypeGhost:
		return 'G', ColorGhost, ColorBlack
	case common.EntityTypeOgre:
		return 'O', ColorOgre, ColorBlack
	case common.EntityTypeSnakeMage:
		return 'S', ColorSnakeMage, ColorBlack
	case common.FoodTypePotatoes, common.FoodTypeBread, common.FoodTypeMeat,
		common.FoodTypeMistery, common.FoodTypeBeer:
		return 'ð', ColorFood, ColorBlack
	case common.ElixirTypeStrength, common.ElixirTypeAgility, common.ElixirTypeDwarfism,
		common.ElixirTypeGiantism, common.ElixirTypeMystery:
		return '¶', ColorElixir, ColorBlack
	case common.ScrollTypeStrength, common.ScrollTypeAgility, common.ScrollTypeUltimate,
		common.ScrollTypeMaxHealth, common.ScrollTypeMystery:
		return '!', ColorScroll, ColorBlack
	case common.Weapon:
		return 'ƒ', ColorWeapon, ColorBlack
	default:
		return ' ', ColorWhite, ColorBlack
	}
}

func (g *Game) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.container.SetInputCapture(capture)
}
