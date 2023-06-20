package runner

import (
	"fmt"
	"strings"

	"kusionstack.io/kclvm-go/pkg/kcl"
)

// The pattern of the external package argument.
const EXTERNAL_PKGS_ARG_PATTERN = "%s=%s"

// Compiler is a wrapper of kcl compiler.
type Compiler struct {
	kclCliArgs string
	kclOpts    *kcl.Option
}

// DefaultCompiler will create a default compiler.
func DefaultCompiler() *Compiler {
	return &Compiler{
		kclOpts: kcl.NewOption(),
	}
}

// SetKclCliArgs will set the kcl cli args.
func (compiler *Compiler) SetKclCliArgs(kclCliArgs string) *Compiler {
	compiler.kclCliArgs = kclCliArgs
	return compiler
}

// AddKFile will add a file path to the entry file list.
func (compiler *Compiler) AddKFile(entryKFile string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithKFilenames(entryKFile))
	return compiler
}

// AddDep will add a file path to the dependency list.
func (compiler *Compiler) AddDepPath(depName string, depPath string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithExternalPkgs(fmt.Sprintf(EXTERNAL_PKGS_ARG_PATTERN, depName, depPath)))
	return compiler
}

// AddWorkDir will add a directory path to the work directory list.
func (compiler *Compiler) AddWorkDir(workDir string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithWorkDir(workDir))
	return compiler
}

// AddSettings will add a settings file path to the settings file list.
func (compiler *Compiler) AddSettings(settingsFile string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithSettings(settingsFile))
	return compiler
}

// SetOptions will set the kcl options.
func (compiler *Compiler) AddOptions(opt string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithOptions(opt))
	return compiler
}

// SetOverride will set the override option.
func (compiler *Compiler) AddOverrides(override string) *Compiler {
	compiler.kclOpts.Merge(kcl.WithOverrides(override))
	return compiler
}

// SetDisableNone will set the disable none option.
func (compiler *Compiler) SetDisableNone(disableNone bool) *Compiler {
	compiler.kclOpts.Merge(kcl.WithDisableNone(disableNone))
	return compiler
}

// SetSortKeys will set the sort keys option.
func (compiler *Compiler) SetSortKeys(sortKeys bool) *Compiler {
	compiler.kclOpts.Merge(kcl.WithSortKeys(sortKeys))
	return compiler
}

// Call KCL Compiler and return the result.
func (compiler *Compiler) Run() (*kcl.KCLResultList, error) {
	// Parse all the kcl options.
	kclFlags, err := ParseArgs(strings.Fields(compiler.kclCliArgs))
	if err != nil {
		return nil, err
	}

	// Merge the kcl options from kcl.mod and kpm cli.
	compiler.kclOpts.Merge(kclFlags.IntoKclOptions())

	return kcl.RunWithOpts(*compiler.kclOpts)
}
