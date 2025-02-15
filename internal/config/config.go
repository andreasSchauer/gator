package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbURL           	string `json:"db_url"`
	CurrentUserName 	string `json:"current_user_name"`
}


func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}


func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("file doesn't exist: %v", err)
	}

	var cfg = Config{}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("problem decoding json: %v", err)
	}

	return cfg, nil
}


func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}


func write (cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("problem encoding json: %v", err)
	}

	err = os.WriteFile(configFilePath, data, 0600)
	if err != nil {
		return fmt.Errorf("problem writing to config file: %v", err)
	}

	return nil
}