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

type RunOption func(*RunOptions) error

// WithKclOptions sets the kcl options for the kcl compiler.
func WithKclOptions(opts kcl.Option) RunOption {
	return func(ro *RunOptions) error {
		ro.Option = &opts
		return nil
	}
}

// WithRunSources sets the sources for the kcl package.
func WithRunSources(sources []*downloader.Source) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Sources = sources
		return nil
	}
}

// WithSource sets the source for the kcl package.
func WithSource(source *downloader.Source) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Sources = append(ro.Sources, source)
		return nil
	}
}

// WithEntries sets the entries for the kcl package.
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

// WithSettingsFiles sets the settings files for the kcl package.
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

// WithArguments sets the arguments for the kcl package.
func WithArguments(args []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithOptions(args...))
		return nil
	}
}

// WithOverrides sets the overrides for the kcl package.
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

// WithPathSelectors sets the path selectors for the kcl package.
func WithPathSelectors(pathSelectors []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithSelectors(pathSelectors...))
		return nil
	}
}

// WithDebug sets the debug flag for the kcl package.
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

// WithDisableNone sets the disable none flag for the kcl package.
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

// WithExternalPkgs sets the external packages for the kcl package.
func WithExternalPkgs(externalPkgs []string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.Merge(kcl.WithExternalPkgs(externalPkgs...))
		return nil
	}
}

// WithSortKeys sets the sort keys flag for the kcl package.
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

// WithShowHidden sets the show hidden flag for the kcl package.
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

// WithStrictRange sets the strict range flag for the kcl package.
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

// WithCompileOnly sets the compile only flag for the kcl package.
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

// WithVendor sets the vendor flag for the kcl package.
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

// WithWorkDir sets the work directory for the kcl package.
func WithWorkDir(workDir string) RunOption {
	return func(ro *RunOptions) error {
		if ro.Option == nil {
			ro.Option = kcl.NewOption()
		}
		ro.workDir = workDir
		return nil
	}
}

// RootPkgSource returns the root package source.
func (o *RunOptions) RootPkgSource() (*downloader.Source, error) {
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

// NoCompileEntries returns true if there is no compile entries.
func (o *RunOptions) NoCompileEntries() bool {
	return o.Sources == nil || len(o.Sources) == 0
}

// Merge merges the options from kcl.yaml which is located in the work directory or inputed by the cli.
func (o *RunOptions) loadYamlSettingsFromLocalAndCli(rootPath string) error {
	err := ErrNoYamlSettings
	// If the settings are inputed by the cli
	if o.hasSettingsYaml {
		for _, setting := range o.SettingYamls {
			o.Merge(kcl.WithSettings(setting))
			err = nil
		}
	} else {
		// if the settings are in the work directory.
		localSettingsYamlPath := filepath.Join(o.workDir, constants.KCL_YAML)
		if utils.DirExists(localSettingsYamlPath) {
			o.Merge(kcl.WithSettings(localSettingsYamlPath))
			o.hasSettingsYaml = true
			err = nil
		}
	}

	// Iterate all the *.k files get from kcl.yaml and make them absolute path.
	// If the compilation entries are all local *.k files or directories, the root path is the work directory.
	// if the compilation entries are git repo, oci registry or tar files, the root path is the home path of the package.
	for i, kfile := range o.KFilenameList {
		if !filepath.IsAbs(kfile) {
			o.KFilenameList[i] = filepath.Join(rootPath, kfile)
		}
	}

	return err
}

var ErrNoCliSettings = errors.New("no cli settings")
var ErrNoYamlSettings = errors.New("no yaml settings")

// loadCliSettings loads the settings from the cli.
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
			if filepath.IsAbs(source.Local.Path) {
				o.Merge(kcl.WithKFilenames(source.Local.Path))
			} else {
				o.Merge(kcl.WithKFilenames(filepath.Join(rootPath, source.Local.Path)))
			}
		}
	}
	return nil
}

// loadCompileSettings loads the compile settings from the kcl.yaml and cli.
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

	rootPkgSource, err := o.RootPkgSource()
	if err != nil {
		return nil, err
	}

	var pkgMap map[string]string
	var res *kcl.KCLResultList
	err = NewVisitor(rootPkgSource, c).Visit(rootPkgSource, func(basePkg *pkg.KclPkg) error {
		err = o.loadCompileSettings(rootPkgSource, basePkg)
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
