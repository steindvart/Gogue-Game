package cli

import (
	"fmt"
	"gogue/internal/common"
	"gogue/internal/presentation/dto"

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
	ColorBlack   = "#000000"
	ColorWhite   = "#FFFFFF"
	ColorSkyBlue = "#87CEEB"

	// Цвета игрока
	ColorPlayer = "#FFD700" // Gold

	// Цвета стен и структур
	ColorWall    = "#808080" // Gray
	ColorFloor   = "#808080" // Gray
	ColorPassage = "#C0C0C0" // Silver
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
	ColorEffects  = ColorSkyBlue
)

type Game struct {
	container    *tview.Flex
	infoPanel    *tview.Flex // Контейнер для statsPanel и effectsPanel
	statsPanel   *tview.TextView
	effectsPanel *tview.TextView
	fieldPanel   *tview.TextView

	currentItemInfo  *dto.ItemInfo
	playerIsOnPortal bool
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

func (g *Game) GetRootPrimitive() *tview.Flex {
	return g.container
}

func (g *Game) UpdatePlayerInfo(info *dto.PlayerInfo) {
	g.updateStatsPanel(info)
	g.updateEffectsPanel(info.TemporaryEffects)
}

func (g *Game) ResetInteraction() {
	g.UpdateItemInfo(nil)
	g.SetPlayerIsOnPortal(false)
}

func (g *Game) SetPlayerIsOnPortal(val bool) {
	g.playerIsOnPortal = val
}

func (g *Game) UpdateItemInfo(info *dto.ItemInfo) {
	g.currentItemInfo = info
}

func (g *Game) updateStatsPanel(info *dto.PlayerInfo) {
	g.statsPanel.Clear()
	fmt.Fprintln(g.statsPanel)
	fmt.Fprintf(g.statsPanel, " [%s::b]HP:[-:-:-]       %.f/%-.f\n", ColorHP, info.Health, info.MaxHealth)
	fmt.Fprintf(g.statsPanel, " [%s::b]Strength:[-:-:-] %.f\n", ColorStrength, info.Strength)
	fmt.Fprintf(g.statsPanel, " [%s::b]Agility:[-:-:-]  %.f\n", ColorAgility, info.Agility)

	g.updateInteractionInfo()
}

func (g *Game) updateInteractionInfo() {
	if g.currentItemInfo != nil {
		g.renderItemInfo()
	}

	if g.playerIsOnPortal {
		printInteractPortalInfo(g.statsPanel)
	}
}

func printInteractPortalInfo(view *tview.TextView) {
	fmt.Fprintln(view)
	fmt.Fprintln(view, "───────────────────────────────────")
	fmt.Fprintf(view, " [%s::b]You are on the portal\n to the next level.\n[-:-:-]", ColorSkyBlue)
	fmt.Fprintln(view)
	printInteractionTips(view, true, false)
	fmt.Fprintf(view, " [%s::i]You can't go back...\n[-:-:-]", ColorSkyBlue)
}

func (g *Game) updateEffectsPanel(effects []dto.EffectInfo) {
	g.effectsPanel.Clear()
	fmt.Fprintf(g.effectsPanel, " [%s::b]ACTIVE EFFECTS:[-:-:-] (%d)\n\n", ColorEffects, len(effects))

	for i, effect := range effects {
		if i > 0 {
			fmt.Fprintln(g.effectsPanel, "───────────────────────────────────")
		}
		g.renderEffect(&effect)
	}
}

func (g *Game) renderEffect(effect *dto.EffectInfo) {
	fmt.Fprintf(g.effectsPanel, "[%s::b]Duration:[-:-:-] %d steps\n", ColorEffects, effect.DurationSteps)
	fmt.Fprintf(g.effectsPanel, "  HP:  %+6.1f\n", effect.HealthModify)
	fmt.Fprintf(g.effectsPanel, "  STR: %+6.1f\n", effect.StrengthModify)
	fmt.Fprintf(g.effectsPanel, "  AGI: %+6.1f\n", effect.AgilityModify)
}

func (g *Game) renderItemInfo() {
	fmt.Fprintln(g.statsPanel)
	fmt.Fprintln(g.statsPanel, "───────────────────────────────────")
	fmt.Fprintf(g.statsPanel, " [yellow::b]ITEM: [-:-:-]")

	var color string
	switch g.currentItemInfo.Type {
	case "Food":
		color = ColorFood
	case "Elixir":
		color = ColorElixir
	case "Scroll":
		color = ColorScroll
	case "Weapon":
		color = ColorWeapon
	case "Treasure":
		// @todo - добавить сокровище как предмет на карте
		// color = ColorTreasure
	default:
		color = ColorWhite
	}

	fmt.Fprintf(g.statsPanel, "[%s::b]%s %s[-:-:-]\n", color, g.currentItemInfo.Type, g.currentItemInfo.Name)

	// Для сокровищ показываем только стоимость
	if g.currentItemInfo.CanTake && !g.currentItemInfo.CanUse {
		fmt.Fprintf(g.statsPanel, " Value: %d\n", g.currentItemInfo.TreasureValue)
	} else {
		// Для остальных предметов показываем эффекты
		if g.currentItemInfo.DurationSteps > 0 {
			fmt.Fprintf(g.statsPanel, " Duration: %d steps\n", g.currentItemInfo.DurationSteps)
		}
		if g.currentItemInfo.MaxHealthModify != 0 {
			fmt.Fprintf(g.statsPanel, "  Max HP:  %+.1f\n", g.currentItemInfo.MaxHealthModify)
		}
		if g.currentItemInfo.HealthModify != 0 {
			fmt.Fprintf(g.statsPanel, "  HP:  %+8.1f\n", g.currentItemInfo.HealthModify)
		}
		if g.currentItemInfo.StrengthModify != 0 {
			fmt.Fprintf(g.statsPanel, "  STR: %+8.1f\n", g.currentItemInfo.StrengthModify)
		}
		if g.currentItemInfo.AgilityModify != 0 {
			fmt.Fprintf(g.statsPanel, "  AGI: %+8.1f\n", g.currentItemInfo.AgilityModify)
		}
	}

	fmt.Fprintln(g.statsPanel)
	printInteractionTips(g.statsPanel, g.currentItemInfo.CanUse, g.currentItemInfo.CanTake)
}

func printInteractionTips(view *tview.TextView, canUse, canTake bool) {
	if canUse && canTake {
		fmt.Fprintf(view, " [green]Use 'e' / Take 'r'[-]\n")
	} else if canUse {
		fmt.Fprintf(view, " [green]Use 'e'[-]\n")
	} else if canTake {
		fmt.Fprintf(view, " [green]Take 'r'[-]\n")
	}
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
		g.resetColorWithNewLine()
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
		// Полный формат: [foreground:background:modifier]char
		fmt.Fprintf(g.fieldPanel, "[%s:%s:-]%c", colorFg, colorBg, ch)
	}
}

func (g *Game) resetColorWithNewLine() {
	fmt.Fprintf(g.fieldPanel, "[-:-:-]\n")
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
		return '@', ColorPortal, ColorBlack
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
