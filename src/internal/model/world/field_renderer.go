package world

import (
	"gogue/internal/common"
	"gogue/internal/model/entities"
	"gogue/internal/model/items"
	"gogue/internal/model/primitives"
	"gogue/internal/utils"
)

type FieldRenderer interface {
	RenderField(width, height int, level *Level) [][]common.GameEntityType
}

type DefaultFieldRenderer struct{}

func NewDefaultFieldRenderer() *DefaultFieldRenderer {
	return &DefaultFieldRenderer{}
}

func (r *DefaultFieldRenderer) RenderField(width, height int, level *Level) [][]common.GameEntityType {
	field := utils.CreateEmpty2DSlice[common.GameEntityType](height, width)

	// Порядок рендеринга (от фона к переднему плану):
	// 1. Комнаты (стены и пол)
	// 2. Коридоры и двери
	// 3. Предметы (еда, зелья, свитки, оружие)
	// 4. Враги
	// 5. Игрок (всегда поверх всех)

	r.renderRooms(level.Rooms, level.FinishPortal, field)
	r.renderPassages(level.Passages, field)
	r.renderItems(level.Items, field, width, height)
	r.renderEnemies(level.Enemies, field, width, height)
	r.renderPlayer(level.Player, field, width, height)

	return field
}

func (r *DefaultFieldRenderer) renderRooms(rooms []Room, finishPortal primitives.Box, field [][]common.GameEntityType) {
	if len(rooms) == 0 {
		return
	}

	for _, room := range rooms {
		r.renderSingleRoom(room, field)
	}

	// Рендерим портал поверх комнаты
	r.renderPortal(finishPortal, field)
}

func (r *DefaultFieldRenderer) renderSingleRoom(room Room, field [][]common.GameEntityType) {
	width := int(room.Box.Size.Width)
	height := int(room.Box.Size.Height)

	startX := room.Box.Point.X
	endX := startX + width - 1
	startY := room.Box.Point.Y
	endY := startY + height - 1

	// Горизонтальные стены (верх и низ)
	for col := startX; col <= endX; col++ {
		if r.isInBounds(col, startY, field) {
			field[startY][col] = common.WorldTypeWall
		}
		if r.isInBounds(col, endY, field) {
			field[endY][col] = common.WorldTypeWall
		}
	}

	// Вертикальные стены (левая и правая)
	for row := startY; row <= endY; row++ {
		if r.isInBounds(startX, row, field) {
			field[row][startX] = common.WorldTypeWall
		}
		if r.isInBounds(endX, row, field) {
			field[row][endX] = common.WorldTypeWall
		}
	}

	// Пространство внутри комнаты
	for row := startY + 1; row <= endY-1; row++ {
		for col := startX + 1; col <= endX-1; col++ {
			if r.isInBounds(col, row, field) {
				field[row][col] = common.WorldTypeRoomFloor
			}
		}
	}
}

func (r *DefaultFieldRenderer) renderPortal(portal primitives.Box, field [][]common.GameEntityType) {
	x, y := portal.Point.X, portal.Point.Y
	if r.isInBounds(x, y, field) {
		field[y][x] = common.WorldTypePortal
	}
}

func (r *DefaultFieldRenderer) renderPassages(passages []Passage, field [][]common.GameEntityType) {
	for _, passage := range passages {
		r.renderSinglePassage(passage, field)
	}
}

func (r *DefaultFieldRenderer) renderSinglePassage(passage Passage, field [][]common.GameEntityType) {
	// Отрисовка пути коридора
	for _, point := range passage.Way {
		if r.isInBounds(point.X, point.Y, field) {
			field[point.Y][point.X] = common.WorldTypePassage
		}
	}

	// Отрисовка дверей
	if r.isInBounds(passage.DoorOne.X, passage.DoorOne.Y, field) {
		field[passage.DoorOne.Y][passage.DoorOne.X] = common.WorldTypeDoor
	}
	if r.isInBounds(passage.DoorTwo.X, passage.DoorTwo.Y, field) {
		field[passage.DoorTwo.Y][passage.DoorTwo.X] = common.WorldTypeDoor
	}
}

func (r *DefaultFieldRenderer) renderItems(
	itemsList []primitives.Positional2D[int],
	field [][]common.GameEntityType,
	width, height int,
) {
	for _, item := range itemsList {
		pt := item.GetPosition()
		if !r.isInBoundsWH(pt.X, pt.Y, width, height) {
			continue
		}

		// Type switch для определения типа предмета
		switch item.(type) {
		case *items.Food:
			field[pt.Y][pt.X] = common.Food
		case *items.Elixir:
			field[pt.Y][pt.X] = common.Elixir
		case *items.Scroll:
			field[pt.Y][pt.X] = common.Scroll
		case *items.Weapon:
			field[pt.Y][pt.X] = common.Weapon
		case *items.Treasure:
			// @todo: Добавить константу Treasure в common.GameEntityType
		}
	}
}

func (r *DefaultFieldRenderer) renderEnemies(enemies []primitives.Positional2D[int], field [][]common.GameEntityType, width, height int) {
	for _, enemy := range enemies {
		// Невидимый Ghost не отображается на карте
		if ghost, ok := enemy.(*entities.Ghost); ok && !ghost.IsVisible {
			continue
		}

		// Замаскированный Mimic отображается как предмет (оружие), а не как враг
		if mimic, ok := enemy.(*entities.Mimic); ok && mimic.IsDisguised {
			pt := enemy.GetPosition()
			if r.isInBoundsWH(pt.X, pt.Y, width, height) {
				field[pt.Y][pt.X] = common.Weapon
			}
			continue
		}

		pt := enemy.GetPosition()
		if r.isInBoundsWH(pt.X, pt.Y, width, height) {
			field[pt.Y][pt.X] = convertEnemyToEntityType(enemy)
		}
	}
}

func (r *DefaultFieldRenderer) renderPlayer(player *entities.Player, field [][]common.GameEntityType, width, height int) {
	if player == nil {
		return
	}

	px := player.Character.Box.Point.X
	py := player.Character.Box.Point.Y

	if r.isInBoundsWH(px, py, width, height) {
		field[py][px] = common.EntityTypePlayer
	}
}

func (r *DefaultFieldRenderer) isInBounds(x, y int, field [][]common.GameEntityType) bool {
	if len(field) == 0 {
		return false
	}
	return y >= 0 && y < len(field) && x >= 0 && x < len(field[0])
}

func (r *DefaultFieldRenderer) isInBoundsWH(x, y, width, height int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
}

// Пример того, как можем из интерфейса определить конкретный тип, без доп. полей
func convertEnemyToEntityType(enemy any) common.GameEntityType {
	switch enemy.(type) {
	case *entities.Zombie:
		return common.EntityTypeZombie
	case *entities.Vampire:
		return common.EntityTypeVampire
	case *entities.Ghost:
		return common.EntityTypeGhost
	case *entities.Ogre:
		return common.EntityTypeOgre
	case *entities.SnakeMage:
		return common.EntityTypeSnakeMage
	case *entities.Mimic:
		return common.EntityTypeMimic
	}

	return common.EntityTypeNone
}
