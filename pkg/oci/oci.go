package oci

import (
	"context"
	"path/filepath"

	"kusionstack.io/kpm/pkg/errors"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

// Pull will pull the oci atrifacts from oci registry to local path.
func Pull(localPath, reg, repoName, tag string) (string, error) {
	// 0. Create a file store
	fs, err := file.New(localPath)
	if err != nil {
		return "", err
	}
	defer fs.Close()

	// 1. Connect to a remote repository
	ctx := context.Background()
	repo, err := remote.NewRepository(filepath.Join(reg, repoName))
	if err != nil {
		return "", err
	}

	// 2. Copy from the remote repository to the file store
	_, err = oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		return "", err
	}

	// 3.Get the (*.tar) file path.
	matches, err := filepath.Glob(filepath.Join(localPath, "*.tar"))
	if err != nil && len(matches) != 1 {
		return "", errors.FailedPullFromOci
	}

	return matches[0], nil
}
