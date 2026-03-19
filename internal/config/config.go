package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFile = ".gatorconfig.json"

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	dat, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(dat, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(user_name string) error {
	cfg.CurrentUserName = user_name
	return write(*cfg)
}

func write(cfg Config) error {
	dat, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, dat, 0644); err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	wd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(wd, configFile)
	return fullPath, nil
}
