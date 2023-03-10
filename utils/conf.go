package conf

type Config struct {
	Name     string
	Edition  string
	Version  string
	ExecPath string
}

const defaultVerion = "0.0.1"
const defaultEdition = "0.0.1"

func NewEmptyConf() Config {
	return Config{
		Edition: defaultEdition,
		Version: defaultVerion,
	}
}

func (conf Config) SetName(name string) Config {
	conf.Name = name
	return conf
}

func (conf Config) SetEdition(edit string) Config {
	conf.Edition = edit
	return conf
}

func (conf Config) SetVersion(version string) Config {
	conf.Version = version
	return conf
}

func (conf Config) SetExecPath(execPath string) Config {
	conf.ExecPath = execPath
	return conf
}
