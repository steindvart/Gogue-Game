package cli

import (
	"fmt"
	"gogue/internal/common"
	"gogue/internal/presentation/dto"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MaxPanelHeight    = 40
	MaxInfoPanelWidth = 40
	MaxEffectsToShow  = 3
	FieldMarginY      = 3
	FieldMarginX      = 17

	MaxBackpackItemsNum = 9
)

const (
	// Общие цвета
	ColorBlack   = "#000000"
	ColorWhite   = "#FFFFFF"
	ColorSkyBlue = "#87CEEB"
	ColorAqua    = "#0194A7"

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

type SecondInfoViewMode int

const (
	SecondInfoViewModeEffects SecondInfoViewMode = iota
	SecondInfoViewModeBackpack
)

type BackpackTab int

const (
	BackpackTabWeapons BackpackTab = iota
	BackpackTabFood
	BackpackTabElixirs
	BackpackTabScrolls
)

type Game struct {
	container      *tview.Flex
	leftInfoPanel  *tview.Flex // Контейнер для statsPanel и effectsPanel/backpackPanel
	statsPanel     *tview.TextView
	effectsPanel   *tview.TextView
	backpackPanel  *tview.TextView
	rightInfoPanel *tview.Flex // Контейнер для legendPanel
	legendPanel    *tview.TextView
	fieldPanel     *tview.TextView

	currentItemInfo  *dto.ItemInfo
	playerIsOnPortal bool

	// Состояние рюкзака
	secondInfoViewMode SecondInfoViewMode
	currentTab         BackpackTab
	playerBackpackInfo *dto.BackpackInfo
}

func NewGame() *Game {
	game := &Game{
		container:          tview.NewFlex(),
		rightInfoPanel:     tview.NewFlex(),
		statsPanel:         tview.NewTextView(),
		effectsPanel:       tview.NewTextView(),
		backpackPanel:      tview.NewTextView(),
		leftInfoPanel:      tview.NewFlex(),
		legendPanel:        tview.NewTextView(),
		fieldPanel:         tview.NewTextView(),
		secondInfoViewMode: SecondInfoViewModeEffects,
		currentTab:         BackpackTabWeapons,
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
		SetTitle("Effects (switch 'b')").
		SetBackgroundColor(tcell.ColorBlack)

	g.backpackPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Backpack (switch 'b')").
		SetBackgroundColor(tcell.ColorBlack)

	g.legendPanel.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Legend").
		SetBackgroundColor(tcell.ColorBlack)

	// Левая панель: Stats + Effects/Backpack
	g.leftInfoPanel.
		SetDirection(tview.FlexRow).
		AddItem(g.statsPanel, 0, 1, false).
		AddItem(g.effectsPanel, 0, 1, false).
		SetBackgroundColor(tcell.ColorBlack)

	g.rightInfoPanel.
		SetDirection(tview.FlexRow).
		AddItem(g.legendPanel, 0, 1, false).
		SetBackgroundColor(tcell.ColorBlack)

	g.fieldPanel.
		SetDynamicColors(true).
		SetWordWrap(false).
		SetBorder(true).
		SetTitle("Game").
		SetBackgroundColor(tcell.ColorBlack)

	g.initializeLegend()
}

func (g *Game) setupLayout() {
	// Горизонтальная компоновка: левая панель + игровое поле + правая панель (легенда)
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(g.leftInfoPanel, MaxInfoPanelWidth, 0, false).
		AddItem(g.fieldPanel, 0, 1, true).
		AddItem(g.rightInfoPanel, MaxInfoPanelWidth, 0, false)
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

func (g *Game) initializeLegend() {
	g.legendPanel.Clear()

	fmt.Fprintf(g.legendPanel, " [%s::b]CONTROLS:[-:-:-]\n", ColorSkyBlue)
	fmt.Fprintln(g.legendPanel, " ↑←↓→ or 'WASD' - Move")
	fmt.Fprintln(g.legendPanel, " E - Use item")
	fmt.Fprintln(g.legendPanel, " R - Take item")
	fmt.Fprintln(g.legendPanel, " B - Toggle Backpack")
	fmt.Fprintln(g.legendPanel, " Z/X/C/V - Tabs")
	fmt.Fprintln(g.legendPanel, " 0-9 - Select item")
	fmt.Fprintln(g.legendPanel, " ESC - Exit game")
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]SYMBOLS:[-:-:-]\n", ColorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]☿[-] - Player\n", ColorPlayer)
	fmt.Fprintf(g.legendPanel, " [%s]@[-] - Portal\n", ColorPortal)
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]ENEMIES:[-:-:-]\n", ColorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]Z[-] - Zombie\n", ColorZombie)
	fmt.Fprintf(g.legendPanel, " [%s]V[-] - Vampire\n", ColorVampire)
	fmt.Fprintf(g.legendPanel, " [%s]G[-] - Ghost\n", ColorGhost)
	fmt.Fprintf(g.legendPanel, " [%s]O[-] - Ogre\n", ColorOgre)
	fmt.Fprintf(g.legendPanel, " [%s]S[-] - Snake Mage\n", ColorSnakeMage)
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]ITEMS:[-:-:-]\n", ColorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]ð[-] - Food\n", ColorFood)
	fmt.Fprintf(g.legendPanel, " [%s]¶[-] - Elixir\n", ColorElixir)
	fmt.Fprintf(g.legendPanel, " [%s]![-] - Scroll\n", ColorScroll)
	fmt.Fprintf(g.legendPanel, " [%s]ƒ[-] - Weapon\n", ColorWeapon)
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

	sortedEffects := make([]dto.EffectInfo, len(effects))
	copy(sortedEffects, effects)

	sort.Slice(sortedEffects, func(i, j int) bool {
		return sortedEffects[i].DurationSteps < sortedEffects[j].DurationSteps
	})

	if len(sortedEffects) >= MaxEffectsToShow {
		sortedEffects = sortedEffects[:MaxEffectsToShow]
	}

	for i, effect := range sortedEffects {
		if i > 0 {
			fmt.Fprintln(g.effectsPanel, "───────────────────────────────────")
		}
		printEffect(g.effectsPanel, &effect)
	}
}

