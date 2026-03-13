package world

import (
	"math"

	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type Level struct {
	// Геометрия
	Rooms        []Room
	Passages     []Passage
	FinishPortal primitives.Box

	// Сущности, хранящиеся непосредственно в уровне
	Player  *entities.Player
	Enemies []primitives.Positional2D[int]
	Items   []primitives.Positional2D[int]

	// Метаданные
	Number uint

	// Конфигурация и зависимости
	config        LevelConfig
	random        utils.Randomizer
	roomGen       RoomGenerator
	passageGen    PassageGenerator
	entitySpawner EntitySpawner
	playerSpawner PlayerSpawner
	fieldRenderer FieldRenderer
	FogOfWar      *FogOfWar
}

type CollisionType int

const (
	CollisionTypeNone CollisionType = iota
	CollisionTypeEnemy
	CollisionTypeItem
	CollisionTypeTeleport
)

func NewLevelWithDefaults(random utils.Randomizer, mapSize primitives.Size2D[uint]) *Level {
	cfg := DefaultLevelConfig(mapSize)
	return newLevelWithComponents(
		random,
		cfg,
		1,
		NewGridRoomGenerator(),
		NewConnectionTreePassageGenerator(),
		NewRoomBasedEntitySpawner(),
		NewStartRoomPlayerSpawner(),
		NewDefaultFieldRenderer(),
	)
}

func newLevelWithComponents(
	random utils.Randomizer,
	cfg LevelConfig,
	levelNumber uint,
	roomGen RoomGenerator,
	passageGen PassageGenerator,
	entitySpawner EntitySpawner,
	playerSpawner PlayerSpawner,
	fieldRenderer FieldRenderer,
) *Level {
	return &Level{
		config:        cfg,
		random:        random,
		roomGen:       roomGen,
		passageGen:    passageGen,
		entitySpawner: entitySpawner,
		playerSpawner: playerSpawner,
		Number:        levelNumber,
		fieldRenderer: fieldRenderer,
		FogOfWar:      NewFogOfWar(int(cfg.MapSize.Width), int(cfg.MapSize.Height)),
		Enemies:       []primitives.Positional2D[int]{},
		Items:         []primitives.Positional2D[int]{},
	}
}

func (l *Level) ProcessTurns(turns uint32) []entities.AttackResult {
	l.Player.ProcessTurns(turns)
	return l.processEnemyTurns()
}

// processEnemyTurns обрабатывает ходы всех врагов за один игровой тик.
//
// Логика для каждого врага:
//  1. Вычисляется расстояние (Чебышёва) до игрока.
//  2. Если расстояние <= HostilityRadius И враг имеет прямую видимость (LoS) до игрока -
//     враг переходит в режим преследования (IsChasing = true).
//  3. Если враг преследует - строится путь до игрока через FindPathBFS.
//  4. Если враг преследует, но потерял LoS - сбрасываем преследование.
//  5. Если первый шаг пути == позиция игрока - враг атакует, но не двигается.
//  6. Если первый шаг пути занят другим врагом - враг пропускает ход.
//  7. Иначе враг перемещается на первый шаг пути.
//  8. Если путь не найден - враг продолжает бродить по своему паттерну (idle move).
//  9. Если враг не преследует - двигается по своему idle-паттерну.
func (l *Level) processEnemyTurns() []entities.AttackResult {
	if l.Player == nil || len(l.Enemies) == 0 {
		return nil
	}

	var attackResults []entities.AttackResult

	playerPos := l.Player.GetPosition()
	mapW := int(l.config.MapSize.Width)
	mapH := int(l.config.MapSize.Height)

	// Получаем полное поле один раз для всех врагов (используется для LoS-проверки и idle-движения).
	fullField := l.GetFullField(mapW, mapH)

	// Типы клеток, по которым враги могут перемещаться:
	// пол, проходы, двери, предметы (еда, зелья, свитки, оружие), портал.
	walkable := []int{
		int(common.WorldTypeRoomFloor),
		int(common.WorldTypePassage),
		int(common.WorldTypeDoor),
		int(common.WorldTypePortal),
		int(common.Food),
		int(common.Elixir),
		int(common.Scroll),
		int(common.Weapon),
	}

	idleMover := NewEnemyIdleMover(l.random)

	for _, positionalEnemy := range l.Enemies {
		provider, ok := positionalEnemy.(entities.EnemyProvider)
		if !ok {
			continue
		}
		enemy := provider.GetEnemy()

		enemyPos := enemy.GetPosition()
		dist := chebyshevDistance(enemyPos, playerPos)

		// Переходим в режим преследования только если:
		// 1) Игрок в зоне враждебности (HostilityRadius)
		// 2) Враг имеет прямую видимость (Line of Sight) до игрока
		if dist <= int(enemy.HostilityRadius) && HasLineOfSight(fullField, enemyPos, playerPos) {
			enemy.IsChasing = true

			// Ghost становится всегда видимым при переходе в режим преследования
			if ghost, ok := positionalEnemy.(*entities.Ghost); ok {
				ghost.IsVisible = true
			}

			// Mimic сбрасывает маскировку при переходе в режим преследования
			if mimic, ok := positionalEnemy.(*entities.Mimic); ok {
				mimic.IsDisguised = false
			}
		}

		// Если враг преследует, но потерял видимость - сбрасываем преследование.
		// Это позволяет игроку "оторваться" от врага, скрывшись за стеной.
		if enemy.IsChasing && !HasLineOfSight(fullField, enemyPos, playerPos) {
			enemy.IsChasing = false
		}

		// Если не преследуем - выполняем idle-движение по паттерну типа врага
		if !enemy.IsChasing {
			idleMover.MoveIdle(positionalEnemy, fullField, l.Rooms, l.Enemies)
			continue
		}

		// --- Преследование ---

		// Строим навигационное поле: из полного поля карты,
		// но блокируем клетки, занятые другими врагами.
		navField := l.buildEnemyNavigationField(mapW, mapH, positionalEnemy)

		// BFS
		// 4 направления: N, E, S, W
		// Но можно сделать 8 направлений, включая диагонали, если нужно более "естественное" движение
		// Однако так станет значительно сложнее играть. И для этого нужно будет вводить возможность движения по диагонали для
		// игрока. Что в целом просто реализовать, но из-за этого усложниться управление - нужно будет биндить ещё 4 клавиши.
		// Поэтому пока оставим только 4 направления, чтобы враги двигались по "квадратной" сетке так же как и игрок.
		directions := []primitives.Point2D[int]{
			{X: 0, Y: -1}, // North
			// {X: 1, Y: -1},  // North-East
			{X: 1, Y: 0}, // East
			// {X: 1, Y: 1},   // South-East
			{X: 0, Y: 1}, // South
			// {X: -1, Y: 1},  // South-West
			{X: -1, Y: 0}, // West
		}

		path, err := FindPathBFS(enemyPos, playerPos, navField, walkable, directions)
		if err != nil || len(path) == 0 {
			// Путь не найден - продолжаем бродить по idle-паттерну
			idleMover.MoveIdle(positionalEnemy, fullField, l.Rooms, l.Enemies)
			continue
		}

		nextStep := path[0]

		// Если следующий шаг - позиция игрока: атакуем, но не двигаемся
		if nextStep == playerPos {
			attackResult := l.processEnemyAttack(positionalEnemy, enemy)
			if attackResult != nil {
				attackResults = append(attackResults, *attackResult)
			}
			continue
		}

		// Проверяем, что клетка не занята другим врагом
		if l.isEnemyAtPosition(nextStep, positionalEnemy) {
			continue
		}

		// Перемещаем врага
		enemy.SetPosition(nextStep)
	}

	return attackResults
}

// processEnemyAttack обрабатывает атаку конкретного врага на игрока
// с учётом уникальных модификаторов типа врага.
//
// Модификаторы:
//   - Ogre: если IsResting - пропускает ход (отдыхает), снимает флаг и
//     гарантированно наносит удар (игнорирует evasion) на следующий ход.
//     После каждой обычной атаки устанавливает IsResting = true.
//   - SnakeMage: при успешном попадании с вероятностью 30% усыпляет игрока на 1 ход.
//   - Прочие враги: стандартная атака через Character.Attack.
func (l *Level) processEnemyAttack(
	positionalEnemy primitives.Positional2D[int],
	enemy *entities.Enemy,
) *entities.AttackResult {
	name := enemyTypeName(positionalEnemy)

	// --- Огр: боевой цикл Ready → Resting → Enraged → Resting → ... ---
	if ogre, ok := positionalEnemy.(*entities.Ogre); ok {
		switch ogre.CombatPhase {
		case entities.OgrePhaseResting:
			// Отдых - пропуск хода. Переход в фазу ярости.
			ogre.CombatPhase = entities.OgrePhaseEnraged
			return nil

		case entities.OgrePhaseEnraged:
			// Гарантированная контратака (без проверки evasion).
			damage := enemy.Character.MakeDamage()
			l.Player.Character.TakeDamage(damage)

			ogre.CombatPhase = entities.OgrePhaseResting

			return &entities.AttackResult{
				AttackerName:        name,
				DefenderName:        "Player",
				Damage:              damage,
				DefenderHealthAfter: l.Player.Character.Attributes.Health,
				DefenderMaxHealth:   l.Player.Character.Attributes.MaxHealth,
				DefenderKilled:      !l.Player.Character.IsAlive(),
				Guaranteed:          true,
			}

		default:
			// OgrePhaseReady - обычная атака, после которой огр уходит на отдых.
			result := enemy.Character.Attack(l.Player.Character, l.random)
			result.AttackerName = name
			result.DefenderName = "Player"
			ogre.CombatPhase = entities.OgrePhaseResting
			return &result
		}
	}

	// --- Змей-маг: шанс усыпления ---
	if _, ok := positionalEnemy.(*entities.SnakeMage); ok {
		result := enemy.Character.Attack(l.Player.Character, l.random)
		result.AttackerName = name
		result.DefenderName = "Player"

		const sleepChance = 0.3
		if !result.Evaded && l.random.Float64() < sleepChance {
			l.Player.Stun(1)
			result.AppliedStun = true
		}
		return &result
	}

	// --- Стандартная атака ---
	result := enemy.Character.Attack(l.Player.Character, l.random)
	result.AttackerName = name
	result.DefenderName = "Player"
	return &result
}

// buildEnemyNavigationField создаёт навигационную карту для конкретного врага.
// Берёт полное поле уровня (со всеми сущностями) и блокирует клетки,
// занятые другими врагами, заменяя их на EntityTypeNone (непроходимый тип).
// Клетка текущего врага не блокируется - это стартовая позиция поиска.
func (l *Level) buildEnemyNavigationField(w, h int, currentEnemy primitives.Positional2D[int]) [][]int {
	fullField := l.GetFullField(w, h)
	navField := make([][]int, h)

	for y := 0; y < h; y++ {
		navField[y] = make([]int, w)
		for x := 0; x < w; x++ {
			navField[y][x] = int(fullField[y][x])
		}
	}

	// Блокируем клетки других врагов
	for _, other := range l.Enemies {
		if other == currentEnemy {
			continue
		}
		pos := other.GetPosition()
		if pos.X >= 0 && pos.X < w && pos.Y >= 0 && pos.Y < h {
			navField[pos.Y][pos.X] = int(common.EntityTypeNone)
		}
	}

	// Позицию текущего врага помечаем как проходимую.
	// В fullField она отрендерена как EntityTypeZombie/Vampire/etc.,
	// которые не входят в walkable - без этого BFS вернёт ErrInvalidStart.
	curPos := currentEnemy.GetPosition()
	if curPos.X >= 0 && curPos.X < w && curPos.Y >= 0 && curPos.Y < h {
		navField[curPos.Y][curPos.X] = int(common.WorldTypeRoomFloor)
	}

	// Позицию игрока помечаем как проходимую (чтобы BFS нашёл путь к нему).
	// В fullField она отрендерена как EntityTypePlayer, который тоже не в walkable.
	// Используем тип RoomFloor, т.к. он гарантированно в walkable.
	if l.Player != nil {
		pp := l.Player.GetPosition()
		if pp.X >= 0 && pp.X < w && pp.Y >= 0 && pp.Y < h {
			navField[pp.Y][pp.X] = int(common.WorldTypeRoomFloor)
		}
	}

	return navField
}

// isEnemyAtPosition проверяет, есть ли враг (кроме excludeEnemy) на указанной позиции.
func (l *Level) isEnemyAtPosition(pos primitives.Point2D[int], excludeEnemy primitives.Positional2D[int]) bool {
	for _, other := range l.Enemies {
		if other == excludeEnemy {
			continue
		}
		if other.GetPosition() == pos {
			return true
		}
	}
	return false
}

// enemyTypeName возвращает читаемое имя типа врага для отображения в UI.
func enemyTypeName(positional primitives.Positional2D[int]) string {
	switch positional.(type) {
	case *entities.Zombie:
		return "Zombie"
	case *entities.Vampire:
		return "Vampire"
	case *entities.Ghost:
		return "Ghost"
	case *entities.Ogre:
		return "Ogre"
	case *entities.SnakeMage:
		return "Snake Mage"
	case *entities.Mimic:
		return "Mimic"
	default:
		return "Unknown"
	}
}

// chebyshevDistance вычисляет расстояние Чебышёва (шахматное расстояние) между двумя точками.
// Это максимум из абсолютных разностей координат: max(|dx|, |dy|).
// Используется для определения, попадает ли игрок в зону враждебности врага
// при 8-направленном движении.
func chebyshevDistance(a, b primitives.Point2D[int]) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	absX := int(math.Abs(float64(dx)))
	absY := int(math.Abs(float64(dy)))
	if absX > absY {
		return absX
	}
	return absY
}

