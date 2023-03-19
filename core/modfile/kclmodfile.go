package kclmodfile

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	gogit "github.com/go-git/go-git/v5"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/reporter"
)

const (
	File     = "kcl.mod"
	LockFile = "kcl.mod.lock"
)

var (
	ErrUpdateJB = errors.New("jsonnetfile version unknown, update jb")
)

type KclModFile struct {
	HomePath     string
	Version      string
	Dependencies map[string]Dependency
}

type Dependency struct {
	Name    string
	Source  Source
	Version string
	Sum     string
}

type Source struct {
	GitSource *Git `json:"git,omitempty"`
}

type Git struct {
	Url    string
	Branch string
	Commit string
	Tag    string
}

func NewGitDependency() Dependency {
	// todo
	return Dependency{}
}

func (dep *Dependency) Download(localPath string) (*Dependency, error) {
	if dep.Source.GitSource != nil {
		localpath, err := dep.Source.GitSource.DownloadFromGitSource(localPath)
		if err != nil {
			return dep, err
		}
		dep.Sum = HashDir(localpath)
		depModFile, err := LoadModFile(localpath)
		dep.Version = depModFile.Version
		if err != nil {
			return dep, err
		}

	}
	return dep, nil
}

func (depGit *Git) DownloadFromGitSource(vendorPath string) (string, error) {
	repoURL := depGit.Url

	localPath := git.ParseLocalPathFromGitUrl(vendorPath, repoURL)
	repo, err := git.Clone(repoURL, localPath)

	if err != nil {
		reporter.Report("kpm: git clone error:", err)
		return localPath, err
	}
	// checkout branch

	err = git.CheckoutBranch(repo, depGit.Branch)

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	// checkout commit

	err = git.CheckoutCommit(repo, depGit.Commit)

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	// checkout tag

	err = git.CheckoutTag(repo, depGit.Tag)

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	return localPath, err
}

func report_checkout_err(repo *gogit.Repository, origin_err error) error {
	ref, err := repo.Head()

	if err != nil {
		reporter.Fatal("kpm: internal bug, please contact us to fix it.")
	}

	reporter.Report("kpm: checkout error:", origin_err, ".")
	reporter.Report("kpm:", ref.Hash().String(), " is used.")

	return origin_err
}

// Exists returns whether the file at the given path exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func LoadModFile(homePath string) (*KclModFile, error) {

	readFile, err := os.OpenFile(homePath, os.O_RDWR, 0644)
	defer readFile.Close()

	modFile := new(KclModFile)
	_, err = toml.NewDecoder(readFile).Decode(&modFile)

	modFile.HomePath = homePath

	if err != nil {
		return nil, err
	}
	return modFile, nil
}

// hashDir computes the checksum of a directory by concatenating all files and
// hashing this data using sha256. This can be memory heavy with lots of data,
// but jsonnet files should be fairly small
func HashDir(dir string) string {
	hasher := sha256.New()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(hasher, f); err != nil {
			return err
		}

		return nil
	})

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
