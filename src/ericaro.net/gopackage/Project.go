package gopackage

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
a Project is an in memory representation of a real go project on the disk.
It is serialized on the local project disk
*/
type Project struct {
	Root            string  `json:"-"` // absolute path of the root, other pathes are relative 
	Name string   // the identity  of the project
	Version         *Version //optional, only for deployed instances 
	Snapshot        *bool    // optional 
	//	Target       string   // path to the target dir where we can generate stuff
	Dependencies []ProjectReference
}

//NewProject creates a new Project object with default values.
func NewProject() *Project {
	return &Project{
		Name:        "",
		Dependencies: make([]ProjectReference, 0),
	}
}

//ReadProject local info from the current dir
func ReadProject() (p *Project, err error) {
	return ReadProjectFile(GopackageFile)
}

func (prj *Project) Reference() (p ProjectReference) {
	return NewProjectReference(prj.Name, prj.Version.Reference())
}

//ReadProjectFile local info from the specified gopackage file
func ReadProjectFile(gopackagepath string) (p *Project, err error) {
	p = NewProject()
	f, err := os.Open(gopackagepath)
	if err != nil {
		return
	}
	defer f.Close()

	err = Decode(p, f)

	p.Root, _ = filepath.Abs(path.Dir(gopackagepath))
	return
}

//ReadProjectTar reads the .gopackage file within the tar in memory. It does not set the Root
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
				err = errors.New(fmt.Sprintf("Invalid package format, %v is missing", GopackageFile))
			}
		}
		if hdr.Name == GopackageFile {
			p = NewProject()
			err = Decode(p, tr)
			break
		}
	}
	return
}

//Untar reads the .gopackage file within the tar in memory. It does not set the Root
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

//TODO there are too many ways to write a project. And in particular, it gets sometimes
// with unwanted values (like the local path etc.) analyse the calling track and split the
// use cases

//WriteProjectSrc into a source project
func WriteProjectSrc(p *Project) error {
	return writeProject( *p) 
}
//WriteProjectPkg into a packaged project
func WriteProjectPkg(p *Project) error {
	return writeProject(*p) 
}
func writeProject(p Project) (err error) {
	dst := filepath.Join(p.Root, GopackageFile)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = EncodeProject(&p, f)
	return err
}
func EncodeProject(p *Project, w io.Writer) (err error) {
	err = Encode(p, w)
	return
}

func (p *Project) AppendDependency(ref ...ProjectReference) {
	p.Dependencies = append(p.Dependencies, ref...)
}

func (p *Project) RemoveDependency(ref ProjectReference) {
	src := p.Dependencies
	is := make([]int, 0, len(src))
	for i, r := range src {
		if ref.Equals(r) {
			is = append(is, i)
		}
	}
	if len(is) == 0 { // nothing to do
		return
	}
	dep := make([]ProjectReference, 0, len(src)-len(is))
	length := len(is)
	if is[0] > 0 {
		dep = append(dep, src[0:is[0]]...)
	}
	for j := 0; j < length-1; j++ {
		s, e := is[j]+1, is[j+1]
		dep = append(dep, src[s:e]...)
	}
	// last bit of slice
	p.Dependencies = dep
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
	TarBuff(filepath.Join("/", GopackageFile), buf, tw)

	return
}

func (p *Project) String() string {
	dependencies := ""
	for _, dep := range p.Dependencies {
		dependencies += fmt.Sprintf("\n  -> %v", dep)
	}
	return fmt.Sprintf("%v:%v\n", p.Name, dependencies)
}