func printEffect(view *tview.TextView, effect *dto.EffectInfo) {
	if effect.DurationSteps > 0 {
		fmt.Fprintf(view, " [%s::b]Duration:[-:-:-] %d steps\n", ColorAqua, effect.DurationSteps)
	}
	if effect.MaxHealthModify != 0 {
		fmt.Fprintf(view, "  Max HP:  %+.1f\n", effect.MaxHealthModify)
	}
	if effect.HealthModify != 0 {
		fmt.Fprintf(view, "  HP:  %+8.1f\n", effect.HealthModify)
	}
	if effect.StrengthModify != 0 {
		fmt.Fprintf(view, "  STR: %+8.1f\n", effect.StrengthModify)
	}
	if effect.AgilityModify != 0 {
		fmt.Fprintf(view, "  AGI: %+8.1f\n", effect.AgilityModify)
	}
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
		printEffect(g.statsPanel, g.currentItemInfo.EffectInfo)
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

func (g *Game) UpdateBackpackInfo(backpack *dto.BackpackInfo) {
	g.playerBackpackInfo = backpack
	g.updateBackpackPanel()
}

func (g *Game) ToggleViewMode() {
	if g.secondInfoViewMode == SecondInfoViewModeEffects {
		g.secondInfoViewMode = SecondInfoViewModeBackpack
		g.switchToBackpackPanel()
	} else {
		g.secondInfoViewMode = SecondInfoViewModeEffects
		g.switchToEffectsPanel()
	}
}

func (g *Game) switchToBackpackPanel() {
	g.leftInfoPanel.RemoveItem(g.effectsPanel)
	g.leftInfoPanel.AddItem(g.backpackPanel, 0, 1, false)
	g.updateBackpackPanel()
}

func (g *Game) switchToEffectsPanel() {
	g.leftInfoPanel.RemoveItem(g.backpackPanel)
	g.leftInfoPanel.AddItem(g.effectsPanel, 0, 1, false)
}

func (g *Game) SetBackpackTab(tab BackpackTab) {
	g.currentTab = tab
	if g.secondInfoViewMode == SecondInfoViewModeBackpack {
		g.updateBackpackPanel()
	}
}

func (g *Game) GetCurrentBackpackTab() BackpackTab {
	return g.currentTab
}

func (g *Game) updateBackpackPanel() {
	g.backpackPanel.Clear()

	if g.playerBackpackInfo == nil {
		fmt.Fprintf(g.backpackPanel, " [%s::b]BACKPACK IS EMPTY[-:-:-]\n", ColorSkyBlue)
		return
	}

	// Отображаем заголовок с текущей вкладкой
	var tabName string
	var items []*dto.ItemInfo

	switch g.currentTab {
	case BackpackTabWeapons:
		tabName = "WEAPONS"
		items = g.playerBackpackInfo.Weapons
	case BackpackTabFood:
		tabName = "FOOD"
		items = g.playerBackpackInfo.Foods
	case BackpackTabElixirs:
		tabName = "ELIXIRS"
		items = g.playerBackpackInfo.Elixirs
	case BackpackTabScrolls:
		tabName = "SCROLLS"
		items = g.playerBackpackInfo.Scrolls
	}

	fmt.Fprintf(g.backpackPanel, " [%s::b]%s:[-:-:-] (%d/%d)\n", ColorSkyBlue, tabName, g.playerBackpackInfo.ItemsNum, g.playerBackpackInfo.Capacity)
	fmt.Fprintln(g.backpackPanel)

	// Отображаем список предметов
	if len(items) == 0 {
		if g.currentTab == BackpackTabWeapons {
			fmt.Fprintf(g.backpackPanel, " 0: [%s::i]put current weapon to backpack[-:-:-]\n\n", ColorSkyBlue)
		}
		fmt.Fprintf(g.backpackPanel, " [%s::i]No items in this category[-:-:-]\n", ColorSkyBlue)
	} else {
		g.renderBackpackItems(items)
	}
}

func (g *Game) renderBackpackItems(items []*dto.ItemInfo) {
	startIndex := 0
	if g.currentTab == BackpackTabWeapons {
		startIndex = 0 // Для оружия отображаем 0-9
	} else {
		startIndex = 1 // Для остальных предметов 1-9
	}

	for i, item := range items {
		if i == 0 {
			fmt.Fprintf(g.backpackPanel, " %d: [-:-:i]put current weapon to backpack[-:-:-]\n", i)
			continue
		}

		if i >= MaxBackpackItemsNum {
			break
		}

		num := startIndex + i
		icon, color := g.getItemIconAndColor(item.Type)

		// Отображаем номер, иконку, тип и название
		fmt.Fprintf(g.backpackPanel, " %d: [%s]%c %s %s[-:-:-]\n", num, color, icon, item.Type, item.Name)

		// Отображаем эффекты предмета
		if item.EffectInfo != nil {
			g.getShortEffectInfoAsStrings(item.EffectInfo)
		}

		if i < len(items)-1 {
			effects := g.getShortEffectInfoAsStrings(item.EffectInfo)
			if len(effects) > 0 {
				fmt.Fprintf(g.backpackPanel, "(%s)\n", strings.Join(effects, ", "))
			}
		}
	}
}

func (g *Game) getItemIconAndColor(itemType string) (rune, string) {
	switch itemType {
	case "Weapon":
		return 'ƒ', ColorWeapon
	case "Food":
		return 'ð', ColorFood
	case "Elixir":
		return '¶', ColorElixir
	case "Scroll":
		return '!', ColorScroll
	default:
		return '?', ColorWhite
	}
}

func (g *Game) getShortEffectInfoAsStrings(effect *dto.EffectInfo) []string {
	var effects []string

	if effect.DurationSteps > 0 {
		effects = append(effects, fmt.Sprintf("Dur: %d", effect.DurationSteps))
	}
	if effect.MaxHealthModify != 0 {
		effects = append(effects, fmt.Sprintf("Max HP: %+.0f", effect.MaxHealthModify))
	}
	if effect.HealthModify != 0 {
		effects = append(effects, fmt.Sprintf("HP: %+.0f", effect.HealthModify))
	}
	if effect.StrengthModify != 0 {
		effects = append(effects, fmt.Sprintf("STR: %+.0f", effect.StrengthModify))
	}
	if effect.AgilityModify != 0 {
		effects = append(effects, fmt.Sprintf("AGI: %+.0f", effect.AgilityModify))
	}

	return effects
}
