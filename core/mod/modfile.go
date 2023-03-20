// Copyright 2022 The KCL Authors. All rights reserved.

// Package mod is the core abstraction in kpm for working with a package or module of kcl.
package modfile

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	gogit "github.com/go-git/go-git/v5"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/opt"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/utils"
)

const (
	File     = "kcl.mod"
	LockFile = "kcl.mod.lock"
)

type Package struct {
	Name    string `toml:"name"`
	Edition string `toml:"edition"`
	Version string `toml:"version"`
}

type ModFile struct {
	HomePath string  `toml:"-"`
	Pkg      Package `toml:"package"`
	Dependencies
}

type ModLockFile struct {
	HomePath string `toml:"-"`
	Dependencies
}

type Dependencies struct {
	Deps map[string]Dependency `toml:"dependencies"`
}

func (d *Dependencies) MarshalTOML() ([]byte, error) {
	if d == nil {
		return nil, nil
	}

	buff := new(bytes.Buffer)
	encoder := toml.NewEncoder(buff)

	for _, v := range d.Deps {
		err := encoder.Encode(struct {
			Name string
		}{
			Name: string(v.Url),
		})

		if err != nil {
			return nil, fmt.Errorf("kpm: Internal bug")
		}
	}

	return buff.Bytes(), nil
}

type Dependency struct {
	Name string `toml:"name"`
	Source
	Version string `toml:"version"`
	Sum     string `toml:"sum"`
}

type Source struct {
	*Git
}

type Git struct {
	Url    string `toml:"url"`
	Branch string `toml:"branch"`
	Commit string `toml:"commit"`
	Tag    string `toml:"tag"`
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
	modFile := new(ModFile)
	err := loadFile(homePath, File, modFile)
	if err != nil {
		return nil, err
	}

	modFile.HomePath = filepath.Join(homePath, File)

	if modFile.Dependencies.Deps == nil {
		modFile.Dependencies.Deps = make(map[string]Dependency)
	}

	return modFile, nil
}

func LoadModLockFile(homePath string) (*ModLockFile, error) {

	locks := new(ModLockFile)
	err := loadFile(homePath, LockFile, locks)
	if err != nil {
		return nil, err
	}

	if locks.Dependencies.Deps == nil {
		locks.Dependencies.Deps = make(map[string]Dependency)
	}

	locks.HomePath = filepath.Join(homePath, LockFile)

	return locks, nil
}

func (mfile *ModFile) Store() error {
	return storeToFile(mfile.HomePath, mfile)
}

func (mfile *ModLockFile) Store() error {
	return storeToFile(mfile.HomePath, mfile)
}

const defaultVerion = "0.0.1"
const defaultEdition = "0.0.1"

func NewModFile(opt *opt.InitOptions, homePath string) *ModFile {
	return &ModFile{
		HomePath: homePath,
		Pkg: Package{
			Name:    opt.Name,
			Version: defaultVerion,
			Edition: defaultEdition,
		},
	}
}

func NewModLockFile(opt *opt.InitOptions, homePath string) *ModLockFile {
	return &ModLockFile{
		HomePath: homePath,
		Dependencies: Dependencies{
			Deps: make(map[string]Dependency),
		},
	}
}

func (dep *Dependency) Download(localPath string) (*Dependency, error) {
	if dep.Source.Git != nil {
		dep.Source.Git.Download(localPath)
		dep.Sum = utils.HashDir(localPath)
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

	// checkout branch/commit/tag
	err = checkout(repo, dep.Branch, dep.Commit, dep.Tag)

	if err != nil {
		return localPath, reportCheckoutErr(repo, err)
	}

	return localPath, err
}

func checkout(repo *gogit.Repository, branch string, commit string, tag string) error {
	var err error = nil
	if len(branch) != 0 {
		err = git.CheckoutBranch(repo, branch)
	} else if len(commit) != 0 {
		err = git.CheckoutCommit(repo, commit)
	} else if len(tag) != 0 {
		err = git.CheckoutTag(repo, tag)
	}
	return err
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

func ParseOpt(opt *opt.RegistryOption) *Dependency {

	if opt.Git != nil {
		gitSource := Git{
			Url:    opt.Git.Url,
			Branch: opt.Git.Branch,
			Commit: opt.Git.Commit,
			Tag:    opt.Git.Tag,
		}

		name := ParseRepoNameFromGitUrl(gitSource.Url)

		return &Dependency{
			Name: name,
			Source: Source{
				Git: &gitSource,
			},
		}
	}

	return nil
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

func storeToFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		reporter.ExitWithReport("kpm: failed to create file: ", filePath, err)
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(data); err != nil {
		reporter.ExitWithReport("kpm: failed to encode TOML:", err)
		return err
	}
	return nil
}

func loadFile(homePath string, fileName string, v interface{}) error {
	readFile, err := os.OpenFile(filepath.Join(homePath, fileName), os.O_RDWR, 0644)
	if err != nil {
		reporter.Report("kpm: failed to load", fileName)
		return err
	}
	defer readFile.Close()

	_, err = toml.NewDecoder(readFile).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