// Generate генерирует геометрию, сущности и игрока
func (l *Level) Generate() error {
	err := l.generateEnvironment()
	if err != nil {
		return err
	}

	player, err := l.playerSpawner.SpawnPlayer(l.Rooms, l.random)
	if err != nil {
		return err
	}
	l.Player = player

	return nil
}

// GenerateWithExistingPlayer генерирует мир, но использует переданного игрока
func (l *Level) GenerateWithExistingPlayer(player *entities.Player) error {
	l.Number++

	// Пересчитываем сложность для нового уровня
	l.config.ScaleForLevel(l.Number)

	err := l.generateEnvironment()
	if err != nil {
		return err
	}

	l.Player = player
	startPos, err := l.playerSpawner.GetStartPosition(l.Rooms, l.random)
	if err != nil {
		return err
	}
	l.Player.SetPosition(*startPos)
	return nil
}

func (l *Level) generateEnvironment() error {
	// Сбрасываем туман войны при генерации нового уровня
	l.FogOfWar.Reset()

	// Геометрия
	rooms, finishPortal, err := l.roomGen.GenerateRooms(l.config, l.random)
	if err != nil {
		return err
	}
	l.Rooms = rooms
	l.FinishPortal = finishPortal

	passages, err := l.passageGen.GeneratePassages(l.Rooms, l.config, l.random)
	if err != nil {
		return err
	}
	l.Passages = passages

	// Сущности
	spawned, err := l.entitySpawner.SpawnEntities(l.Rooms, l.config.ItemsSpawnConfig, l.config.EnemiesSpawnConfig, l.random)
	if err != nil {
		return err
	}

	l.Items = spawned.Items
	l.Enemies = spawned.Enemies

	return nil
}

