package pkg

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func (dep *Dependency) Download() (string, error) {
	if d, ok := dep.GetDependency().(*Dependency_Git); ok {
		return d.Git.Download()
	}
	return "", nil
}

const kpmHome = "KPM_HOME"
const gitPkg = "git"

func (dep *GitDependency) Download() (string, error) {
	repoURL := dep.GetGit()
	localPath := parseLocalPathFromGitUrl(dep.GetGit())

	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	// switch branch

	// switch commit

	// switch tag

	return localPath, err
}

func parseLocalPathFromGitUrl(gitUrl string) string {
	parsedUrl, _ := url.Parse(gitUrl)
	pathWithoutScheme := strings.TrimPrefix(gitUrl, parsedUrl.Scheme+"://")

	fileExt := filepath.Ext(pathWithoutScheme)
	return filepath.Join(kpmHome, filepath.Base(pathWithoutScheme[:len(pathWithoutScheme)-len(fileExt)]))
}
