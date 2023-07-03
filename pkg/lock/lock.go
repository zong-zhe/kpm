package lock

type FileLock struct {
	filePath string
}

func NewFileLock(filePath string) *FileLock {
	return &FileLock{
		filePath: filePath,
	}
}

func (fileLock *FileLock) TryLock() error {
	return nil
}
