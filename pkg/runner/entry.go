package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/golang-collections/collections/set"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/errors"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/utils"
)

// EntryKind is the kind of the entry.
// Including:
// 1. WithKclMod: local file which can find 'kcl.mod' in the parent dir of the file.
// 2. WithoutKclMod: local file which can find 'kcl.mod' in the parent dir of the file.
// 3. TarEntry: kcl package tar file.
// 4. UrlEntry: kcl package url.
// 5. RefEntry: kcl package ref.
type EntryKind string

// Entry is the entry of 'kpm run'.
// 一个编译入口，主要包括
// 1. packageSource: 包源，可以是文件路径、tar路径、url或者ref
// 2. entryFiles: 要编译的全部文件，注意，这个里面只能是文件，或者文件目录，不能再是别的了
// 3. kind: 入口的类型，可以是文件、tar、url或者ref
// 4. KclPkg: 如果是 kcl 包，这个包里主要负责承载一些 kcl.mod 中的配置信息
type Entry struct {
	// The package source of the entry, filepath, tar path, url or ref.
	packageSource string
	// The start files for one compilation.
	entryFiles []string
	// The kind of the entry, file, tar, url or ref.
	kind EntryKind
}

// SetKind will set the kind of the entry.
func (e *Entry) SetKind(kind EntryKind) {
	e.kind = kind
}

// Kind will return the kind of the entry.
func (e *Entry) Kind() EntryKind {
	return e.kind
}

// IsLocalFileWithKclMod will return true if the entry is a local file with 'kcl.mod'.
func (e *Entry) IsLocalFileWithKclMod() bool {
	return e.kind == constants.FileWithKclModEntry
}

// IsLocalFile will return true if the entry is a local file.
func (e *Entry) IsLocalFile() bool {
	return e.kind == constants.FileEntry
}

// IsUrl will return true if the entry is a url.
func (e *Entry) IsUrl() bool {
	return e.kind == constants.UrlEntry
}

// IsRef will return true if the entry is a ref.
func (e *Entry) IsRef() bool {
	return e.kind == constants.RefEntry
}

// IsTar will return true if the entry is a tar.
func (e *Entry) IsTar() bool {
	return e.kind == constants.TarEntry
}

func (e *Entry) IsGit() bool {
	return e.kind == constants.GitEntry
}

// IsEmpty will return true if the entry is empty.
func (e *Entry) IsEmpty() bool {
	return len(e.packageSource) == 0
}

// PackageSource will return the package source of the entry.
func (e *Entry) PackageSource() string {
	return e.packageSource
}

// EntryFiles will return the entry files of the entry.
func (e *Entry) EntryFiles() []string {
	return e.entryFiles
}

// SetPackageSource will set the package source of the entry.
func (e *Entry) SetPackageSource(packageSource string) {
	e.packageSource = packageSource
}

// AddEntryFile will add a entry file to the entry.
func (e *Entry) AddEntryFile(entrySource string) {
	e.entryFiles = append(e.entryFiles, entrySource)
}

// FindRunEntryFrom will find the entry of the compilation from the entry sources.
func FindRunEntryFrom(sources []string) (*Entry, *reporter.KpmEvent) {
	entry := Entry{}
	// modPathSet is used to check if there are multiple packages to be compiled at the same time.
	// It is a set of the package source so that the same package source will only be added once.
	var modPathSet = set.New()
	for _, source := range sources {
		// If the entry is a local file but not a tar file,
		if utils.DirExists(source) && !utils.IsTar(source) {
			// Find the 'kcl.mod'
			modPath, err := FindModRootFrom(source)
			if err != (*reporter.KpmEvent)(nil) {
				// If the 'kcl.mod' is not found,
				if err.Type() == reporter.KclModNotFound {
					if utils.IsKfile(source) {
						// If the entry is a kcl file, the parent dir of the kcl file will be package path.
						modPath = filepath.Dir(source)
					} else {
						// If the entry is a dir, the dir will be package path.
						modPath = source
					}
				} else {
					return nil, err
				}
			}
			entry.SetPackageSource(modPath)
			entry.AddEntryFile(source)
			if !utils.DirExists(filepath.Join(modPath, constants.KCL_MOD)) {
				entry.SetKind(constants.FileEntry)
			} else {
				entry.SetKind(constants.FileWithKclModEntry)
			}
			absModPath, bugerr := filepath.Abs(modPath)
			if bugerr != nil {
				return nil, reporter.NewErrorEvent(reporter.Bug, bugerr, errors.InternalBug.Error())
			}
			modPathSet.Insert(absModPath)
		} else if utils.IsURL(source) || utils.IsRef(source) || utils.IsTar(source) {
			modPathSet.Insert(source)
			entry.SetPackageSource(source)
			entry.SetKind(GetSourceKindFrom(source))
		}
	}

	// kpm only allows one package to be compiled at a time.
	if modPathSet.Len() > 1 {
		// sort the mod paths to make the error message more readable.
		var modPaths []string
		setModPathsMethod := func(modpath interface{}) {
			p, ok := modpath.(string)
			if !ok {
				modPaths = append(modPaths, "")
			} else {
				modPaths = append(modPaths, p)
			}
		}
		modPathSet.Do(setModPathsMethod)
		sort.Strings(modPaths)
		return nil, reporter.NewErrorEvent(
			reporter.CompileFailed,
			fmt.Errorf("cannot compile multiple packages %s at the same time", modPaths),
			"only allows one package to be compiled at a time",
		)
	}

	return &entry, nil
}

// GetSourceKindFrom will return the kind of the source.
func GetSourceKindFrom(source string) EntryKind {
	if utils.DirExists(source) && !utils.IsTar(source) {
		return constants.FileEntry
	} else if utils.IsTar(source) {
		return constants.TarEntry
	} else if utils.IsGitRepoUrl(source) {
		return constants.GitEntry
	} else if utils.IsURL(source) {
		return constants.UrlEntry
	} else if utils.IsRef(source) {
		return constants.RefEntry
	}
	return ""
}

// FindModRootFrom will find the kcl.mod path from the start path.
func FindModRootFrom(startPath string) (string, *reporter.KpmEvent) {
	info, err := os.Stat(startPath)
	if err != nil {
		return "", reporter.NewErrorEvent(reporter.CompileFailed, err, fmt.Sprintf("failed to access path '%s'", startPath))
	}
	var start string
	// If the start path is a kcl file, find from the parent dir of the kcl file.
	if utils.IsKfile(startPath) {
		start = filepath.Dir(startPath)
	} else if info.IsDir() {
		// If the start path is a dir, find from the start path.
		start = startPath
	} else {
		return "", reporter.NewErrorEvent(reporter.CompileFailed, err, fmt.Sprintf("invalid file path '%s'", startPath))
	}

	if _, err := os.Stat(filepath.Join(start, constants.KCL_MOD)); err == nil {
		return start, nil
	} else {
		parent := filepath.Dir(startPath)
		if parent == startPath {
			return "", reporter.NewErrorEvent(reporter.KclModNotFound, fmt.Errorf("cannot find kcl.mod in '%s'", startPath))
		}
		return FindModRootFrom(filepath.Dir(startPath))
	}
}
