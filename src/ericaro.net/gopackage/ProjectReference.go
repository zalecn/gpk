package gopackage

import (
	"strings"
	"path/filepath"
	"fmt"
	"errors"
)

// ProjectReference is just a way to keep references to another project
type ProjectReference struct {
	Group, Artifact string // the symbolic name for this project.
	Version         VersionReference
}

func ParseProjectReference(value string) (p ProjectReference, err error) {
	p = ProjectReference{}
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		errors.New("Invalid Project Reference Format") 
	}
	p.Group = parts[0]
	p.Artifact = parts[1]
	p.Version = ParseVersionReference(parts[2])
	return
}

func NewProjectReference(group, artifact string, version VersionReference) ProjectReference {
	return ProjectReference{
	Group: group,
	Artifact: artifact,
	Version: version,
	}
}

func (this ProjectReference) Equals(that ProjectReference) bool {
	return this.String() == that.String()
}

func (p *ProjectReference) Project() (prj *Project) {
	return &Project{
	Group: p.Group,
	Artifact: p.Artifact,
	Version: p.Version.Version(),
	}
}

func NewGoGetProjectReference(pack string, version VersionReference) ProjectReference {
	parts := strings.SplitN(pack, "/", 2)
	if len(parts) !=2 {
		panic("Not a valid go get package "+pack)
	}
	lefties := strings.Split(parts[0], ".")
	righties := strings.Split(parts[1], "/")
	
	//reverse lefties
	for i,j:=0,len(lefties)-1; i< len(lefties)/2;i,j=i+1,j-1 {
		lefties[i], lefties[j] = lefties[j], lefties[i]
	}
	names:=make([]string, 0, len(lefties)+ len(righties) )
	names = append(names, lefties...) 
	names = append(names, righties...) 
	
	group := strings.Join(names[:len(names)-1], ".")
	artifact:= names[len(names)-1]
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