func (l *Level) CheckEntityCollision(pos primitives.Point2D[int]) CollisionType {
	if l.isCollisionWithEnemy(pos) {
		return CollisionTypeEnemy
	}
	if l.isCollisionWithItem(pos) {
		return CollisionTypeItem
	}
	if l.isCollisionWithTeleport(pos) {
		return CollisionTypeTeleport
	}

	return CollisionTypeNone
}

func (l *Level) IsCollisionWithBorders(pos primitives.Point2D[int]) bool {
	if l.isCollisionWithMapBorders(pos) {
		return true
	}
	if l.isCollisionWithRoomsWall(pos) {
		return true
	}
	if !isInSomeRoom(pos, l.Rooms) && !l.isCollisionWithPassages(pos) {
		return true
	}

	return false
}

func (l *Level) isCollisionWithMapBorders(pos primitives.Point2D[int]) bool {
	if pos.X < 0 || pos.Y < 0 {
		return true
	}
	if pos.X >= int(l.config.MapSize.Width) || pos.Y >= int(l.config.MapSize.Height) {
		return true
	}
	return false
}

func (l *Level) isCollisionWithRoomsWall(pos primitives.Point2D[int]) bool {
	for _, room := range l.Rooms {
		if isInRoom(pos, room) && isCollisionWithRoomWall(pos, room) {
			return true
		}
	}
	return false
}

