package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

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
		fmt.Println(path)
		if _, err := io.Copy(hasher, f); err != nil {
			return err
		}

		return nil
	})

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func RenameDir(oldDir string, newDir string) {
	// 复制目录及其内容到新位置
	err := copy.Copy(oldDir, newDir)
	if err != nil {
		log.Fatal(err)
	}

	// 删除旧目录
	err = os.RemoveAll(oldDir)
	if err != nil {
		log.Fatal(err)
	}
}
