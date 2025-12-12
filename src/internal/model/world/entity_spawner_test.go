package world

import (
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
	entities, err := spawner.SpawnEntities(rooms, config.ItemCounts, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Проверяем что сущности созданы
	if len(entities.Foods) != int(config.ItemCounts.FoodsQuntity) {
		t.Errorf("Expected %d foods, got %d", config.ItemCounts.FoodsQuntity, len(entities.Foods))
	}
	if len(entities.Elixirs) != int(config.ItemCounts.ElixirsQuntity) {
		t.Errorf("Expected %d elixirs, got %d", config.ItemCounts.ElixirsQuntity, len(entities.Elixirs))
	}
	if len(entities.Scrolls) != int(config.ItemCounts.ScrollsQuntity) {
		t.Errorf("Expected %d scrolls, got %d", config.ItemCounts.ScrollsQuntity, len(entities.Scrolls))
	}
	if len(entities.Weapons) != int(config.ItemCounts.WeaponsQuntity) {
		t.Errorf("Expected %d weapons, got %d", config.ItemCounts.WeaponsQuntity, len(entities.Weapons))
	}

	// Проверяем что все позиции уникальны
	positions := make(map[primitives.Point2D[int]]bool)
	for _, food := range entities.Foods {
		if positions[food.Shape.Point] {
			t.Errorf("Duplicate position found for food: %+v", food.Shape.Point)
		}
		positions[food.Shape.Point] = true
	}
	for _, elixir := range entities.Elixirs {
		if positions[elixir.Shape.Point] {
			t.Errorf("Duplicate position found for elixir: %+v", elixir.Shape.Point)
		}
		positions[elixir.Shape.Point] = true
	}
	for _, scroll := range entities.Scrolls {
		if positions[scroll.Shape.Point] {
			t.Errorf("Duplicate position found for scroll: %+v", scroll.Shape.Point)
		}
		positions[scroll.Shape.Point] = true
	}
	for _, weapon := range entities.Weapons {
		if positions[weapon.Shape.Point] {
			t.Errorf("Duplicate position found for weapon: %+v", weapon.Shape.Point)
		}
		positions[weapon.Shape.Point] = true
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
	entities, err := spawner.SpawnEntities(rooms, config.ItemCounts, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Функция для проверки что точка находится в какой-то комнате
	isInRoom := func(point primitives.Point2D[int]) bool {
		for _, room := range rooms {
			minX := room.Shape.Point.X
			minY := room.Shape.Point.Y
			maxX := room.Shape.Point.X + int(room.Shape.Size.Width)
			maxY := room.Shape.Point.Y + int(room.Shape.Size.Height)

			if point.X >= minX && point.X < maxX && point.Y >= minY && point.Y < maxY {
				return true
			}
		}
		return false
	}

	// Проверяем все сущности
	for i, food := range entities.Foods {
		if !isInRoom(food.Shape.Point) {
			t.Errorf("Food %d at %+v is not in any room", i, food.Shape.Point)
		}
	}
	for i, elixir := range entities.Elixirs {
		if !isInRoom(elixir.Shape.Point) {
			t.Errorf("Elixir %d at %+v is not in any room", i, elixir.Shape.Point)
		}
	}
	for i, scroll := range entities.Scrolls {
		if !isInRoom(scroll.Shape.Point) {
			t.Errorf("Scroll %d at %+v is not in any room", i, scroll.Shape.Point)
		}
	}
	for i, weapon := range entities.Weapons {
		if !isInRoom(weapon.Shape.Point) {
			t.Errorf("Weapon %d at %+v is not in any room", i, weapon.Shape.Point)
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
	entities, err := spawner.SpawnEntities(rooms, config.ItemCounts, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	// Функция для проверки что точка находится в стартовой комнате
	isInStartRoom := func(point primitives.Point2D[int]) bool {
		minX := startRoom.Shape.Point.X
		minY := startRoom.Shape.Point.Y
		maxX := startRoom.Shape.Point.X + int(startRoom.Shape.Size.Width)
		maxY := startRoom.Shape.Point.Y + int(startRoom.Shape.Size.Height)

		return point.X >= minX && point.X < maxX && point.Y >= minY && point.Y < maxY
	}

	// Проверяем что НИ ОДНА сущность не в стартовой комнате
	for i, food := range entities.Foods {
		if isInStartRoom(food.Shape.Point) {
			t.Errorf("Food %d at %+v is in start room (should not be)", i, food.Shape.Point)
		}
	}
	for i, elixir := range entities.Elixirs {
		if isInStartRoom(elixir.Shape.Point) {
			t.Errorf("Elixir %d at %+v is in start room (should not be)", i, elixir.Shape.Point)
		}
	}
	for i, scroll := range entities.Scrolls {
		if isInStartRoom(scroll.Shape.Point) {
			t.Errorf("Scroll %d at %+v is in start room (should not be)", i, scroll.Shape.Point)
		}
	}
	for i, weapon := range entities.Weapons {
		if isInStartRoom(weapon.Shape.Point) {
			t.Errorf("Weapon %d at %+v is in start room (should not be)", i, weapon.Shape.Point)
		}
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_EmptyRooms(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	spawner := NewRoomBasedEntitySpawner()
	_, err := spawner.SpawnEntities([]Room{}, config.ItemCounts, source)

	if err == nil {
		t.Error("Expected error for empty rooms, got nil")
	}
}

func TestRoomBasedEntitySpawner_SpawnEntities_ZeroCounts(t *testing.T) {
	// Проверяем что можно создать spawner с нулевыми значениями
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)
	config.ItemCounts.FoodsQuntity = 0
	config.ItemCounts.ElixirsQuntity = 0
	config.ItemCounts.ScrollsQuntity = 0
	config.ItemCounts.WeaponsQuntity = 0

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	spawner := NewRoomBasedEntitySpawner()
	entities, err := spawner.SpawnEntities(rooms, config.ItemCounts, source)
	if err != nil {
		t.Fatalf("SpawnEntities failed: %v", err)
	}

	if len(entities.Foods) != 0 {
		t.Errorf("Expected 0 foods, got %d", len(entities.Foods))
	}
	if len(entities.Elixirs) != 0 {
		t.Errorf("Expected 0 elixirs, got %d", len(entities.Elixirs))
	}
	if len(entities.Scrolls) != 0 {
		t.Errorf("Expected 0 scrolls, got %d", len(entities.Scrolls))
	}
	if len(entities.Weapons) != 0 {
		t.Errorf("Expected 0 weapons, got %d", len(entities.Weapons))
	}
}
