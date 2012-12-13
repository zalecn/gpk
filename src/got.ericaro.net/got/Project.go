package got

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

/*
a Project is an in memomry representation of a real go project on the disk.
*/
type Project struct {
	Root            string   // absolute path of the root, other pathes are relative 
	Group, Artifact string   // the identity  of the project
	Version         *Version //optional, only for deployed instances 
	Snapshot        *bool    // optional 
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
	return ReadProjectFile(GotFile)
}


func (prj *Project) Reference() (p ProjectReference){
	return NewProjectReference(prj.Group, prj.Artifact, prj.Version.Reference() )
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

	p.Root, _ = filepath.Abs(path.Dir(gotpath))
	return
}

//ReadProjectTar reads the .got file within the tar in memory. It does not set the Root
func ReadProjectInPackage(in io.Reader) (p *Project, err error) {
	gz, err := gzip.NewReader(in)
	if err != nil {
		return
	}
	tr := tar.NewReader(in)
	defer gz.Close()
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				err = errors.New(fmt.Sprintf("Invalid package format, %v is missing", GotFile))
			}
		}
		if hdr.Name == GotFile {
			p = NewProject()
			err = Decode(p, tr)
			break
		}
	}
	return
}

//Untar reads the .got file within the tar in memory. It does not set the Root
func (p *Project) UnpackageProject(in io.Reader) (err error) {
	gz, err := gzip.NewReader(in)
	tr := tar.NewReader(in)
	defer gz.Close()
	dst := p.Root
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			break
		}
		// make the target file
		ndst := filepath.Join(dst, hdr.Name)
		df, err := os.Create(ndst)
		if err != nil {
			break
		}
		io.Copy(df, tr)
		df.Close()
	}
	return
}

//WriteProject local info from the current dir
func WriteProject(p *Project) error {
	return WriteProjectFile(GotFile, p)
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
func (p Project) PackageProject(in io.Writer) (err error) {
	p.Root = "" // delete the root before package, Note that we are working on a "copy" of the project
	gz, err := gzip.NewWriterLevel(in, gzip.BestCompression)
	if err != nil {
		return
	}
	defer gz.Close()
	
	tw := tar.NewWriter(gz)
	
	//prepare recursive handlers
	dirHandler := func(ldst, lsrc string) (err error) {
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		err = TarFile(ldst, lsrc, tw)
		return
	}
	walkDir("/", filepath.Join(p.Root, "src"), dirHandler, fileHandler)

	buf := new(bytes.Buffer)
	EncodeProject(&p, buf)
	TarBuff(filepath.Join("/", GotFile), buf, tw)
	
	return
}

func (p *Project) String() string {
	dependencies := ""
	for _, dep := range p.Dependencies {
		dependencies += fmt.Sprintf("\n  -> %v", dep)
	}
	return fmt.Sprintf("%v:%v%v\n", p.Group, p.Artifact, dependencies)
}
