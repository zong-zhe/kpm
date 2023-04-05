package modfile

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type TOML interface {
	MarshalTOML() string
	UnmarshalTOML(data interface{}) error
}

const NEWLINE = "\n"

func (mod *ModFile) MarshalTOML() string {
	var sb strings.Builder
	sb.WriteString(mod.Pkg.MarshalTOML())
	sb.WriteString(NEWLINE)
	sb.WriteString(mod.Dependencies.MarshalTOML())
	return sb.String()
}

const PACKAGE_PATTERN = "[package]"

func (pkg *Package) MarshalTOML() string {
	var sb strings.Builder
	sb.WriteString(PACKAGE_PATTERN)
	sb.WriteString(NEWLINE)
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(pkg); err != nil {
		fmt.Println(err)
		return ""
	}
	sb.WriteString(buf.String())
	return sb.String()
}

const DEPS_PATTERN = "[dependencies]"

func (dep *Dependencies) MarshalTOML() string {
	var sb strings.Builder
	sb.WriteString(DEPS_PATTERN)
	for _, dep := range dep.Deps {
		sb.WriteString(NEWLINE)
		sb.WriteString(dep.MarshalTOML())
	}
	return sb.String()
}

const DEP_PATTERN = "%s = %s"

func (dep *Dependency) MarshalTOML() string {
	source := dep.Source.MarshalTOML()
	var sb strings.Builder
	if len(source) != 0 {
		sb.WriteString(fmt.Sprintf(DEP_PATTERN, dep.Name, source))
	}
	return sb.String()
}

const SOURCE_PATTERN = "{ %s }"

func (source *Source) MarshalTOML() string {
	gitToml := source.Git.MarshalTOML()
	var sb strings.Builder
	if len(gitToml) != 0 {
		sb.WriteString(fmt.Sprintf(SOURCE_PATTERN, gitToml))
	}
	return sb.String()
}

const GTI_URL_PATTERN = "git = \"%s\""
const GTI_TAG_PATTERN = "tag = \"%s\""
const SEPARATOR = ", "

func (git *Git) MarshalTOML() string {
	var sb strings.Builder
	if len(git.Url) != 0 {
		sb.WriteString(fmt.Sprintf(GTI_URL_PATTERN, git.Url))
	}
	if len(git.Tag) != 0 {
		sb.WriteString(SEPARATOR)
		sb.WriteString(fmt.Sprintf(GTI_TAG_PATTERN, git.Tag))
	}
	return sb.String()
}

const PACKAGE_FLAG = "package"
const DEPS_FLAG = "dependencies"

func (mod *ModFile) UnmarshalTOML(data interface{}) error {
	meta, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", data)
	}

	if v, ok := meta[PACKAGE_FLAG]; ok {
		pkg := Package{}
		pkg.UnmarshalTOML(v)
		mod.Pkg = pkg
	}

	if v, ok := meta[DEPS_FLAG]; ok {
		deps := Dependencies{
			Deps: make(map[string]Dependency),
		}
		deps.UnmarshalTOML(v)
		mod.Dependencies = deps
	}

	return nil
}

const NAME_FLAG = "name"
const EDITION_FLAG = "edition"
const VERSION_FLAG = "version"

func (pkg *Package) UnmarshalTOML(data interface{}) error {
	meta, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", data)
	}

	if v, ok := meta[NAME_FLAG].(string); ok {
		pkg.Name = v
	}

	if v, ok := meta[EDITION_FLAG].(string); ok {
		pkg.Edition = v
	}

	if v, ok := meta[VERSION_FLAG].(string); ok {
		pkg.Version = v
	}
	return nil
}

func (deps *Dependencies) UnmarshalTOML(data interface{}) error {
	meta, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", data)
	}

	for k, v := range meta {
		dep := Dependency{}
		dep.UnmarshalTOML(v)
		dep.Name = k
		deps.Deps[k] = dep
	}

	return nil
}

func (dep *Dependency) UnmarshalTOML(data interface{}) error {
	source := Source{}
	err := source.UnmarshalTOML(data)
	if err != nil {
		return err
	}
	dep.Source = source
	return nil
}

func (source *Source) UnmarshalTOML(data interface{}) error {
	git := Git{}
	err := git.UnmarshalTOML(data)
	if err != nil {
		return err
	}
	source.Git = &git
	return nil
}

const GTI_URL_FLAG = "git"
const GTI_TAG_FLAG = "tag"

func (git *Git) UnmarshalTOML(data interface{}) error {
	meta, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", data)
	}

	if v, ok := meta[GTI_URL_FLAG].(string); ok {
		git.Url = v
	}

	if v, ok := meta[GTI_TAG_FLAG].(string); ok {
		git.Tag = v
	}

	return nil
}
