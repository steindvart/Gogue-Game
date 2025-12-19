package world

import (
	"fmt"
	"gogue/internal/model/entities"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type StartRoomPlayerSpawner struct{}

func NewStartRoomPlayerSpawner() *StartRoomPlayerSpawner {
	return &StartRoomPlayerSpawner{}
}

func (s *StartRoomPlayerSpawner) SpawnPlayer(
	rooms []Room,
	random utils.Randomizer,
) (*entities.Player, error) {
	startPos, err := s.GetStartPosition(rooms, random)
	if err != nil {
		return nil, err
	}

	player := entities.NewPlayer(primitives.Box{
		Point: *startPos,
		Size:  primitives.Size2D[uint]{Height: 1, Width: 1},
	})

	return player, nil
}

func (s *StartRoomPlayerSpawner) GetStartPosition(
	rooms []Room,
	random utils.Randomizer,
) (*primitives.Point2D[int], error) {
	for _, room := range rooms {
		if room.Type == RoomTypeStart {
			pos, err := room.GetRandomFreePosition(random)
			if err != nil {
				return nil, fmt.Errorf("no free position in starting room: %w", err)
			}
			return pos, nil
		}
	}

	return nil, fmt.Errorf("starting room was not found")
}
