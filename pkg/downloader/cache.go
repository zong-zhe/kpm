package downloader

import "os"

// key is the full source url of the package.
type Cache interface {
	// Save is copy or move the package under 'srcPath' to the cache path.
	Update(source Source, updateFunc func(cachePath string) error) error
	// Find is find the package by the key and return the path of the package in cache.
	Find(source Source) (string, error)
	// Remove is remove the package from the cache, delete the package in local path.
	Remove(source Source) error
	// RemoveAll is remove all the packages from the cache.
	RemoveAll() error
}

type PkgCache struct {
	cacheDir string
}

func (p *PkgCache) RemoveAll() error {
	if err := os.RemoveAll(p.cacheDir); err != nil {
		return err
	}
	return nil
}
