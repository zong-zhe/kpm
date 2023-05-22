package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"kusionstack.io/kpm/pkg/errors"
)

// The config.json used to persist user information
const CONFIG_JSON_PATH = ".kpm/config/config.json"

// The kpm.json used to describe the default configuration of kpm.
const KPM_JSON_PATH = ".kpm/config/kpm.json"

type KpmConf struct {
	DefaultOciRegistry string
	DefaultOciRepo     string
}

func DefaultKpmConf() KpmConf {
	return KpmConf{
		DefaultOciRegistry: "ghcr.io",
		DefaultOciRepo:     "KusionStack",
	}
}

type Settings struct {
	CredentialsFile string
	Conf            KpmConf
}

func (settings *Settings) DefauleOciRegistry() string {
	return settings.Conf.DefaultOciRegistry
}

func (settings *Settings) DefauleOciRepo() string {
	return settings.Conf.DefaultOciRepo
}

// GetFullJsonPath returns the full path of config.json and kpm.json file path under '$HOME/.kpm/config/'
func GetFullJsonPath(jsonFileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.InternalBug
	}

	return filepath.Join(home, jsonFileName), nil
}

// Init returns default kpm settings.
func Init() (*Settings, error) {
	credentialsFile, err := GetFullJsonPath(CONFIG_JSON_PATH)
	if err != nil {
		return nil, err
	}

	conf, err := loadOrCreateDefaultKpmJson()
	if err != nil {
		return nil, err
	}

	return &Settings{
		CredentialsFile: credentialsFile,
		Conf:            *conf,
	}, nil
}

func loadOrCreateDefaultKpmJson() (*KpmConf, error) {
	kpmConfpath, err := GetFullJsonPath(KPM_JSON_PATH)
	if err != nil {
		return nil, err
	}

	defaultKpmConf := DefaultKpmConf()

	b, err := ioutil.ReadFile(kpmConfpath)
	if os.IsNotExist(err) {
		b, err := json.Marshal(defaultKpmConf)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(kpmConfpath, b, 0644)
		if err != nil {
			return nil, err
		}
		return &defaultKpmConf, nil
	} else {
		err = json.Unmarshal(b, &defaultKpmConf)
		if err != nil {
			return nil, err
		}
		return &defaultKpmConf, nil
	}
}
