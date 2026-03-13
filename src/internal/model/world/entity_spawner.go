package world

import (
	"errors"
	"fmt"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type RoomBasedEntitySpawner struct {
	foodTypes   []items.FoodType
	elixirTypes []items.ElixirType
	scrollTypes []items.ScrollType
	weaponTypes []items.WeaponType
}

func NewRoomBasedEntitySpawner() *RoomBasedEntitySpawner {
	return &RoomBasedEntitySpawner{
		foodTypes: []items.FoodType{
			items.FoodTypePotatoes,
			items.FoodTypeBread,
			items.FoodTypeMeat,
			items.FoodTypeMistery,
			items.FoodTypeBeer,
		},
		elixirTypes: []items.ElixirType{
			items.ElixirTypeStrength,
			items.ElixirTypeAgility,
			items.ElixirTypeDwarfism,
			items.ElixirTypeGiantism,
			items.ElixirTypeMystery,
		},
		scrollTypes: []items.ScrollType{
			items.ScrollTypeStrength,
			items.ScrollTypeAgility,
			items.ScrollTypeUltimate,
			items.ScrollTypeMaxHealth,
			items.ScrollTypeMystery,
		},
		weaponTypes: []items.WeaponType{
			items.WeaponTypeDagger,
			items.WeaponTypeSpear,
			items.WeaponTypeSword,
			items.WeaponTypeAxe,
			items.WeaponTypeMaul,
			items.WeaponTypeMystery,
		},
	}
}

func (s *RoomBasedEntitySpawner) SpawnEntities(
	rooms []Room,
	itemsConfig ItemSpawnConfig,
	enemiesConfig EnemySpawnConfig,
	random utils.Randomizer,
) (*SpawnedEntities, error) {
	if len(rooms) == 0 {
		return nil, errors.New("no rooms provided for entity spawning")
	}

	result := &SpawnedEntities{
		Items:   make([]primitives.Positional2D[int], 0),
		Enemies: make([]primitives.Positional2D[int], 0),
	}

	// Создаем изменяемые копии комнат для отслеживания занятых позиций
	roomsForSpawning := make([]*Room, len(rooms))
	for i := range rooms {
		roomsForSpawning[i] = &rooms[i]
	}

	// Спавним предметы
	var err error
	foods, err := s.spawnFoods(roomsForSpawning, itemsConfig.FoodsQuantity, random)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn foods: %w", err)
	}
	for i := range foods {
		result.Items = append(result.Items, &foods[i])
	}

	elixirs, err := s.spawnElixirs(roomsForSpawning, itemsConfig.ElixirsQuantity, random)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn elixirs: %w", err)
	}
	for i := range elixirs {
		result.Items = append(result.Items, &elixirs[i])
	}

	scrolls, err := s.spawnScrolls(roomsForSpawning, itemsConfig.ScrollsQuantity, random)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn scrolls: %w", err)
	}
	for i := range scrolls {
		result.Items = append(result.Items, &scrolls[i])
	}

	weapons, err := s.spawnWeapons(roomsForSpawning, itemsConfig.WeaponsQuantity, random)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn weapons: %w", err)
	}
	for i := range weapons {
		result.Items = append(result.Items, &weapons[i])
	}

	result.Enemies, err = s.spawnEnemies(roomsForSpawning, enemiesConfig.Quantity, random, enemiesConfig.AttributeMultiplier)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn enemies: %w", err)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) spawnFoods(
	rooms []*Room,
	quantity uint,
	random utils.Randomizer,
) ([]items.Food, error) {
	result := make([]items.Food, 0, quantity)

	for i := uint(0); i < quantity; i++ {
		room, pos, err := s.findAvailablePositionInSomeRoom(rooms, random)
		if err != nil {
			return nil, err
		}

		foodType := utils.GetRandomElement(random, s.foodTypes)
		itemBox := primitives.Box{
			Point: *pos,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		}

		food := items.NewFoodBuiltin(random, itemBox, foodType)
		result = append(result, *food)

		room.MarkOccupied(*pos)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) spawnElixirs(
	rooms []*Room,
	quantity uint,
	random utils.Randomizer,
) ([]items.Elixir, error) {
	result := make([]items.Elixir, 0, quantity)

	for i := uint(0); i < quantity; i++ {
		room, pos, err := s.findAvailablePositionInSomeRoom(rooms, random)
		if err != nil {
			return nil, err
		}

		elixirType := utils.GetRandomElement(random, s.elixirTypes)
		itemBox := primitives.Box{
			Point: *pos,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		}

		elixir := items.NewElixirBuiltin(random, itemBox, elixirType)
		result = append(result, *elixir)

		room.MarkOccupied(*pos)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) spawnScrolls(
	rooms []*Room,
	quantity uint,
	random utils.Randomizer,
) ([]items.Scroll, error) {
	result := make([]items.Scroll, 0, quantity)

	for i := uint(0); i < quantity; i++ {
		room, pos, err := s.findAvailablePositionInSomeRoom(rooms, random)
		if err != nil {
			return nil, err
		}

		scrollType := utils.GetRandomElement(random, s.scrollTypes)
		itemBox := primitives.Box{
			Point: *pos,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		}

		scroll := items.NewScrollBuiltin(random, itemBox, scrollType)
		result = append(result, *scroll)

		room.MarkOccupied(*pos)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) spawnWeapons(
	rooms []*Room,
	quantity uint,
	random utils.Randomizer,
) ([]items.Weapon, error) {
	result := make([]items.Weapon, 0, quantity)

	for i := uint(0); i < quantity; i++ {
		room, pos, err := s.findAvailablePositionInSomeRoom(rooms, random)
		if err != nil {
			return nil, err
		}

		weaponType := utils.GetRandomElement(random, s.weaponTypes)
		itemBox := primitives.Box{
			Point: *pos,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		}

		weapon := items.NewWeaponBuiltin(random, itemBox, weaponType)
		result = append(result, *weapon)

		room.MarkOccupied(*pos)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) spawnEnemies(
	rooms []*Room,
	quantity uint,
	random utils.Randomizer,
	attributeMultiplier float64,
) ([]primitives.Positional2D[int], error) {
	result := make([]primitives.Positional2D[int], 0, quantity)

	for i := uint(0); i < quantity; i++ {
		room, pos, err := s.findAvailablePositionInSomeRoom(rooms, random)
		if err != nil {
			return nil, err
		}

		enemyType := utils.GetRandomElement(random, entities.EnemyTypes)
		box := &primitives.Box{
			Point: *pos,
			Size:  primitives.Size2D[uint]{Width: 1, Height: 1},
		}

		var enemy primitives.Positional2D[int]

		// Пример как мы можем помещать разные типы врагов в одну через общий интерфейс
		switch enemyType {
		case entities.EnemyTypeZombie:
			enemy = entities.NewZombie(box)
		case entities.EnemyTypeVampire:
			enemy = entities.NewVampire(box)
		case entities.EnemyTypeGhost:
			enemy = entities.NewGhost(box)
		case entities.EnemyTypeOgre:
			enemy = entities.NewOgre(box)
		case entities.EnemyTypeSnakeMage:
			enemy = entities.NewSnakeMage(box)
		case entities.EnemyTypeMimic:
			enemy = entities.NewMimic(box)
		}

		// Применяем множитель атрибутов для повышения сложности на поздних уровнях
		if attributeMultiplier > 1.0 {
			if provider, ok := enemy.(entities.EnemyProvider); ok {
				e := provider.GetEnemy()
				e.ScaleAttributes(attributeMultiplier)
			}
		}

		result = append(result, enemy)
		room.MarkOccupied(*pos)
	}

	return result, nil
}

func (s *RoomBasedEntitySpawner) findAvailablePositionInSomeRoom(
	rooms []*Room,
	random utils.Randomizer,
) (*Room, *primitives.Point2D[int], error) {
	indexes := random.Perm(len(rooms))

	for _, idx := range indexes {
		room := rooms[idx]

		// Не размещаем в стартовой комнате
		if room.Type == RoomTypeStart {
			continue
		}

		if room.GetCountFreePosition() <= 0 {
			continue
		}

		pos, err := room.GetRandomFreePosition(random)
		if err != nil {
			continue
		}

		return room, pos, nil
	}

	return nil, nil, errors.New("no available positions in rooms for entity spawning")
}
