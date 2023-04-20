package settings

const CONFIG_JSON_PATH = "~/.kpm/config/config.json"

type Settings struct {
	CredentialsFile string
}

func Init() Settings {
	return Settings{
		CredentialsFile: CONFIG_JSON_PATH,
	}
}
