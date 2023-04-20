package settings

import (
	"os"
	"path/filepath"
)

const CONFIG_JSON_PATH = ".kpm/config/config.json"

type Settings struct {
	CredentialsFile string
}

func Init() (*Settings, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Settings{
		CredentialsFile: filepath.Join(homeDir, CONFIG_JSON_PATH),
	}, nil
}
