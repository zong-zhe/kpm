package store

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter"
)

const kpmHome = "KPM_HOME"
const pkgMetaDataPath = "metadata"
const gitPkg = "git"



func GetGitDepInLocalPath(localPath string, gitUrl *url.URL) (string, error) {
	gitGetter := getter.GitGetter{}
	storePath := filepath.Join(os.Getenv(kpmHome), gitPkg, localPath)
	err := gitGetter.Get(storePath, gitUrl)
	return storePath, err
}

func ParseLocalPathFromGitUrl(gitUrl string) string {
	parsedUrl, _ := url.Parse(gitUrl)
	pathWithoutScheme := strings.TrimPrefix(gitUrl, parsedUrl.Scheme+"://")

	fileExt := filepath.Ext(pathWithoutScheme)
	return filepath.Base(pathWithoutScheme[:len(pathWithoutScheme)-len(fileExt)])
}
