package world

import (
	"math"

	"gogue/internal/model/primitives"
)

type LevelConfig struct {
	MapSize primitives.Size2D[uint]

	// Параметры комнат
	RoomMinWidth   int
	RoomMinHeight  int
	MapMaxWidth    uint
	MapMaxHeight   uint
	MinRoomPadding int

	// Параметры генерации мира
	MaxExtraPassageCount int
	NumberXYSections     int
	RoomsCount           int
	ItemsSpawnConfig     ItemSpawnConfig
	EnemiesSpawnConfig   EnemySpawnConfig

	// Конфигурация сложности
	DifficultyScaling DifficultyScaling
}

type ItemSpawnConfig struct {
	FoodsQuantity   uint
	ElixirsQuantity uint
	ScrollsQuantity uint
	WeaponsQuantity uint
}

type EnemySpawnConfig struct {
	Quantity uint

	// Множитель атрибутов врагов (1.0 = базовые значения).
	// Увеличивается с каждым уровнем подземелья.
	AttributeMultiplier float64
}

// TreasureDropConfig определяет проценты выпадения типов сокровищ с побежденных врагов.
// Сумма всех полей должна быть равна 100.
type TreasureDropConfig struct {
	GoldPercent     uint
	GemPercent      uint
	ArtifactPercent uint
	MysteryPercent  uint
}

// DifficultyScaling содержит параметры, управляющие тем, как сложность
// масштабируется с ростом номера уровня подземелья.
type DifficultyScaling struct {
	// Враги: базовое количество + прирост за уровень
	BaseEnemyCount       uint
	EnemyCountPerLevel   uint
	MaxEnemyCount        uint
	EnemyAttributeGrowth float64 // Прирост множителя атрибутов за уровень (например, 0.15 => +15% за уровень)

	// Предметы: базовое количество - убыль за уровень (с минимальным порогом)
	BaseFoods     uint
	BaseElixirs   uint
	BaseScrolls   uint
	BaseWeapons   uint
	ItemDecayRate float64 // Доля убыли предметов за уровень (0.08 => -8% от базы за уровень)
	MinFoods      uint
	MinElixirs    uint
	MinScrolls    uint
	MinWeapons    uint

	// Сокровища: сдвиг процентов от Gold к более ценным типам
	BaseTreasureDrop       TreasureDropConfig
	GemGrowthPerLevel      float64 // Прирост % гемов за уровень
	ArtifactGrowthPerLevel float64 // Прирост % артефактов за уровень
}

func DefaultDifficultyScaling() DifficultyScaling {
	return DifficultyScaling{
		// Враги
		BaseEnemyCount:       10,
		EnemyCountPerLevel:   2,
		MaxEnemyCount:        30,
		EnemyAttributeGrowth: 0.12,

		// Предметы
		BaseFoods:     7,
		BaseElixirs:   5,
		BaseScrolls:   3,
		BaseWeapons:   1,
		ItemDecayRate: 0.08,
		MinFoods:      2,
		MinElixirs:    1,
		MinScrolls:    1,
		MinWeapons:    0,

		// Сокровища
		BaseTreasureDrop: TreasureDropConfig{
			GoldPercent:     65,
			GemPercent:      20,
			ArtifactPercent: 5,
			MysteryPercent:  10,
		},
		GemGrowthPerLevel:      1.5,
		ArtifactGrowthPerLevel: 0.8,
	}
}

func DefaultLevelConfig(mapSize primitives.Size2D[uint]) LevelConfig {
	scaling := DefaultDifficultyScaling()
	return LevelConfig{
		MapSize:              mapSize,
		RoomMinWidth:         3,
		RoomMinHeight:        3,
		MapMaxWidth:          200,
		MapMaxHeight:         150,
		MinRoomPadding:       1,
		MaxExtraPassageCount: 2,
		NumberXYSections:     3,
		RoomsCount:           9,
		ItemsSpawnConfig: ItemSpawnConfig{
			FoodsQuantity:   scaling.BaseFoods,
			ElixirsQuantity: scaling.BaseElixirs,
			ScrollsQuantity: scaling.BaseScrolls,
			WeaponsQuantity: scaling.BaseWeapons,
		},
		EnemiesSpawnConfig: EnemySpawnConfig{
			Quantity:            scaling.BaseEnemyCount,
			AttributeMultiplier: 1.0,
		},
		DifficultyScaling: scaling,
	}
}

