// Copyright 2022 The KCL Authors. All rights reserved.

// Package mod is the core abstraction in kpm for working with a package or module of kcl.
package mod

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

	modFile.HomePath = filepath.Join(homePath, File)
	if modFile.Dependencies.Deps == nil {
		modFile.Dependencies.Deps = make(map[string]Dependency)
	}

	if err != nil {
		return nil, err
	}
	return modFile, nil
}

func LoadModLockFile(homePath string) (*ModLockFile, error) {

	readFile, err := os.OpenFile(filepath.Join(homePath, LockFile), os.O_RDWR, 0644)
	defer readFile.Close()

	locks := new(ModLockFile)
	_, err = toml.NewDecoder(readFile).Decode(&locks)

	if locks.Dependencies.Deps == nil {
		locks.Dependencies.Deps = make(map[string]Dependency)
	}

	locks.HomePath = filepath.Join(homePath, LockFile)

	if err != nil {
		return nil, err
	}
	return locks, nil
}

// 这个里面要把sum字段屏蔽
func (mfile *ModFile) Store() error {
	file, err := os.Create(mfile.HomePath)
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
	file, err := os.Create(mfile.HomePath)
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

func NewModLockFile(conf *conf.Config, homePath string) *ModLockFile {
	return &ModLockFile{
		HomePath: homePath,
		Dependencies: Dependencies{
			Deps: make(map[string]Dependency),
		},
	}
}

func (dep *Dependency) Download(localPath string) (*Dependency, error) {
	if dep.Source.GitSource != nil {
		dep.Source.GitSource.Download(localPath)

		dep.Sum = utils.HashDir(filepath.Join(localPath, dep.Name))

		// modfile, _ := LoadModFile(localPath)
		// utils.RenameDir(localPath, filepath.Join(filepath.Dir(localPath), modfile.Pkg.Name))
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

func ParseOpt(opt *git.GitOption) Dependency {

	gitSource := Git{
		Url:    opt.Url,
		Branch: opt.Branch,
		Commit: opt.Commit,
		Tag:    opt.Tag,
	}

	name := ParseRepoNameFromGitUrl(gitSource.Url)

	return Dependency{
		Name: name,
		Source: Source{
			GitSource: &gitSource,
		},
	}
}

func DepEqual(d1, d2 Dependency) bool {
	name := d1.Name == d2.Name
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}

func ParseLocalPathFromGitUrl(rootPath string, gitUrl string) string {
	parsedUrl, _ := url.Parse(gitUrl)
	pathWithoutScheme := strings.TrimPrefix(gitUrl, parsedUrl.Scheme+"://")

	fileExt := filepath.Ext(pathWithoutScheme)
	return filepath.Join(rootPath, filepath.Base(pathWithoutScheme[:len(pathWithoutScheme)-len(fileExt)]))
}

func ParseRepoNameFromGitUrl(gitUrl string) string {
	name := filepath.Base(gitUrl)
	return name[:len(name)-len(filepath.Ext(name))]
}
