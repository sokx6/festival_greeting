package utils

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Email   Email    `toml:"email"`
	Friends []Friend `toml:"friends"`
}

type Email struct {
	Port     int    `toml:"port"`
	Addr     string `toml:"addr"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Friend struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &config, nil
}
