package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	configPath, err := getConfigPath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	data, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	configPath, err := getConfigPath()

	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0666)

	if err != nil {
		return err
	}

	return nil
}

func getConfigPath() (string, error) {
	homeDir, error := os.UserHomeDir()

	if error != nil {
		return "", error
	}

	return homeDir + "/" + configFileName, nil
}