func isInSomeRoom(pos primitives.Point2D[int], rooms []Room) bool {
	for _, room := range rooms {
		if isInRoom(pos, room) {
			return true
		}
	}
	return false
}

func isInRoom(pos primitives.Point2D[int], room Room) bool {
	leftEndX := room.Box.Point.X + 1
	rightEndX := leftEndX + int(room.Box.Size.Width) - 3
	topEndY := room.Box.Point.Y + 1
	downEndY := topEndY + int(room.Box.Size.Height) - 3

	return (pos.X >= leftEndX && pos.X <= rightEndX) && (pos.Y >= topEndY && pos.Y <= downEndY)
}

func isCollisionWithRoomWall(pos primitives.Point2D[int], room Room) bool {
	if isCollisionWithDoors(pos, room.Doors) {
		return false
	}

	leftEndX := room.Box.Point.X
	rightEndX := leftEndX + int(room.Box.Size.Width)
	topEndY := room.Box.Point.Y
	downEndY := topEndY + int(room.Box.Size.Height)

	if (pos.X == leftEndX || pos.X == rightEndX) || (pos.Y == topEndY || pos.Y == downEndY) {
		return true
	}
	return false
}

func isCollisionWithDoors(pos primitives.Point2D[int], doors []primitives.Point2D[int]) bool {
	for _, door := range doors {
		if pos == door {
			return true
		}
	}
	return false
}

