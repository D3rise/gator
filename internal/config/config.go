package config

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
	CurrentUserId   uuid.UUID
	configPath      string
}

const (
	configFileName = ".gatorconfig.json"
)

func GetConfigPath() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(homeDirectory, configFileName)
}

func NewConfig(customConfigPath string) (Config, error) {
	var configPath string
	if customConfigPath == "" {
		configPath = GetConfigPath()
	} else {
		configPath = customConfigPath
	}

	return read(configPath)
}

func (c *Config) SetCurrentUserId(userId uuid.UUID) error {
	c.CurrentUserId = userId
	return nil
}

func (c *Config) SetCurrentUserName(name string) error {
	c.CurrentUserName = name
	return write(c.configPath, *c)
}

func read(configPath string) (Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	conf := Config{configPath: configPath}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

func write(configPath string, conf Config) error {
	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
