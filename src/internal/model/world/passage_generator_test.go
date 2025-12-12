package world

import (
	"gogue/internal/model/primitives"
	"math/rand"
	"testing"
)

func TestConnectionTreePassageGenerator_GeneratePassages(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	// Сначала генерируем комнаты
	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	// Генерируем коридоры
	passageGen := NewConnectionTreePassageGenerator()
	passages, err := passageGen.GeneratePassages(rooms, config, source)
	if err != nil {
		t.Fatalf("GeneratePassages failed: %v", err)
	}

	// Проверяем что коридоры созданы
	if len(passages) == 0 {
		t.Error("Expected at least one passage, got none")
	}

	// Проверяем что количество коридоров разумно (для 9 комнат минимум 8 для связности + extra edges)
	minPassages := config.RoomsCount - 1
	maxPassages := minPassages + config.MaxExtraPassageCount
	if len(passages) < minPassages || len(passages) > maxPassages+1 { // +1 для погрешности
		t.Errorf("Expected %d-%d passages, got %d", minPassages, maxPassages, len(passages))
	}

	// Проверяем структуру каждого коридора
	for i, passage := range passages {
		// Проверяем что двери определены
		if passage.DoorOne.X == 0 && passage.DoorOne.Y == 0 {
			t.Errorf("Passage %d: DoorOne is not set", i)
		}
		if passage.DoorTwo.X == 0 && passage.DoorTwo.Y == 0 {
			t.Errorf("Passage %d: DoorTwo is not set", i)
		}

		// Проверяем что путь существует
		if len(passage.Way) == 0 {
			t.Errorf("Passage %d: Way is empty", i)
		}

		// Проверяем что путь находится в пределах карты
		for j, point := range passage.Way {
			if point.X < 0 || point.Y < 0 {
				t.Errorf("Passage %d, point %d: negative coordinates %+v", i, j, point)
			}
			if point.X >= int(mapSize.Width) || point.Y >= int(mapSize.Height) {
				t.Errorf("Passage %d, point %d: exceeds map boundaries %+v", i, j, point)
			}
		}
	}
}

func TestConnectionTreePassageGenerator_GeneratePassages_EmptyRooms(t *testing.T) {
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	passageGen := NewConnectionTreePassageGenerator()
	_, err := passageGen.GeneratePassages([]Room{}, config, source)

	if err == nil {
		t.Error("Expected error for empty rooms, got nil")
	}
}

func TestConnectionTreePassageGenerator_GeneratePassages_Connectivity(t *testing.T) {
	// Проверяем что коридоры действительно соединяют комнаты
	source := rand.New(rand.NewSource(randomSeedTest))
	mapSize := primitives.Size2D[uint]{Height: 30, Width: 90}
	config := DefaultLevelConfig(mapSize)

	roomGen := NewGridRoomGenerator()
	rooms, _, err := roomGen.GenerateRooms(config, source)
	if err != nil {
		t.Fatalf("Failed to generate rooms: %v", err)
	}

	passageGen := NewConnectionTreePassageGenerator()
	passages, err := passageGen.GeneratePassages(rooms, config, source)
	if err != nil {
		t.Fatalf("GeneratePassages failed: %v", err)
	}

	// Проверяем что каждая дверь принадлежит какой-то комнате
	for i, passage := range passages {
		doorOneInRoom := false
		doorTwoInRoom := false

		for _, room := range rooms {
			if isPointInOrOnBorder(passage.DoorOne, room.Shape) {
				doorOneInRoom = true
			}
			if isPointInOrOnBorder(passage.DoorTwo, room.Shape) {
				doorTwoInRoom = true
			}
		}

		if !doorOneInRoom {
			t.Errorf("Passage %d: DoorOne %+v is not in any room", i, passage.DoorOne)
		}
		if !doorTwoInRoom {
			t.Errorf("Passage %d: DoorTwo %+v is not in any room", i, passage.DoorTwo)
		}
	}
}

// isPointInOrOnBorder проверяет находится ли точка внутри или на границе прямоугольника
func isPointInOrOnBorder(point primitives.Point2D[int], box primitives.Box) bool {
	minX := box.Point.X
	minY := box.Point.Y
	maxX := box.Point.X + int(box.Size.Width) - 1
	maxY := box.Point.Y + int(box.Size.Height) - 1

	return point.X >= minX && point.X <= maxX && point.Y >= minY && point.Y <= maxY
}

func TestConnectionTreePassageGenerator_DoorPlacer(t *testing.T) {
	// Тестируем DoorPlacer отдельно
	placer := NewDoorPlacer()

	room1 := Room{
		Shape: primitives.Box{
			Point: primitives.Point2D[int]{X: 10, Y: 10},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
		},
	}

	room2 := Room{
		Shape: primitives.Box{
			Point: primitives.Point2D[int]{X: 20, Y: 10},
			Size:  primitives.Size2D[uint]{Width: 5, Height: 5},
		},
	}

	// Комнаты расположены горизонтально
	doorOne := primitives.Point2D[int]{X: room1.Shape.Point.X + int(room1.Shape.Size.Width) - 1, Y: 12}
	doorTwo := primitives.Point2D[int]{X: room2.Shape.Point.X, Y: 12}

	placer.PlaceDoors(&room1, &room2, doorOne, doorTwo)

	// Проверяем что двери добавлены
	if len(room1.Doors) == 0 {
		t.Error("room1 should have doors added")
	}
	if len(room2.Doors) == 0 {
		t.Error("room2 should have doors added")
	}

	// Проверяем что двери установлены правильно
	if room1.Doors[0] != doorOne {
		t.Errorf("Expected doorOne %+v in room1, got %+v", doorOne, room1.Doors[0])
	}
	if room2.Doors[0] != doorTwo {
		t.Errorf("Expected doorTwo %+v in room2, got %+v", doorTwo, room2.Doors[0])
	}
}
