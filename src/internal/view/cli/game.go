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
	// Символы игровых сущностей
	SymbolPlayer    = '☿'
	SymbolPortal    = '@'
	SymbolFloor     = '.'
	SymbolEmpty     = ' '
	SymbolZombie    = 'Z'
	SymbolVampire   = 'V'
	SymbolGhost     = 'G'
	SymbolOgre      = 'O'
	SymbolSnakeMage = 'S'
	SymbolFood      = 'ð'
	SymbolElixir    = '¶'
	SymbolScroll    = '!'
	SymbolWeapon    = 'ƒ'
	SymbolUnknown   = '?'
)

const (
	// Базовые цвета
	colorBlack      = "#000000"
	colorWhite      = "#FFFFFF"
	colorGray       = "#808080"
	colorLightGray  = "#E0E0E0"
	colorSilver     = "#C0C0C0"
	colorGold       = "#FFD700"
	colorOrange     = "#FFA500"
	colorRed        = "#FF0000"
	colorTomato     = "#FF6347"
	colorGreen      = "#00FF00"
	colorLimeGreen  = "#32CD32"
	colorLightGreen = "#90EE90"
	colorPurple     = "#FF00FF"
	colorMedPurple  = "#9370DB"
	colorRoyalBlue  = "#4169E1"
	colorSkyBlue    = "#87CEEB"
	colorAqua       = "#0194A7"
	colorTurquoise  = "#00CED1"
	colorBrown      = "#8B4513"
)

const (
	// Цвета игрока
	ColorPlayer = colorGold

	// Цвета стен и структур
	ColorWall    = colorGray
	ColorFloor   = colorGray
	ColorPassage = colorSilver
	ColorPortal  = colorPurple

	// Цвета врагов
	ColorZombie    = colorGreen
	ColorVampire   = colorRed
	ColorGhost     = colorLightGray
	ColorOgre      = colorBrown
	ColorSnakeMage = colorLimeGreen

	// Цвета предметов
	ColorFood   = colorOrange
	ColorElixir = colorMedPurple
	ColorScroll = colorTurquoise
	ColorWeapon = colorRoyalBlue

	// Цвета текста (для информационной панели)
	ColorLevel    = "#9999FF"
	ColorHP       = colorGold
	ColorStrength = colorTomato
	ColorAgility  = colorLightGreen
	ColorEffects  = colorSkyBlue
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
	SecondInfoViewMode SecondInfoViewMode
	currentTab         BackpackTab
	playerBackpackInfo *dto.BackpackInfo
	backpackDropMode   bool

	infoErrorMessage string
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
		SecondInfoViewMode: SecondInfoViewModeEffects,
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

	fmt.Fprintf(g.legendPanel, " [%s::b]CONTROLS:[-:-:-]\n", colorSkyBlue)
	fmt.Fprintln(g.legendPanel, " ↑←↓→ or 'wasd' - Movement")
	fmt.Fprintln(g.legendPanel, " e - Use item")
	fmt.Fprintln(g.legendPanel, " r - Take item")
	fmt.Fprintln(g.legendPanel, " ESC - Exit game")
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]BACKPACK:[-:-:-]\n", colorSkyBlue)
	fmt.Fprintln(g.legendPanel, " b - Switch Effects/Backpack panels")
	fmt.Fprintln(g.legendPanel, " n - Switch Use/Drop mode")
	fmt.Fprintln(g.legendPanel, " 0-9 - Select item or action")
	fmt.Fprintln(g.legendPanel, " z - Weapon tab")
	fmt.Fprintln(g.legendPanel, " x - Food tab")
	fmt.Fprintln(g.legendPanel, " c - Elixir tab")
	fmt.Fprintln(g.legendPanel, " v - Scroll tab")
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]SYMBOLS:[-:-:-]\n", colorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Player\n", ColorPlayer, SymbolPlayer)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Portal\n", ColorPortal, SymbolPortal)
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]ENEMIES:[-:-:-]\n", colorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Zombie\n", ColorZombie, SymbolZombie)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Vampire\n", ColorVampire, SymbolVampire)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Ghost\n", ColorGhost, SymbolGhost)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Ogre\n", ColorOgre, SymbolOgre)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Snake Mage\n", ColorSnakeMage, SymbolSnakeMage)
	fmt.Fprintln(g.legendPanel)

	fmt.Fprintf(g.legendPanel, " [%s::b]ITEMS:[-:-:-]\n", colorSkyBlue)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Food\n", ColorFood, SymbolFood)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Elixir\n", ColorElixir, SymbolElixir)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Scroll\n", ColorScroll, SymbolScroll)
	fmt.Fprintf(g.legendPanel, " [%s]%c[-] - Weapon\n", ColorWeapon, SymbolWeapon)
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
	g.ClearInfoMessages()
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
	fmt.Fprintf(g.statsPanel, " [%s::b]Level:[-:-:-]    %d\n", ColorLevel, info.LevelNumber)
	fmt.Fprintf(g.statsPanel, " [%s::b]HP:[-:-:-]       %.f/%.f\n", ColorHP, info.Health, info.MaxHealth)
	fmt.Fprintf(g.statsPanel, " [%s::b]Strength:[-:-:-] %.f\n", ColorStrength, info.Strength)
	fmt.Fprintf(g.statsPanel, " [%s::b]Agility:[-:-:-]  %.f\n", ColorAgility, info.Agility)
	if info.WeaponInfo == nil {
		fmt.Fprintf(g.statsPanel, " [%s::b]Weapon [-:-:-]   None\n", ColorWeapon)
	} else {
		fmt.Fprintf(
			g.statsPanel, " [%s::b]Weapon %s:[-:-:-] %s\n",
			ColorWeapon, info.WeaponInfo.Type, strings.Join(g.getShortEffectInfoAsStrings(info.WeaponInfo.EffectInfo), ", "),
		)
	}

	g.updateInteractionInfo()
}

