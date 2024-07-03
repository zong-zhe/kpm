package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"kcl-lang.io/kpm/pkg/downloader"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/utils"
)

type Visitor interface {
	Visit(source *downloader.Source, visitor func(*pkg.KclPkg) error) error
}

type PkgVisitor struct {
	kpmcli *KpmClient
}

type RemotePkgVisitor struct {
	*PkgVisitor
}

func NewRemotePkgVisitor(kpmcli *KpmClient) *RemotePkgVisitor {
	return &RemotePkgVisitor{
		PkgVisitor: &PkgVisitor{
			kpmcli: kpmcli,
		},
	}
}

func (v *RemotePkgVisitor) Visit(source *downloader.Source, visit func(*pkg.KclPkg) error) error {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if source.Git != nil {
		tmpDir = filepath.Join(tmpDir, "tmp")
	}

	err = v.kpmcli.DepDownloader.Download(*downloader.NewDownloadOptions(
		downloader.WithLocalPath(tmpDir),
		downloader.WithSource(*source),
		downloader.WithLogWriter(v.kpmcli.GetLogWriter()),
		downloader.WithSettings(*v.kpmcli.GetSettings()),
	))

	if err != nil {
		return err
	}

	pkg, err := v.kpmcli.LoadPkgFromPath(tmpDir)
	if err != nil {
		return err
	}

	return visit(pkg)
}

type FileLocalPkgVisitor struct {
	*PkgVisitor
}

func NewFileLocalPkgVisitor(kpmcli *KpmClient) *FileLocalPkgVisitor {
	return &FileLocalPkgVisitor{
		PkgVisitor: &PkgVisitor{
			kpmcli: kpmcli,
		},
	}
}

func (v *FileLocalPkgVisitor) Visit(source *downloader.Source, visit func(*pkg.KclPkg) error) error {
	if !source.IsLocalPath() {
		return fmt.Errorf("source is not a local path")
	}

	rootPath, err := source.FindRootPath()
	if err != nil {
		return err
	}

	vKclModPath := filepath.Join(rootPath, pkg.MOD_FILE)
	vKclModLockPath := filepath.Join(rootPath, pkg.MOD_LOCK_FILE)
	if !utils.DirExists(vKclModPath) {
		logWriter := v.kpmcli.GetLogWriter()
		v.kpmcli.SetLogWriter(nil)

		defer func() {
			if err := os.Remove(vKclModPath); err != nil {
				log.Printf("Failed to remove %s: %v", vKclModPath, err)
			}
			if err := os.Remove(vKclModLockPath); err != nil {
				log.Printf("Failed to remove %s: %v", vKclModLockPath, err)
			}
		}()

		initOpts := opt.InitOptions{
			Name:     "vPkg_" + uuid.New().String(),
			InitPath: rootPath,
		}

		kclPkg := pkg.NewKclPkg(&initOpts)
		err := v.kpmcli.createIfNotExist(kclPkg.ModFile.GetModFilePath(), kclPkg.ModFile.StoreModFile)
		if err != nil {
			return err
		}

		v.kpmcli.SetLogWriter(logWriter)
	}

	pkg, err := v.kpmcli.LoadPkgFromPath(rootPath)
	if err != nil {
		return err
	}

	return visit(pkg)
}

type TarLocalPkgVisitor struct {
	*PkgVisitor
}

func NewTarLocalPkgVisitor(kpmcli *KpmClient) *TarLocalPkgVisitor {
	return &TarLocalPkgVisitor{
		PkgVisitor: &PkgVisitor{
			kpmcli: kpmcli,
		},
	}
}

func (v *TarLocalPkgVisitor) Visit(source *downloader.Source, visit func(*pkg.KclPkg) error) error {
	if !source.IsLocalTarPath() {
		return fmt.Errorf("source is not a local tar path")
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	sourcePath, err := source.ToFilePath()
	if err != nil {
		return err
	}
	err = utils.UnTarDir(sourcePath, tmpDir)
	if err != nil {
		return err
	}

	pkg, err := v.kpmcli.LoadPkgFromPath(tmpDir)
	if err != nil {
		return err
	}

	return visit(pkg)
}

func NewVisitor(source *downloader.Source, kpmcli *KpmClient) Visitor {
	if source.IsLocalTarPath() {
		return NewTarLocalPkgVisitor(kpmcli)
	} else if source.IsLocalPath() {
		return NewFileLocalPkgVisitor(kpmcli)
	} else {
		return NewRemotePkgVisitor(kpmcli)
	}
}
