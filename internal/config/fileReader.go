package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
)

const fileName = ".gatorconfig.json"

func Read() Config {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("config::Read: Failed to get current user: %v", err)
	}

	configFilePath := getFilePath(currentUser)
	fileBytes, err := readFile(configFilePath)
	if err != nil {
		log.Fatalf("config::Read: Failed to read file: %v", err)
	}

	config := Config{}
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalf("config::Read: Failed to map filecontent to type Config: %v", err)
	}

	return config
}

func getFilePath(u *user.User) string {
	return fmt.Sprintf("%v/%v", u.HomeDir, fileName)
}

func readFile(fp string) ([]byte, error) {
	file, err := os.Open(fp)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return fileBytes, nil
}