func (l *Level) isCollisionWithPassages(newPos primitives.Point2D[int]) bool {
	for _, passage := range l.Passages {
		if isInPassage(newPos, passage) {
			return true
		}
	}
	return false
}

func isInPassage(pos primitives.Point2D[int], passage Passage) bool {
	if pos == passage.DoorOne || pos == passage.DoorTwo {
		return true
	}
	for _, wayPoint := range passage.Way {
		if pos == wayPoint {
			return true
		}
	}
	return false
}

func (l *Level) isCollisionWithTeleport(delta primitives.Point2D[int]) bool {
	return delta == l.FinishPortal.Point
}

func (l *Level) GetEnemyAtPosition(pos primitives.Point2D[int]) primitives.Positional2D[int] {
	for _, enemy := range l.Enemies {
		if pos == enemy.GetPosition() {
			return enemy
		}
	}
	return nil
}

func (l *Level) isCollisionWithEnemy(pos primitives.Point2D[int]) bool {
	return l.GetEnemyAtPosition(pos) != nil
}

func (l *Level) isCollisionWithItem(pos primitives.Point2D[int]) bool {
	for _, item := range l.Items {
		if pos == item.GetPosition() {
			return true
		}
	}
	return false
}

// MakeCurrentField создаёт двумерное представление карты уровня с учётом тумана войны
func (l *Level) MakeCurrentField(w, h int) [][]common.GameEntityType {
	// Сначала рендерим полное поле
	fullField := l.GetFullField(w, h)

	// Если игрока нет, возвращаем пустое поле (всё скрыто туманом войны)
	if l.Player == nil {
		return utils.CreateEmpty2DSlice[common.GameEntityType](h, w)
	}

	// Фильтруем поле с учётом тумана войны и радиуса обзора игрока
	return l.FogOfWar.ApplyFogOfWar(fullField, l.Player.GetPosition(), l.Player.ViewRadius)
}