func (g *Game) updateInteractionInfo() {
	if g.currentItemInfo != nil {
		g.renderItemInfo()
	}

	if g.playerIsOnPortal {
		printInteractPortalInfo(g.statsPanel)
	}

	// @todo - небольшой баг с отрисовкой: если уже есть сообщение об ошибке,
	// то мелькнёт ещё одно аналогичное сообщение.
	// Не хочется делать "костыль" для решения проблемы, нужно подумать над этим отдельно.
	if g.infoErrorMessage != "" {
		printErrorInfoMessage(g.statsPanel, g.infoErrorMessage)
	}
}

func printErrorInfoMessage(view *tview.TextView, msg string) {
	fmt.Fprintln(view)
	fmt.Fprintln(view, "───────────────────────────────────")
	fmt.Fprintf(view, "[%s::b]Error:[-:-:-]\n %s\n", colorTomato, msg)
	fmt.Fprintln(view)
}

func (g *Game) SetInfoErrorMessage(msg string) {
	g.infoErrorMessage = msg
	g.updateInteractionInfo()
}

func (g *Game) ClearInfoMessages() {
	g.infoErrorMessage = ""
}

func printInteractPortalInfo(view *tview.TextView) {
	fmt.Fprintln(view)
	fmt.Fprintln(view, "───────────────────────────────────")
	fmt.Fprintf(view, " [%s::b]You are on the portal\n to the next level.\n[-:-:-]", colorSkyBlue)
	fmt.Fprintln(view)
	printInteractionTips(view, true, false)
	fmt.Fprintf(view, " [%s::i]You can't go back...\n[-:-:-]", colorSkyBlue)
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
		fmt.Fprintf(view, " [%s::b]Duration:[-:-:-] %d steps\n", colorAqua, effect.DurationSteps)
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
		color = colorWhite
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
		return SymbolPlayer, ColorPlayer, colorBlack
	case common.WorldTypeWall:
		return SymbolEmpty, ColorWall, ColorWall
	case common.WorldTypeRoomFloor:
		return SymbolFloor, ColorFloor, colorBlack
	case common.WorldTypePortal:
		return SymbolPortal, ColorPortal, colorBlack
	case common.WorldTypePassage, common.WorldTypeDoor:
		return SymbolEmpty, ColorPassage, ColorPassage
	case common.EntityTypeZombie:
		return SymbolZombie, ColorZombie, colorBlack
	case common.EntityTypeVampire:
		return SymbolVampire, ColorVampire, colorBlack
	case common.EntityTypeGhost:
		return SymbolGhost, ColorGhost, colorBlack
	case common.EntityTypeOgre:
		return SymbolOgre, ColorOgre, colorBlack
	case common.EntityTypeSnakeMage:
		return SymbolSnakeMage, ColorSnakeMage, colorBlack
	case common.Food:
		return SymbolFood, ColorFood, colorBlack
	case common.Elixir:
		return SymbolElixir, ColorElixir, colorBlack
	case common.Scroll:
		return SymbolScroll, ColorScroll, colorBlack
	case common.Weapon:
		return SymbolWeapon, ColorWeapon, colorBlack
	default:
		return SymbolEmpty, colorWhite, colorBlack
	}
}

