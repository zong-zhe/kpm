// Copyright 2021 The KCL Authors. All rights reserved.

package opt

// Config represents some configurations used by kpm.
type Config struct {
	opt      *Options
	HomePath string
}

type Options struct {
	init *InitOptions
	add  *AddOptions
}

type InitOptions struct {
	Name     string
	InitPath string
}

type AddOptions struct {
	LocalPath    string
	RegistryOpts RegistryOption
}

type RegistryOption struct {
	Git *GitOption
}

type GitOption struct {
	Url    string
	Branch string
	Commit string
	Tag    string
}

func NewGitOption() *GitOption {
	return &GitOption{}
}
