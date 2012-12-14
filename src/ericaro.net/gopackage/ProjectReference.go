package gopackage

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// ProjectReference is just a way to keep references to another project
type ProjectReference struct {
	Name string // the symbolic name for this project.
	Version         VersionReference
}

func ParseProjectReference(value string) (p ProjectReference, err error) {
	p = ProjectReference{}
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		errors.New("Invalid Project Reference Format")
	}
	p.Name = parts[0]
	p.Version = ParseVersionReference(parts[1])
	return
}

func NewProjectReference(name string, version VersionReference) ProjectReference {
	return ProjectReference{
		Name:    name,
		Version:  version,
	}
}

func (this ProjectReference) Equals(that ProjectReference) bool {
	return this.String() == that.String()
}

func (p *ProjectReference) Project() (prj *Project) {
	return &Project{
		Name:    p.Name,
		Version:  p.Version.Version(),
	}
}

func NewGoGetProjectReference(pack string, version VersionReference) ProjectReference {
	parts := strings.SplitN(pack, "/", 2)
	if len(parts) != 2 {
		panic("Not a valid go get package " + pack)
	}
	return ProjectReference{
		Name:    pack,
		Version:  version,
	}
}

//Path converts this project reference into the path it should have in the repository layout
func (d ProjectReference) Path() string {
	return filepath.Join(d.Name, d.Version.Path())
}

func (d ProjectReference) String() string {
	return fmt.Sprintf("%v:%v", d.Name, d.Version)
}
