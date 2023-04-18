// Copyright 2023 The KCL Authors. All rights reserved.

package opt

import (
	"fmt"
	"path/filepath"
	"strings"

	"kusionstack.io/kpm/pkg/errors"
	"kusionstack.io/kpm/pkg/reporter"
)

// Input options of 'kpm init'.
type InitOptions struct {
	Name     string
	InitPath string
}

func (opts *InitOptions) Validate() error {
	if len(opts.Name) == 0 {
		return errors.InvalidInitOptions
	} else if len(opts.InitPath) == 0 {
		return errors.InternalBug
	}
	return nil
}

type AddOptions struct {
	LocalPath    string
	RegistryOpts RegistryOptions
}

func (opts *AddOptions) Validate() error {
	if len(opts.LocalPath) == 0 {
		return errors.InternalBug
	} else if opts.RegistryOpts.Git != nil {
		return opts.RegistryOpts.Git.Validate()
	} else if opts.RegistryOpts.Git == nil {
		return errors.InvalidAddOptionsWithoutRegistry
	}
	return nil
}

type RegistryOptions struct {
	Git *GitOptions
	Oci *OciOptions
}

type GitOptions struct {
	Url    string
	Branch string
	Commit string
	Tag    string
}

func (opts *GitOptions) Validate() error {
	if len(opts.Url) == 0 {
		return errors.InvalidAddOptionsInvalidGitUrl
	} else if len(opts.Tag) == 0 {
		return errors.InvalidAddOptionsInvalidTag
	}
	return nil
}

const DEFAULT_REGISTRY = "docker.io"
const DEFAULT_OCI_TAG = "latest"

func GetOCIReg() string {
	return DEFAULT_REGISTRY
}

func GetDefaultOCITag() string {
	return DEFAULT_OCI_TAG
}

type OciOptions struct {
	Reg  string
	Repo string
	Tag  string
}

func (opts *OciOptions) Validate() error {
	opts.Reg = GetOCIReg()
	if len(opts.Repo) == 0 {
		return errors.InvalidAddOptionsInvalidOciRepo
	} else if len(opts.Tag) == 0 {
		return errors.InvalidAddOptionsInvalidTag
	}
	return nil
}

const OCI_SEPARATOR = ":"

// ParseOciOptionFromString will parser '<repo_name>:<repo_tag>' into an 'OciOptions' with an OCI registry.
// the default OCI registry is 'docker.io'.
// if the 'ociUrl' is only '<repo_name>', ParseOciOptionFromString will take 'latest' as the default tag.
func ParseOciOptionFromString(ociUrl string) (*OciOptions, error) {
	oci_address := strings.Split(ociUrl, OCI_SEPARATOR)
	if len(oci_address) == 1 {
		reporter.Report("kpm: using default tag: latest")
		return &OciOptions{
			Reg:  GetOCIReg(),
			Repo: oci_address[0],
			Tag:  GetDefaultOCITag(),
		}, nil
	} else if len(oci_address) == 2 {
		return &OciOptions{
			Reg:  GetOCIReg(),
			Repo: oci_address[0],
			Tag:  oci_address[1],
		}, nil
	} else {
		return nil, errors.InvalidOciRef
	}
}

// AddStoragePathSuffix will take 'Registry/Repo/Tag' as a path suffix.
// e.g. Take '/usr/test' as input,
// and oci options is
//
// OciOptions {
//   Reg: 'docker.io',
//   Repo: 'test/testRepo',
//   Tag: 'v0.0.1'
// }
//
// You will get a path '/usr/test/docker.io/test/testRepo/v0.0.1'.
func (oci *OciOptions) AddStoragePathSuffix(pathPrefix string) string {
	return filepath.Join(filepath.Join(filepath.Join(pathPrefix, oci.Reg), oci.Repo), oci.Tag)
}

// The parameters needed to compile the kcl program.
type KclvmOptions struct {
	Deps       map[string]string
	EntryFiles []string
	// todo: add all kclvm options.
}

func (opts *KclvmOptions) Validate() error {
	if len(opts.EntryFiles) == 0 {
		return errors.InvalidRunOptionsWithoutEntryFiles
	}
	return nil
}

func NewKclvmOpts() *KclvmOptions {
	return &KclvmOptions{
		Deps:       make(map[string]string),
		EntryFiles: make([]string, 0),
	}
}

// Generate the kcl compile command arguments based on 'KclvmOptions'.
func (kclOpts *KclvmOptions) Args() []string {
	var args []string
	args = append(args, kclOpts.EntryFiles...)
	args = append(args, kclOpts.PkgPathMapArgs()...)
	return args
}

const EXTERNAL_ARG = "-E"
const EXTERNAL_ARG_PATTERN = "%s=%s"

// Generate the kcl compile command arguments '-E <pkg_name>=<pkg_path>'.
func (kclOpts *KclvmOptions) PkgPathMapArgs() []string {
	var args []string
	for k, v := range kclOpts.Deps {
		args = append(args, EXTERNAL_ARG)
		args = append(args, fmt.Sprintf(EXTERNAL_ARG_PATTERN, k, v))
	}
	return args
}