// GetFullField возвращает двумерное представление карты уровня бещ учёта тумана войны
func (l *Level) GetFullField(w, h int) [][]common.GameEntityType {
	return l.fieldRenderer.RenderField(w, h, l)
}

func (l *Level) PlayerUseItemAtPosition(pos primitives.Point2D[int]) {
	item := l.GetItemAtPosition(pos)
	if item == nil {
		return
	}

	// @todo - кривая логика, нужно декомпозировать и перенести в Player
	if weapon, ok := item.(*items.Weapon); ok {
		_ = l.Player.EquipWeapon(weapon)
	} else if usableItem, ok := item.(items.Usable); ok {
		l.Player.Character.Use(usableItem)
	}

	l.RemoveItem(item)
}

func (l *Level) PlayerTakeItemAtPosition(pos primitives.Point2D[int]) error {
	item := l.GetItemAtPosition(pos)
	if item == nil {
		return nil
	}

	if takeableItem, ok := item.(items.Takeable); ok {
		if err := l.Player.Backpack.AddItem(takeableItem); err != nil {
			return err
		}
		l.RemoveItem(item)
	}

	return nil
}

func (l *Level) GetItemAtPosition(pos primitives.Point2D[int]) primitives.Positional2D[int] {
	for _, item := range l.Items {
		if item.GetPosition() == pos {
			return item
		}
	}
	return nil
}

func (l *Level) RemoveItem(item primitives.Positional2D[int]) {
	for i, it := range l.Items {
		if it == item {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			return
		}
	}
}

func (l *Level) AddItem(item primitives.Positional2D[int]) {
	l.Items = append(l.Items, item)
}

func (l *Level) RemoveEnemy(enemy primitives.Positional2D[int]) {
	for i, it := range l.Enemies {
		if it == enemy {
			l.Enemies = append(l.Enemies[:i], l.Enemies[i+1:]...)
			return
		}
	}
}

type Attacker interface {
	Attack(defender *entities.Character, rnd utils.Randomizer) entities.AttackResult
}

// Attack выполняет атаку и возвращает результат.
// Имена атакующего и защитника заполняются вызывающим кодом.
// Перед стандартной атакой проверяются уникальные модификаторы защитника.
func (l *Level) Attack(attacker Attacker, defender primitives.Positional2D[int]) *entities.AttackResult {
	if attacker == nil || defender == nil {
		return nil
	}

	charProvider, ok := defender.(entities.CharacterProvider)
	if !ok {
		return nil
	}

	// Модификатор Вампира: абсолютное уклонение (первый удар - гарантированный промах).
	if vampire, ok := defender.(*entities.Vampire); ok {
		if vampire.AbsoluteEvasions > 0 {
			vampire.AbsoluteEvasions--
			return &entities.AttackResult{
				Evaded:              true,
				DefenderHealthAfter: vampire.Character.Attributes.Health,
				DefenderMaxHealth:   vampire.Character.Attributes.MaxHealth,
			}
		}
	}

	result := attacker.Attack(charProvider.GetCharacter(), l.random)
	return &result
}

// GenerateTreasure генерирует сокровище с процентами выпадения,
// масштабированными в соответствии с текущим уровнем подземелья.
func (l *Level) GenerateTreasure() (*items.Treasure, error) {
	drop := l.config.GetTreasureDropConfig(l.Number)

	item, err := items.GetTreasureType(l.random, drop.GoldPercent, drop.GemPercent, drop.ArtifactPercent, drop.MysteryPercent)
	if err != nil {
		return nil, err
	}

	return items.NewTreasureBuiltin(
		l.random,
		primitives.Box{
			Point: primitives.Point2D[int]{X: 0, Y: 0},
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		},
		item), nil
}
