package world

import (
	"gogue/internal/model/primitives"
	"math/rand"
	"testing"
)

func TestStartRoomPlayerSpawner_SpawnPlayer(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	// Генерируем комнаты
	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	// Создаём игрока
	spawner := NewStartRoomPlayerSpawner()
	player, err := spawner.SpawnPlayer(rooms, source)
	if err != nil {
		t.Fatalf("SpawnPlayer failed: %v", err)
	}

	// Проверяем что игрок создан
	if player == nil {
		t.Fatal("Player is nil")
	}

	// Проверяем что у игрока есть позиция
	if player.Box == nil {
		t.Fatal("Player shape is nil")
	}

	// Проверяем базовые характеристики игрока
	if player.Attributes.Health <= 0 {
		t.Errorf("Player health should be positive, got %f", player.Attributes.Health)
	}
	if player.CharacterLevel != 1 {
		t.Errorf("Player level should be 1, got %d", player.CharacterLevel)
	}
}

func TestStartRoomPlayerSpawner_SpawnPlayer_InStartRoom(t *testing.T) {
	// Проверяем что игрок спавнится в стартовой комнате
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

	spawner := NewStartRoomPlayerSpawner()
	player, err := spawner.SpawnPlayer(rooms, source)
	if err != nil {
		t.Fatalf("SpawnPlayer failed: %v", err)
	}

	// Проверяем что игрок в пределах стартовой комнаты
	playerPos := player.Box.Point
	minX := startRoom.Box.Point.X
	minY := startRoom.Box.Point.Y
	maxX := startRoom.Box.Point.X + int(startRoom.Box.Size.Width)
	maxY := startRoom.Box.Point.Y + int(startRoom.Box.Size.Height)

	if playerPos.X < minX || playerPos.X >= maxX || playerPos.Y < minY || playerPos.Y >= maxY {
		t.Errorf("Player position %+v is not in start room bounds (X: %d-%d, Y: %d-%d)",
			playerPos, minX, maxX, minY, maxY)
	}
}

func TestStartRoomPlayerSpawner_GetStartPosition(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	spawner := NewStartRoomPlayerSpawner()
	position, err := spawner.GetStartPosition(rooms, source)
	if err != nil {
		t.Fatalf("GetStartPosition failed: %v", err)
	}

	// Проверяем что позиция не нулевая
	if position.X == 0 && position.Y == 0 {
		t.Error("Position should not be zero")
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

	// Проверяем что позиция в стартовой комнате
	minX := startRoom.Box.Point.X
	minY := startRoom.Box.Point.Y
	maxX := startRoom.Box.Point.X + int(startRoom.Box.Size.Width)
	maxY := startRoom.Box.Point.Y + int(startRoom.Box.Size.Height)

	if position.X < minX || position.X >= maxX || position.Y < minY || position.Y >= maxY {
		t.Errorf("Position %+v is not in start room bounds (X: %d-%d, Y: %d-%d)",
			position, minX, maxX, minY, maxY)
	}
}

func TestStartRoomPlayerSpawner_SpawnPlayer_EmptyRooms(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	spawner := NewStartRoomPlayerSpawner()

	_, err := spawner.SpawnPlayer([]Room{}, source)
	if err == nil {
		t.Error("Expected error for empty rooms, got nil")
	}
}

func TestStartRoomPlayerSpawner_SpawnPlayer_NoStartRoom(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	spawner := NewStartRoomPlayerSpawner()

	// Создаём комнаты без стартовой
	rooms := []Room{
		{
			Type: RoomTypeOrdinary,
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 10, Y: 10},
				Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
			},
		},
		{
			Type: RoomTypeFinish,
			Box: &primitives.Box{
				Point: primitives.Point2D[int]{X: 20, Y: 20},
				Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
			},
		},
	}

	_, err := spawner.SpawnPlayer(rooms, source)
	if err == nil {
		t.Error("Expected error when no start room exists, got nil")
	}
}

func TestStartRoomPlayerSpawner_Consistency(t *testing.T) {
	// Проверяем что SpawnPlayer и GetStartPosition возвращают согласованные позиции
	source1 := rand.New(rand.NewSource(randomSeedTest))
	source2 := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	roomGen := NewGridRoomGenerator()
	rooms1, _, err := roomGen.GenerateRooms(config, source1)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	spawner := NewStartRoomPlayerSpawner()
	player, err := spawner.SpawnPlayer(rooms1, source1)
	if err != nil {
		t.Fatalf("SpawnPlayer failed: %v", err)
	}

	rooms2, _, err := roomGen.GenerateRooms(config, source2)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	position, err := spawner.GetStartPosition(rooms2, source2)
	if err != nil {
		t.Fatalf("GetStartPosition failed: %v", err)
	}

	// С одинаковым seed должны получить одинаковые позиции
	if player.Box.Point != *position {
		t.Errorf("SpawnPlayer and GetStartPosition returned different positions: %+v vs %+v",
			player.Box.Point, *position)
	}
}