// ScaleForLevel пересчитывает конфигурацию спавна врагов и предметов
// в соответствии с номером уровня подземелья.
// levelNumber начинается с 1 (первый уровень — базовая сложность).
func (c *LevelConfig) ScaleForLevel(levelNumber uint) {
	if levelNumber == 0 {
		levelNumber = 1
	}

	s := c.DifficultyScaling
	progression := float64(levelNumber - 1) // 0 для первого уровня

	// --- Враги ---
	enemyCount := s.BaseEnemyCount + s.EnemyCountPerLevel*uint(progression)
	if enemyCount > s.MaxEnemyCount {
		enemyCount = s.MaxEnemyCount
	}
	c.EnemiesSpawnConfig.Quantity = enemyCount
	c.EnemiesSpawnConfig.AttributeMultiplier = 1.0 + s.EnemyAttributeGrowth*progression

	// --- Предметы (экспоненциальное затухание, чтобы не обнулялось слишком резко) ---
	decayMultiplier := math.Pow(1.0-s.ItemDecayRate, progression)

	c.ItemsSpawnConfig.FoodsQuantity = clampUint(
		uint(math.Round(float64(s.BaseFoods)*decayMultiplier)),
		s.MinFoods, s.BaseFoods,
	)
	c.ItemsSpawnConfig.ElixirsQuantity = clampUint(
		uint(math.Round(float64(s.BaseElixirs)*decayMultiplier)),
		s.MinElixirs, s.BaseElixirs,
	)
	c.ItemsSpawnConfig.ScrollsQuantity = clampUint(
		uint(math.Round(float64(s.BaseScrolls)*decayMultiplier)),
		s.MinScrolls, s.BaseScrolls,
	)
	c.ItemsSpawnConfig.WeaponsQuantity = clampUint(
		uint(math.Round(float64(s.BaseWeapons)*decayMultiplier)),
		s.MinWeapons, s.BaseWeapons,
	)
}

// GetTreasureDropConfig возвращает конфигурацию выпадения сокровищ
// для указанного номера уровня. С ростом уровня Gold-процент уменьшается
// в пользу Gem и Artifact.
func (c *LevelConfig) GetTreasureDropConfig(levelNumber uint) TreasureDropConfig {
	if levelNumber == 0 {
		levelNumber = 1
	}

	s := c.DifficultyScaling
	progression := float64(levelNumber - 1)
	base := s.BaseTreasureDrop

	gemGrowth := s.GemGrowthPerLevel * progression
	artifactGrowth := s.ArtifactGrowthPerLevel * progression

	// Gem и Artifact растут, Gold — уменьшается, Mystery — фиксирован
	gem := float64(base.GemPercent) + gemGrowth
	artifact := float64(base.ArtifactPercent) + artifactGrowth
	mystery := float64(base.MysteryPercent)

	// Gold забирает остаток, но не может быть отрицательным
	gold := 100.0 - gem - artifact - mystery
	if gold < 5 {
		// Если Gold слишком мал, перераспределяем пропорционально Gem и Artifact
		excess := 5 - gold
		gold = 5
		total := gem + artifact
		if total > 0 {
			gem -= excess * (gem / total)
			artifact -= excess * (artifact / total)
		}
	}

	return TreasureDropConfig{
		GoldPercent:     uint(math.Round(gold)),
		GemPercent:      uint(math.Round(gem)),
		ArtifactPercent: uint(math.Round(artifact)),
		MysteryPercent:  uint(math.Round(mystery)),
	}
}

// clampUint ограничивает значение v диапазоном [minVal, maxVal].
func clampUint(v, minVal, maxVal uint) uint {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
