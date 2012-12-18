package gpk

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ProjectID struct {
	name    string // any valid package name
	version Version
}

type Project struct {
	workingDir   string      // transient workding directory aboslute path
	name         string      // package name
	dependencies []ProjectID // contains the current project's dependencies
	// TO be added build time , and test dependencies
}
type Package struct {
	self      Project
	version   Version
	timestamp time.Time

	// more to come, like sha1,signature, snapshot/release
	// add also go1 , i.e the target go runtime.

}

func NewProjectID(name string, version Version) ProjectID {
	return ProjectID{name: name, version: version}
}

func ReadProject() (p *Project, err error) {
	return ReadProjectFile(GpkFile) // read from the current dir. TODO look up in the hierarchy too
}

//ReadPackageFile local info from the specified gopackage file
func ReadPackageFile(gpkPath string) (p *Package, err error) {
	p = &Package{}
	f, err := os.Open(gpkPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = DecodePackage(f, p)
	p.self.workingDir, _ = filepath.Abs(path.Dir(gpkPath))
	return
}

//ReadProjectFile local info from the specified gopackage file
func ReadProjectFile(gpkPath string) (p *Project, err error) {
	p = &Project{}
	f, err := os.Open(gpkPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = DecodeProject(f, p)
	p.workingDir, _ = filepath.Abs(path.Dir(gpkPath))
	return
}

func (p *Package) Timestamp() time.Time {
	return p.timestamp
}

//Write package  info to where it belongs (package holds working dir info)
func (p Package) Write() (err error) {
	dst := filepath.Join(p.self.workingDir, GpkFile)
	fmt.Printf("writing package to %s\n", dst)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = EncodePackage(f, p)
	return err
}

func (p *Project) WorkingDir() string {
	return p.workingDir
}
func (p *Package) InstallDir() string {
	return p.self.workingDir
}
func (p *Project) Name() string {
	return p.name
}
func (p *Project) SetWorkingDir(pwd string) {
	p.workingDir = pwd
}
func (p *Project) SetName(name string) {
	p.name = name
}

func (p *Project) AppendDependency(ref ...ProjectID) {
	p.dependencies = append(p.dependencies, ref...)
}

func (p *Project) RemoveDependency(ref ProjectID) {
	src := p.dependencies
	is := make([]int, 0, len(src))
	for i, r := range src {
		if ref.Equals(r) {
			is = append(is, i)
		}
	}
	if len(is) == 0 { // nothing to do
		return
	}
	dep := make([]ProjectID, 0, len(src)-len(is))
	length := len(is)
	if is[0] > 0 {
		dep = append(dep, src[0:is[0]]...)
	}
	for j := 0; j < length-1; j++ {
		s, e := is[j]+1, is[j+1]
		dep = append(dep, src[s:e]...)
	}
	// last bit of slice
	p.dependencies = dep
}

//Write project  info to where it belongs (project holds working dir info)
func (p Project) Write() (err error) {
	dst := filepath.Join(p.workingDir, GpkFile)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = EncodeProject(f,p)
	return err
}

//Path converts this project reference into the path it should have in the repository layout
func (d ProjectID) Path() string {
	return filepath.Join(d.name, d.version.String())
}

func (this ProjectID) Equals(that ProjectID) bool {
	return this.ID() == that.ID()
}

func (d ProjectID) ID() string {
	return fmt.Sprintf("%s:%s", d.name, d.version.String())
}
func (d ProjectID) String() string {
	return d.ID()
}

//Path converts this project reference into the path it should have in the repository layout
func (p *Package) Path() string {
	return filepath.Join(p.self.name, p.version.String())
}

func (p *Package) ID() ProjectID {
	return ProjectID{
		name:    p.self.name,
		version: p.version,
	}
}

//ReadProjectTar reads the .gopackage file within the tar in memory. It does not set the Root
func ReadPackageInPackage(in io.Reader) (p *Package, err error) {
	//fmt.Printf("Parsing in memory package\n")
	gz, err := gzip.NewReader(in)
	if err != nil {
		return
	}
	tr := tar.NewReader(gz)
	defer gz.Close()
	
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				err = errors.New(fmt.Sprintf("Invalid package format, %v is missing", GpkFile))
			}
			break
		}
		//fmt.Printf("hdr %v\n", hdr )
		if hdr.Name == GpkFile {
			p = &Package{}
			err = DecodePackage(tr, p)
			return p, err
		}
	}
	return
}

