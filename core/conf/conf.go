// Copyright 2021 The KCL Authors. All rights reserved.

package conf

// Config represents some configurations used by kpm.
type Config struct {
	Name       string
	Edition    string
	Version    string
	KclModPath string
}

const defaultVerion = "0.0.1"
const defaultEdition = "0.0.1"

// NewEmptyConf returns a default configuration for kcl.mod.
// Default version of kcl.mod is "0.0.1".
// Default edition of kcl.mod is "0.0.1".
func NewEmptyConf() Config {
	return Config{
		Edition: defaultEdition,
		Version: defaultVerion,
	}
}

// SetName sets name for struct Config.
func (conf Config) SetName(name string) Config {
	conf.Name = name
	return conf
}

// SetEdition sets edition for struct Config.
func (conf Config) SetEdition(edit string) Config {
	conf.Edition = edit
	return conf
}

// SetVersion sets Version for struct Config.
func (conf Config) SetVersion(version string) Config {
	conf.Version = version
	return conf
}

func (conf Config) SetKclModPath(execPath string) Config {
	conf.KclModPath = execPath
	return conf
}
