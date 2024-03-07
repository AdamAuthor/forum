package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	Driver     string `json:"driver"`
	Path       string `json:"path"`
	FileName   string `json:"fileName"`
	SchemesDir string `json:"schemesDir"`
}

func ReadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