//Untar reads the .gopackage file within the tar in memory. It does not set the Root
func (p *Package) Unpack(in io.Reader) (err error) {
	gz, err := gzip.NewReader(in)
	if err != nil {return}
	
	defer gz.Close()
	tr := tar.NewReader(gz)
	dst := p.self.workingDir
	fmt.Printf("unpacking to %s\n", dst)
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("error in tar %s\n", err)
			return err
		}
		// make the target file
		ndst := filepath.Join(dst, hdr.Name)
		os.MkdirAll(path.Dir(ndst), os.ModeDir|os.ModePerm) // mkdir -p
		//fmt.Printf("%s\n", ndst)
		df, err := os.Create(ndst)
		if err != nil {
		fmt.Printf("error unpacking %v\n", err)
			break
		}
		io.Copy(df, tr)
		df.Close()
	}
	return
}

//PackageProject into a tar writer
func (p *Package) Pack(in io.Writer) (err error) {
	gz, err := gzip.NewWriterLevel(in, gzip.BestCompression)
	if err != nil {
		return
	}

	tw := tar.NewWriter(gz)
	defer gz.Close()
	defer tw.Close()

	//prepare recursive handlers
	dirHandler := func(ldst, lsrc string) (err error) {
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		err = TarFile(ldst, lsrc, tw)
		return
	}
	walkDir("src", filepath.Join(p.self.workingDir, "src"), dirHandler, fileHandler)
	// copy the package .gpk
	TarFile(filepath.Join("", GpkFile), filepath.Join(p.self.workingDir, GpkFile), tw)
	// or rewrite it (and edit it on the fly ?
	//	buf := new(bytes.Buffer)
	//	json.NewEncoder(buf).Encode(p)
	//	TarBuff(filepath.Join("/", GpkFile), buf, tw)
	
	return
}

type projectFile struct {
	Name string
	Dependencies []projectIDFile
}
type packageFile struct {
	Self      projectFile
	Version   string
	Timestamp time.Time
}

type projectIDFile struct {
	Name, Version string
}


func DecodePackage(r io.Reader, p *Package) (err error) {
	pf := new(packageFile)
	err = json.NewDecoder(r).Decode(pf)
	if err != nil {
		return
	}
	decodePackageFile(*pf, p)
	return
}
func EncodePackage(w io.Writer, p Package) (err error) {
	pf := encodePackageFile(p)
	err = json.NewEncoder(w).Encode(pf)
	return 
	
}



func EncodeProject(w io.Writer, p Project) (err error) {
	pf := encodeProjectFile(p)
	err = json.NewEncoder(w).Encode(pf)
	return 
	
}

func DecodeProject(r io.Reader, p *Project) (err error) {
	pf := new(projectFile)
	err = json.NewDecoder(r).Decode(pf)
	if err != nil {
		return
	}
	decodeProjectFile(*pf, p)
	return
}

func decodeProjectFile(pf projectFile, p *Project)  {
	p.name = pf.Name
	for _,d := range pf.Dependencies {
		v,_ := ParseVersion(d.Version)
		p.AppendDependency( NewProjectID( d.Name, v ) )
	}
}
func decodePackageFile(pf packageFile, p *Package)  {
	prj := new(Project)
	decodeProjectFile(pf.Self, prj)
	p.self = *prj
	p.timestamp = pf.Timestamp
	v,_ := ParseVersion(pf.Version)
	p.version = v
}


func encodePackageFile(p Package) *packageFile{
	return &packageFile{
		Self : *encodeProjectFile(p.self),
		Timestamp: p.timestamp,
		Version: p.version.String(),
	}
	
}

func encodeProjectIDFile(p ProjectID) *projectIDFile {
return &projectIDFile{
	Name: p.name,
	Version: p.version.String(),
}
}
func encodeProjectFile(p Project) *projectFile {
	dep := make([]projectIDFile,0, len(p.dependencies))
	
	for _,d:= range p.dependencies{
		dep = append(dep, *encodeProjectIDFile(d))
	}
	
return &projectFile{
		Name: p.name,
		Dependencies: dep,
	}
}
