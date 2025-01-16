package config

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
)

func (c *Config) SetUser(userName string) {
	updatedConfig := Config{DbUrl: c.DbUrl, CurrentUserName: userName}
	configBytes, err := mapConfigToJson(updatedConfig)
	if err != nil {
		log.Fatalf("config::SetUser: Failed to map config to JSON: %v", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("config::SetUser: Faile to get current user: %v", err)
	}

	filePath := getFilePath(currentUser)
	err = os.WriteFile(filePath, configBytes, 0644)
	if err != nil {
		log.Fatalf("config::SetUser: Ran into a issue updating config: %v", err)
	}

	c.CurrentUserName = userName
}

func mapConfigToJson(cfg Config) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "")
}
