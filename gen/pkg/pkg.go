package pkg

import "strings"

func (dep *Dependency) Eq(other *Dependency) bool {
	return strings.EqualFold(dep.Name, other.Name)
}
