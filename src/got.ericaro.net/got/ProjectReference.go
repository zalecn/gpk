package got

import (
	"strings"
	"path/filepath"
	"fmt"
)

// ProjectReference is just a way to keep references to another project
type ProjectReference struct {
	Group, Artifact string // the symbolic name for this project.
	Version         VersionReference
}

func ParseProjectReference(value string) ProjectReference {
	p := ProjectReference{}
	parts := strings.Split(value, ":")
	p.Group = parts[0]
	p.Artifact = parts[1]
	p.Version = ParseVersionReference(parts[2])
	return p
}

func NewProjectReference(group, artifact string, version VersionReference) ProjectReference {
	return ProjectReference{
	Group: group,
	Artifact: artifact,
	Version: version,
	}
}

//Path converts this project reference into the path it should have in the repository layout
func (d ProjectReference) Path() string {
	return filepath.Join(d.Group,d.Artifact, d.Version.Path())
}

func (d ProjectReference) String() string {
	return fmt.Sprintf("%v:%v:%v", d.Group, d.Artifact,  d.Version)
}