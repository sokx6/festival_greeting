package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Email struct {
	From     string `toml:"from"`
	Port     int    `toml:"port"`
	Host     string `toml:"host"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type SendTime struct {
	Hour   int `toml:"hour"`
	Minute int `toml:"minute"`
	Second int `toml:"second"`
}

type Friend struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type Model struct {
	BaseUrl   string `toml:"base_url"`
	ModelName string `toml:"model_name"`
	ApiKey    string `toml:"api_key"`
}

type Config struct {
	Email      Email    `toml:"email"`
	SendTime   SendTime `toml:"send_time"`
	Friends    []Friend `toml:"friends"`
	AvatarURL  string   `toml:"avatar_url"`
	SenderName string   `toml:"sender_name"`
	Model      Model    `toml:"model"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &config, nil
}
