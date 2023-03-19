// Copyright 2022 The KCL Authors. All rights reserved.

// Package mod is the core abstraction in kpm for working with a package or module of kcl.
package mod

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
	gogit "github.com/go-git/go-git/v5"
	"kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/utils"
)

const (
	File     = "kcl.mod"
	LockFile = "kcl.mod.lock"
)

type Package struct {
	Name    string
	Edition string
	Version string
}

type ModFile struct {
	HomePath     string
	Pkg          Package
	Dependencies Dependencies
}

type ModLockFile struct {
	HomePath     string
	Dependencies Dependencies
}

type Dependencies struct {
	Deps map[string]Dependency
}

type Dependency struct {
	Name    string
	Source  Source
	Version string
	Sum     string
}

type Source struct {
	GitSource *Git
}

type Git struct {
	Url    string
	Branch string
	Commit string
	Tag    string
}

func ModFileExists(path string) (bool, error) {
	return exists(filepath.Join(path, File))
}

func ModLockFileExists(path string) (bool, error) {
	return exists(filepath.Join(path, LockFile))
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func LoadModFile(homePath string) (*ModFile, error) {

	readFile, err := os.OpenFile(filepath.Join(homePath, File), os.O_RDWR, 0644)
	defer readFile.Close()

	modFile := new(ModFile)
	_, err = toml.NewDecoder(readFile).Decode(&modFile)

	modFile.HomePath = homePath

	deps, err := LoadModLockFile(homePath)
	modFile.Dependencies = *deps

	if err != nil {
		return nil, err
	}
	return modFile, nil
}

func LoadModLockFile(homePath string) (*Dependencies, error) {

	readFile, err := os.OpenFile(filepath.Join(homePath, LockFile), os.O_RDWR, 0644)
	defer readFile.Close()

	locks := new(Dependencies)
	_, err = toml.NewDecoder(readFile).Decode(&locks)

	if err != nil {
		return nil, err
	}
	return locks, nil
}

// 这个里面要把sum字段屏蔽
func (mfile *ModFile) Store() error {
	file, err := os.Create(filepath.Join(mfile.HomePath, File))
	if err != nil {
		reporter.ExitWithReport("Error creating file:", err)
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(mfile); err != nil {
		reporter.ExitWithReport("Error encoding TOML:", err)
		return err
	}
	return nil
}

func (mfile *ModLockFile) Store() error {
	file, err := os.Create(filepath.Join(mfile.HomePath, LockFile))
	if err != nil {
		reporter.ExitWithReport("Error creating file:", err)
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(mfile); err != nil {
		reporter.ExitWithReport("Error encoding TOML:", err)
		return err
	}
	return nil
}

func NewModFile(conf *conf.Config, homePath string) *ModFile {
	return &ModFile{
		HomePath: homePath,
		Pkg: Package{
			Name:    conf.Name,
			Version: conf.Version,
			Edition: conf.Edition,
		},
	}
}

func (dep *Dependency) Download(localPath string) (*Dependency, error) {
	if dep.Source.GitSource != nil {
		dep.Source.GitSource.Download(localPath)
		dep.Sum = utils.HashDir(filepath.Join(localPath, dep.Name))
	}
	return dep, nil
}

func (dep *Git) Download(localPath string) (string, error) {
	repoURL := dep.Url
	repo, err := git.Clone(repoURL, localPath)

	if err != nil {
		reporter.Report("kpm: git clone error:", err)
		return localPath, err
	}
	// checkout branch

	err = git.CheckoutBranch(repo, dep.Branch)

	if err != nil {
		return localPath, reportCheckoutErr(repo, err)
	}

	// checkout commit

	err = git.CheckoutCommit(repo, dep.Commit)

	if err != nil {
		return localPath, reportCheckoutErr(repo, err)
	}

	// checkout tag

	err = git.CheckoutTag(repo, dep.Tag)

	if err != nil {
		return localPath, reportCheckoutErr(repo, err)
	}

	return localPath, err
}

func reportCheckoutErr(repo *gogit.Repository, origin_err error) error {
	ref, err := repo.Head()

	if err != nil {
		reporter.Fatal("kpm: internal bug, please contact us to fix it.")
	}

	reporter.Report("kpm: checkout error:", origin_err, ".")
	reporter.Report("kpm:", ref.Hash().String(), " is used.")

	return origin_err
}

func ParseUrl(localPath string, url string) *Dependency {
	return nil
}

func DepEqual(d1, d2 Dependency) bool {
	name := d1.Name == d2.Name
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}
