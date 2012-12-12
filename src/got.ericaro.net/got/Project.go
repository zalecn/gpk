package got

import (
	"fmt"
	"os"
	"io"
	"path"
	"path/filepath"
	"archive/tar"
	"bytes"
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
//ReadProjectFile local info from the specified got file
func ReadProjectFile(gotpath string) (p *Project, err error) {
	p = NewProject()
	f, err := os.Open(gotpath)
	if err != nil {
		return
	}
	defer f.Close()

	err = Decode(p, f)
	
	p.Root, _ = filepath.Abs( path.Dir(gotpath) ) 
	return
}

//WriteProject local info from the current dir
func WriteProject(p *Project) error {
	return WriteProjectFile(GorFile, p)
}
func WriteProjectFile(file string, p *Project) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = EncodeProject(p, f)
	return err
}
func EncodeProject(p *Project, w io.Writer) (err error) {
	err = Encode(p, w)
	return
}



func (p *Project) AppendDependency(ref ...ProjectReference) {
	p.Dependencies = append(p.Dependencies, ref...)
}


//PackageProject into a tar writer
func (p *Project) PackageProject(tw *tar.Writer) {
		//prepare recursive handlers
	dirHandler := func(ldst, lsrc string)(err error){
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		err = TarFile(ldst, lsrc, tw)
		return
	}
	walkDir("/", filepath.Join(p.Root, "src") , dirHandler, fileHandler)
	
	buf := new(bytes.Buffer)
	EncodeProject(p, buf)
	TarBuff(filepath.Join("/", GorFile), buf, tw)
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
