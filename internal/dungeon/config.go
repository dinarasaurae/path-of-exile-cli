package dungeon

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Config struct {
	Floors   int
	Monsters int
	OpenAt   string
	Duration int
}

type Settings struct {
	Floors        int
	Monsters      int
	OpenAt        time.Duration
	CloseAt       time.Duration
	OrdinaryCount int
}

func LoadSettings(reader io.Reader) (Settings, error) {
	var config Config
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return Settings{}, err
	}

	return NewSettings(config)
}

func NewSettings(config Config) (Settings, error) {
	if config.Floors < 1 {
		return Settings{}, fmt.Errorf("floors must be positive")
	}

	if config.Monsters < 0 {
		return Settings{}, fmt.Errorf("monsters cannot be negative")
	}

	if config.Duration <= 0 {
		return Settings{}, fmt.Errorf("duration must be positive")
	}

	return Settings{
		Floors:        config.Floors,
		Monsters:      config.Monsters,
		OrdinaryCount: config.Floors - 1,
	}, nil
}
