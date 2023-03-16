package pkg

import "strings"

func (dep *Dependency) Eq(other *Dependency) bool {
	return strings.EqualFold(dep.Name, other.Name)
}

func (dep *Dependency) Download() (string, error) {
	if d, ok := dep.GetDependency().(*Dependency_Git); ok {
		return d.Git.Download()
	}
	return "", nil
}

func (dep *GitDependency) Download() (string, error) {
	
}
