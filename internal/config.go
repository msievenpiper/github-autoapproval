package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Branch string   `json:"branch"`
	Repos  []string `json:"repos"`
	Probe  bool     `json:"probe"`
	Merge  bool     `json:"merge"`
}

func LoadConfig(path string) (Config, error) {
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}
		path = filepath.Join(cwd, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
