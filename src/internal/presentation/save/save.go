package save

import (
	"encoding/json"
	"gogue/internal/model/primitives"
	"gogue/internal/model/world"
	"gogue/internal/presentation/dto"
	"gogue/internal/utils"
	"os"
)

func SaveGame(level *world.Level, size primitives.Size2D[uint], filename string) error {
	dto := ConvertLevelToDto(level, size)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(dto)
}

func LoadGame(filename string, random utils.Randomizer) (*world.Level, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dto dto.GameSaveDto
	if err := json.NewDecoder(file).Decode(&dto); err != nil {
		return nil, err
	}

	return RestoreLevelFromDto(dto, random)
}
