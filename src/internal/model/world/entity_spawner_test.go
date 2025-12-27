package world

import (
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"math/rand"
	"testing"
)

func TestRoomBasedEntitySpawner_SpawnEntities(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	// Генерируем комнаты
	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	// Генерируем сущности
	spawner := NewRoomBasedEntitySpawner()
	entities, err := spawner.SpawnEntities(rooms, config.ItemsSpawnConfig, config.EnemiesSpawnConfig, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Подсчитываем предметы по типам
	foodCount, elixirCount, scrollCount, weaponCount := 0, 0, 0, 0
	for _, item := range entities.Items {
		switch item.(type) {
		case *items.Food:
			foodCount++
		case *items.Elixir:
			elixirCount++
		case *items.Scroll:
			scrollCount++
		case *items.Weapon:
			weaponCount++
		}
	}

	// Проверяем что сущности созданы
	if foodCount != int(config.ItemsSpawnConfig.FoodsQuantity) {
		t.Errorf("Expected %d foods, got %d", config.ItemsSpawnConfig.FoodsQuantity, foodCount)
	}
	if elixirCount != int(config.ItemsSpawnConfig.ElixirsQuantity) {
		t.Errorf("Expected %d elixirs, got %d", config.ItemsSpawnConfig.ElixirsQuantity, elixirCount)
	}
	if scrollCount != int(config.ItemsSpawnConfig.ScrollsQuantity) {
		t.Errorf("Expected %d scrolls, got %d", config.ItemsSpawnConfig.ScrollsQuantity, scrollCount)
	}
	if weaponCount != int(config.ItemsSpawnConfig.WeaponsQuantity) {
		t.Errorf("Expected %d weapons, got %d", config.ItemsSpawnConfig.WeaponsQuantity, weaponCount)
	}

	// Проверяем что все позиции уникальны
	positions := make(map[primitives.Point2D[int]]bool)
	for _, item := range entities.Items {
		pos := item.GetPosition()
		if positions[pos] {
			t.Errorf("Duplicate position found for item: %+v", pos)
		}
		positions[pos] = true
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_InRoomBounds(t *testing.T) {
	// Проверяем что все сущности находятся внутри комнат
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	spawner := NewRoomBasedEntitySpawner()
	entities, err := spawner.SpawnEntities(rooms, config.ItemsSpawnConfig, config.EnemiesSpawnConfig, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Функция для проверки что точка находится в какой-то комнате
	isInRoom := func(point primitives.Point2D[int]) bool {
		for _, room := range rooms {
			minX := room.Box.Point.X
			minY := room.Box.Point.Y
			maxX := room.Box.Point.X + int(room.Box.Size.Width)
			maxY := room.Box.Point.Y + int(room.Box.Size.Height)

			if point.X >= minX && point.X < maxX && point.Y >= minY && point.Y < maxY {
				return true
			}
		}
		return false
	}

	// Проверяем все сущности
	for i, item := range entities.Items {
		pos := item.GetPosition()
		if !isInRoom(pos) {
			t.Errorf("Item %d at %+v is not in any room", i, pos)
		}
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_NoStartRoom(t *testing.T) {
	// Проверяем что сущности не спавнятся в стартовой комнате
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	// Находим стартовую комнату
	var startRoom *Room
	for i, room := range rooms {
		if room.Type == RoomTypeStart {
			startRoom = &rooms[i]
			break
		}
	}
	if startRoom == nil {
		t.Fatal("No start room found")
	}

	spawner := NewRoomBasedEntitySpawner()
	entities, err := spawner.SpawnEntities(rooms, config.ItemsSpawnConfig, config.EnemiesSpawnConfig, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Функция для проверки что точка находится в стартовой комнате
	isInStartRoom := func(point primitives.Point2D[int]) bool {
		minX := startRoom.Box.Point.X
		minY := startRoom.Box.Point.Y
		maxX := startRoom.Box.Point.X + int(startRoom.Box.Size.Width)
		maxY := startRoom.Box.Point.Y + int(startRoom.Box.Size.Height)

		return point.X >= minX && point.X < maxX && point.Y >= minY && point.Y < maxY
	}

	// Проверяем что НИ ОДНА сущность не в стартовой комнате
	for i, item := range entities.Items {
		pos := item.GetPosition()
		if isInStartRoom(pos) {
			t.Errorf("Item %d at %+v is in start room (should not be)", i, pos)
		}
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_EmptyRooms(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	spawner := NewRoomBasedEntitySpawner()
	_, err := spawner.SpawnEntities([]Room{}, config.ItemsSpawnConfig, config.EnemiesSpawnConfig, source)

	if err == nil {
		t.Error("Expected error for empty rooms, got nil")
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_ZeroCounts(t *testing.T) {
	// Проверяем что можно создать spawner с нулевыми значениями
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)
	config.ItemsSpawnConfig.FoodsQuantity = 0
	config.ItemsSpawnConfig.ElixirsQuantity = 0
	config.ItemsSpawnConfig.ScrollsQuantity = 0
	config.ItemsSpawnConfig.WeaponsQuantity = 0

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	spawner := NewRoomBasedEntitySpawner()
	entities, err := spawner.SpawnEntities(rooms, config.ItemsSpawnConfig, config.EnemiesSpawnConfig, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	if len(entities.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(entities.Items))
	}
}
