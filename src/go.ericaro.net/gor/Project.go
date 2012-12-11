package gor

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

/*
a Project is an in memomry representation of a real go project on the disk.
*/
type Project struct {
	Root string // absolute path of the root, other pathes are relative 
	Group, Artifact string // the identity  of the project
//	Target       string   // path to the target dir where we can generate stuff
	Dependencies []ProjectReference
}

//NewProject creates a new Project object with default values.
func NewProject() *Project {
	return &Project{
		Group:        "",
		Artifact:     "",
		Dependencies: make([]ProjectReference, 0),
	}
}

//ReadProject local info from the current dir
func ReadProject() (p *Project, err error) {
	return ReadProjectFile(GorFile) 
}
//ReadProjectFile local info from the specified gor file
func ReadProjectFile(gorpath string) (p *Project, err error) {
	p = NewProject()
	f, err := os.Open(gorpath)
	if err != nil {
		return
	}
	defer f.Close()

	err = Decode(p, f)
	
	p.Root, _ = filepath.Abs( path.Dir(gorpath) ) 
	return
}

//WriteProject local info from the current dir
func WriteProject(p *Project) error {
	return WriteProjectFile(GorFile, p)
}
func WriteProjectFile(file string, p *Project) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = Encode(p, f)
	return err
}

func (p *Project) AppendDependency(ref ...ProjectReference) {
	p.Dependencies = append(p.Dependencies, ref...)
}



func (p *Project) String() string {
	dependencies := ""
	for _,dep := range p.Dependencies {
		dependencies+=fmt.Sprintf("                        %v\n", dep)
	}
	return fmt.Sprintf(`
Project             : %v:%v
        Dependencies:
%v`, p.Group, p.Artifact, dependencies)
}
