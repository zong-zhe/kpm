package downloader

type GitCache struct {
	*PkgCache
}

// Save is copy or move the package under 'srcPath' to the cache path.
func (gc *GitCache) Save(key, srcPath string) error {
	// TODO: implement this method
	return nil
}

// Find is find the package by the key and return the path of the package in cache.
func (gc *GitCache) Find(key string) (string, error) {
	// TODO: implement this method
	return "", nil
}

// Exists is check whether the package exists in the cache.
func (gc *GitCache) Exists(key string) bool {
	// TODO: implement this method
	return false
}

// Remove is remove the package from the cache, delete the package in local path.
func (gc *GitCache) Remove(key string) error {
	// TODO: implement this method
	return nil
}

// SetCachePath is set the cache path.
func (gc *GitCache) SetCachePath(path string) {
	// TODO: implement this method
}

// GetCachePath is get the cache path.
func (gc *GitCache) GetCachePath() string {
	// TODO: implement this method
	return ""
}
