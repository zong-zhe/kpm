package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/downloader"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/runner"
	"kcl-lang.io/kpm/pkg/utils"
)

// RunOptions contains the options for running a kcl package.
type RunOptions struct {
	// CompileOptions is the options for kcl compiler.
	hasSettingsYaml bool
	SettingYamls    []string
	Vendor          bool
	workDir         string
	// Sources is the sources of the package.
	// It can be a local *.k path, a local *.tar path, a local directory, a remote git/oci path,.
	Sources []*downloader.Source
	*kcl.Option
}

// Only one kcl module can be compiled at a time.
// So, the sources must have the same root path.
func (ro *RunOptions) Validate() error {
	if len(ro.Sources) == 0 {
		return errors.New("no source provided")
	}

	// More than one source, all sources must have the same root path.
	if len(ro.Sources) > 1 {
		rootPath, err := ro.Sources[0].FindRootPath()
		if err != nil {
			return err
		}
		for _, source := range ro.Sources {
			// By now, each remote path is a kcl module.
			// And, only one remote path is allowed.
			// So when more than one remote path or local path and remote path are provided, return an error.
			if (!source.IsLocalPath() || source.IsLocalTarPath()) && len(ro.Sources) > 1 {
				return errors.New("only one kcl module root path is allowed")
			}

			tmpRootPath, err := source.FindRootPath()
			if err != nil {
				return err
			}
			if tmpRootPath != rootPath {
				return errors.New("only one kcl module root path is allowed: root path conflicts between " + rootPath + " with " + tmpRootPath)
			}
		}
	}

	return nil
}

type RunOption func(*RunOptions) error

func WithKclOptions(opts kcl.Option) RunOption {
	return func(ro *RunOptions) error {
		ro.Option = &opts
		return nil
	}
}

func WithRunSources(sources []*downloader.Source) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Sources = sources
		return nil
	}
}

func WithSource(source *downloader.Source) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Sources = append(ro.Sources, source)
		return nil
	}
}

func WithEntries(entries []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if ro.Sources == nil {
			ro.Sources = make([]*downloader.Source, 0)
		}

		for _, entry := range entries {
			source, err := downloader.NewSourceFromStr(entry)
			if err != nil {
				return err
			}
			ro.Sources = append(ro.Sources, source)
		}

		return nil
	}
}

func WithSettingFiles(settingFiles []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.hasSettingsYaml = true
		ro.SettingYamls = settingFiles
		return nil
	}
}

func WithArguments(args []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithOptions(args...))
		return nil
	}
}

func WithOverrides(overrides []string, debug bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithOverrides(overrides...))
		ro.PrintOverrideAst = debug

		return nil
	}
}

func WithPathSelectors(pathSelectors []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithSelectors(pathSelectors...))
		return nil
	}
}

func WithDebug(debug bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if debug {
			ro.Debug = 1
		}

		return nil
	}
}

func WithDisableNone(disableNone bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if disableNone {
			ro.Merge(kcl.WithDisableNone(disableNone))
		}

		return nil
	}
}

func WithExternalPkgs(externalPkgs []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithExternalPkgs(externalPkgs...))
		return nil
	}
}

func WithSortKeys(sortKeys bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if sortKeys {
			ro.Merge(kcl.WithSortKeys(sortKeys))
		}

		return nil
	}
}

func WithShowHidden(showHidden bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if showHidden {
			ro.Merge(kcl.WithShowHidden(showHidden))
		}

		return nil
	}
}

func WithStrictRange(strictRange bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if strictRange {
			ro.StrictRangeCheck = strictRange
		}

		return nil
	}
}

func WithCompileOnly(compileOnly bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if compileOnly {
			ro.CompileOnly = compileOnly
		}

		return nil
	}
}

func WithVendor(vendor bool) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		if vendor {
			ro.Vendor = vendor
		}

		return nil
	}
}

func WithWorkDir(workDir string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.workDir = workDir
		return nil
	}
}

func (o *RunOptions) BaseEntry() (*downloader.Source, error) {
	// 如果参数中没有任何需要编译的东西，那么，当前目录就是 base entry, 直接应用包内情况
	if o.Sources == nil || len(o.Sources) == 0 {
		if o.workDir == "" {
			pwd, err := os.Getwd()
			if err != nil {
				return nil, err
			}

			o.workDir = pwd
		}

		return downloader.NewSourceFromStr(o.workDir)
	} else {
		return o.Sources[0], nil
	}
}

func (o *RunOptions) NoCompileEntries() bool {
	return o.Sources == nil || len(o.Sources) == 0
}