func (g *Game) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	g.container.SetInputCapture(capture)
}

func (g *Game) UpdateBackpackInfo(backpack *dto.BackpackInfo) {
	g.playerBackpackInfo = backpack
	g.updateBackpackPanel()
}

func (g *Game) ToggleSecondInfoViewMode() {
	if g.SecondInfoViewMode == SecondInfoViewModeEffects {
		g.SecondInfoViewMode = SecondInfoViewModeBackpack
		g.switchToBackpackPanel()
	} else {
		g.SecondInfoViewMode = SecondInfoViewModeEffects
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
	if g.SecondInfoViewMode == SecondInfoViewModeBackpack {
		g.updateBackpackPanel()
	}
}

func (g *Game) GetCurrentBackpackTab() BackpackTab {
	return g.currentTab
}

func (g *Game) SetBackpackDropMode(dropMode bool) {
	g.backpackDropMode = dropMode
	if g.SecondInfoViewMode == SecondInfoViewModeBackpack {
		g.updateBackpackPanel()
	}
}

func (g *Game) updateBackpackPanel() {
	g.backpackPanel.Clear()

	if g.playerBackpackInfo == nil {
		fmt.Fprintf(g.backpackPanel, " [%s::b]PLAYER DON'T HAVE BACKPACK[-:-:-]\n", colorSkyBlue)
		return
	}

	// Отображаем заголовок с текущей вкладкой
	var tabName string
	var items []*dto.ItemInfo
	color := colorWhite

	switch g.currentTab {
	case BackpackTabWeapons:
		tabName = "WEAPONS"
		items = g.playerBackpackInfo.Weapons
		color = ColorWeapon
	case BackpackTabFood:
		tabName = "FOOD"
		items = g.playerBackpackInfo.Foods
		color = ColorFood
	case BackpackTabElixirs:
		tabName = "ELIXIRS"
		items = g.playerBackpackInfo.Elixirs
		color = ColorElixir
	case BackpackTabScrolls:
		tabName = "SCROLLS"
		items = g.playerBackpackInfo.Scrolls
		color = ColorScroll
	}

	fmt.Fprintf(g.backpackPanel, "[%s::b]Common capacity: %d/%d\n", colorBrown, g.playerBackpackInfo.ItemsNum, g.playerBackpackInfo.Capacity)
	fmt.Fprintf(g.backpackPanel, " [%s::b]%s:[-:-:-] %d\n", color, tabName, len(items))

	modeStr := "USE"
	modeColor := colorLightGreen
	if g.backpackDropMode {
		modeStr = "DROP"
		modeColor = colorTomato
	}
	fmt.Fprintf(g.backpackPanel, " [%s::b]Mode: %s[-:-:-]\n", modeColor, modeStr)
	fmt.Fprintln(g.backpackPanel)

	// Отображаем список предметов
	if g.currentTab == BackpackTabWeapons {
		if g.backpackDropMode {
			fmt.Fprintf(g.backpackPanel, " 0: [%s::i]drop current weapon[-:-:-]\n", colorSkyBlue)
		} else {
			fmt.Fprintf(g.backpackPanel, " 0: [%s::i]put current weapon to backpack[-:-:-]\n", colorSkyBlue)
		}
	}

	if len(items) == 0 {
		fmt.Fprintf(g.backpackPanel, " [%s::i]No items in this category[-:-:-]\n", colorSkyBlue)
	} else {
		g.renderBackpackItems(items)
	}
}

func (g *Game) renderBackpackItems(items []*dto.ItemInfo) {
	startIndex := 1

	for i, item := range items {
		if i >= MaxBackpackItemsNum {
			break
		}

		num := startIndex + i
		icon, color := g.getItemIconAndColor(item.Type)

		// Отображаем номер, иконку, тип и название
		fmt.Fprintf(g.backpackPanel, " %d: [%s]%c %s %s[-:-:-]\n", num, color, icon, item.Type, item.Name)

		// Отображаем эффекты предмета
		if item.EffectInfo != nil {
			effects := strings.Join(g.getShortEffectInfoAsStrings(item.EffectInfo), ", ")
			fmt.Fprintf(g.backpackPanel, "(%s)\n", effects)
		}
	}
}

func (g *Game) getItemIconAndColor(itemType string) (rune, string) {
	switch itemType {
	case "Weapon":
		return SymbolWeapon, ColorWeapon
	case "Food":
		return SymbolFood, ColorFood
	case "Elixir":
		return SymbolElixir, ColorElixir
	case "Scroll":
		return SymbolScroll, ColorScroll
	default:
		return SymbolUnknown, colorWhite
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