// yaml 中文件路径需要从相对路径替换为绝对路径
func (o *RunOptions) loadYamlSettingsFromLocalAndCli(rootPath string) error {
	err := ErrNoYamlSettings
	if o.hasSettingsYaml {
		for _, setting := range o.SettingYamls {
			o.Merge(kcl.WithSettings(setting))
			err = nil
		}
	} else {
		localSettingsYamlPath := filepath.Join(o.workDir, constants.KCL_YAML)
		if utils.DirExists(localSettingsYamlPath) {
			o.Merge(kcl.WithSettings(localSettingsYamlPath))
			o.hasSettingsYaml = true
			err = nil
		}
	}

	for i, kfile := range o.KFilenameList {
		if !filepath.IsAbs(kfile) {
			o.KFilenameList[i] = filepath.Join(rootPath, kfile)
		}
	}

	return err
}

var ErrNoCliSettings = errors.New("no cli settings")
var ErrNoYamlSettings = errors.New("no yaml settings")

func (o *RunOptions) loadCliSettings(rootPath string, baseEntry *downloader.Source) error {
	if o.NoCompileEntries() {
		return ErrNoCliSettings
	} else {
		if o.Sources[0].IsPackaged() {
			o.Sources = o.Sources[1:]
		}
		if len(o.Sources) == 0 {
			return ErrNoCliSettings
		}
	}

	baseRootPath, err := baseEntry.FindRootPath()
	if err != nil {
		return err
	}

	for _, source := range o.Sources {
		sourceRootPath, err := source.FindRootPath()
		if err != nil {
			return err
		}

		if baseEntry.IsPackaged() && source.IsPackaged() {
			sourceStr, err := source.ToString()
			if err != nil {
				return err
			}
			return reporter.NewErrorEvent(
				reporter.CompileFailed,
				fmt.Errorf("cannot compile multiple packages %s at the same time", []string{rootPath, sourceStr}),
				"only allows one package to be compiled at a time",
			)
		} else if !baseEntry.IsPackaged() && baseRootPath != sourceRootPath {
			return reporter.NewErrorEvent(
				reporter.CompileFailed,
				fmt.Errorf("cannot compile multiple packages %s at the same time", []string{baseRootPath, sourceRootPath}),
				"only allows one package to be compiled at a time",
			)
		} else {
			// 如果 baseEntry 是非 pwd 的本地路径，相对计算时使用 Workdir 计算的
			// 如果 baseEntry 是 pwd,git,oci,tag，那么，相对路径就是相对于包的homepath
			if filepath.IsAbs(source.Local.Path) {
				o.Merge(kcl.WithKFilenames(source.Local.Path))
			} else {
				o.Merge(kcl.WithKFilenames(filepath.Join(rootPath, source.Local.Path)))
			}
		}
	}
	return nil

}

func (o *RunOptions) loadCompileSettings(baseEntry *downloader.Source, basePkg *pkg.KclPkg) error {
	var rootPath string
	if !baseEntry.IsPackaged() {
		rootPath = o.workDir
	} else {
		rootPath = basePkg.HomePath
	}

	err := o.loadCliSettings(rootPath, baseEntry)
	if err != nil && err != ErrNoCliSettings {
		return err
	}

	if err == ErrNoCliSettings {
		err = o.loadYamlSettingsFromLocalAndCli(rootPath)
		if err != nil {
			if err != ErrNoYamlSettings {
				return err
			}
			o.Merge(*basePkg.GetKclOpts())
			// 如果最后没有任何要被编译的入口
			if len(o.KFilenameList) == 0 {
				o.Merge(kcl.WithKFilenames(basePkg.HomePath))
			} else {
				for i, kfile := range o.KFilenameList {
					if !filepath.IsAbs(kfile) {
						o.KFilenameList[i] = filepath.Join(basePkg.HomePath, kfile)
					}
				}
			}
		}
	}

	return nil
}

// 需要一个 source accesser 来访问 source
func (c *KpmClient) Run(options ...RunOption) (*kcl.KCLResultList, error) {
	o := &RunOptions{}
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}

	// acquire the lock of the package cache.
	err := c.AcquirePackageCacheLock()
	if err != nil {
		return nil, err
	}

	defer func() {
		// release the lock of the package cache after the function returns.
		releaseErr := c.ReleasePackageCacheLock()
		if releaseErr != nil && err == nil {
			err = releaseErr
		}
	}()

	// 1. 首先，需要一个 base entry 是后续 kcl.yaml 和其他 main.k 生效的地方
	baseEntry, err := o.BaseEntry()
	if err != nil {
		return nil, err
	}

	var pkgMap map[string]string
	var res *kcl.KCLResultList
	err = NewVisitor(baseEntry, c).Visit(baseEntry, func(basePkg *pkg.KclPkg) error {
		err = o.loadCompileSettings(baseEntry, basePkg)
		if err != nil {
			return err
		}

		pkgMap, err = c.ResolveDepsIntoMap(basePkg)
		if err != nil {
			return err
		}

		// Fill the dependency path.
		for dName, dPath := range pkgMap {
			if !filepath.IsAbs(dPath) {
				dPath = filepath.Join(c.homePath, dPath)
			}
			o.Merge(kcl.WithExternalPkgs(fmt.Sprintf(runner.EXTERNAL_PKGS_ARG_PATTERN, dName, dPath)))
		}

		res, err = kcl.RunWithOpts(*o.Option)

		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
