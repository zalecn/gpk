package gopack

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	GpkFile = ".gpk"
)

type LicenseSet []License
type License struct {
	FullName, Alias string
}
var (
	Licenses LicenseSet = ([]License{
	License{ "Apache License 2.0"          , "ASF" }, 
	License{ "Eclipse Public License 1.0"  , "EPL" },
	License{ "GNU GPL v2"                  , "GPL2"},
	License{ "GNU GPL v3"                  , "GPL3"},
	License{ "GNU Lesser GPL"              , "LGPL"},
	License{ "MIT License"                 , "MIT" },
	License{ "Mozilla Public License 1.1"  , "MPL" },
	License{ "New BSD License"             , "BSD" },
	License{ "Other Open Source"           , "OOS" },
	License{ "Other Closed Source"         , "OCS" },
	})
)

var (
	CentralUrl = url.URL{
		Scheme: "http",
		Host:   "gpk.ericaro.net",
	}
	Central RemoteRepository
)


func (set LicenseSet) String() (licenses string) {
	licenses = ""	
	for _,l:= range set {
		licenses += fmt.Sprintf("%s\n", l.FullName)
	}
	return licenses
}

func (licenses LicenseSet) Get(fullname string) (lic *License, err error){
	for i := range licenses {
		if licenses[i].FullName == fullname {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license %s", fullname))
}
func (licenses LicenseSet) GetAlias(alias string) (lic *License, err error){
	for i := range licenses {
		if licenses[i].Alias == alias {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license's alias %s", alias))
}

type ProjectID struct {
	name    string // any valid package name
	version Version
}

type Project struct {
	workingDir   string      // transient workding directory aboslute path
	name         string      // package name
	dependencies []ProjectID // contains the current project's dependencies
	license      License      // one of the predefined licenses
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
	//	fmt.Printf("writing package to %s\n", dst)
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
func (p *Project) License() License {
	return p.license
}

func (p *Package) Name() string {
	return p.self.name
}
func (p *Package) Version() Version {
	return p.version
}

func (p *ProjectID) Name() string {
	return p.name
}
func (p *ProjectID) Version() Version {
	return p.version
}

func (p *Project) SetWorkingDir(pwd string) {
	p.workingDir = pwd
}
func (p *Project) SetName(name string) {
	p.name = name
}

func (p *Project) SetLicense(license License) {
	if _,err := Licenses.Get(license.FullName) ; err != nil {
		panic(err)
	}
	p.license = license
}


func (p *Project) Dependencies() []ProjectID {
	return p.dependencies[:]
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
	err = EncodeProject(f, p)
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
	if err != nil {
		return
	}

	defer gz.Close()
	tr := tar.NewReader(gz)
	dst := p.self.workingDir
	//fmt.Printf("unpacking to %s\n", dst)
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// make the target file
		ndst := filepath.Join(dst, hdr.Name)
		os.MkdirAll(path.Dir(ndst), os.ModeDir|os.ModePerm) // mkdir -p
		//fmt.Printf("%s\n", ndst)
		df, err := os.Create(ndst)
		if err != nil {
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

type ProjectFile struct { // TODO append a version number to make it possible to handle "format upgrade"
	Name         string
	Dependencies []ProjectIDFile
	License      string // one of the value in the restricted list
}
type PackageFile struct {
	Self      ProjectFile
	Version   string
	Timestamp time.Time
}

type ProjectIDFile struct {
	Name, Version string
}

type LocalRepositoryFile struct {
	Remotes []RemoteFile
}

type RemoteFile struct {
	Name, Url string
}

func DecodePackage(r io.Reader, p *Package) (err error) {
	pf := new(PackageFile)
	err = json.NewDecoder(r).Decode(pf)
	if err != nil {
		return
	}
	DecodePackageFile(*pf, p)
	return
}
func EncodePackage(w io.Writer, p Package) (err error) {
	pf := EncodePackageFile(p)
	err = json.NewEncoder(w).Encode(pf)
	return

}

func EncodeProject(w io.Writer, p Project) (err error) {
	pf := EncodeProjectFile(p)
	err = json.NewEncoder(w).Encode(pf)
	return

}

func DecodeProject(r io.Reader, p *Project) (err error) {
	pf := new(ProjectFile)
	err = json.NewDecoder(r).Decode(pf)
	if err != nil {
		return
	}
	err = DecodeProjectFile(*pf, p)
	return
}

func EncodeLocalRepository(w io.Writer, p LocalRepository) (err error) {
	pf := EncodeLocalRepositoryFile(p)
	err = json.NewEncoder(w).Encode(pf)
	return

}

func DecodeLocalRepository(r io.Reader, p *LocalRepository) (err error) {
	pf := new(LocalRepositoryFile)
	err = json.NewDecoder(r).Decode(pf)
	if err != nil {
		return
	}
	DecodeLocalRepositoryFile(*pf, p)
	return
}

func DecodeRemoteFile(rf RemoteFile) RemoteRepository {
	u, err := url.Parse(rf.Url)
	if err != nil {
		panic(err)
	}
	return NewRemoteRepository(rf.Name, *u)
}

func EncodeRemoteFile(r RemoteRepository) *RemoteFile {
	u := r.Path()
	return &RemoteFile{
		Name: r.Name(),
		Url:  u.String(),
	}
}

func DecodeLocalRepositoryFile(pf LocalRepositoryFile, p *LocalRepository) {
	for _, d := range pf.Remotes {
		p.RemoteAdd(DecodeRemoteFile(d))
	}
}
func DecodeProjectFile(pf ProjectFile, p *Project) (err error){
	p.name = pf.Name
	for _, d := range pf.Dependencies {
		v, _ := ParseVersion(d.Version)
		p.AppendDependency(NewProjectID(d.Name, v))
	}
	if l,e:= Licenses.Get(pf.License); e!= nil {
		err = errors.New(fmt.Sprintf(`Illegal license: "%s" was expecting one of: %s`, pf.License, Licenses) )
	} else {
		p.license = *l
	}
	return
}

func DecodePackageFile(pf PackageFile, p *Package) {
	prj := new(Project)
	DecodeProjectFile(pf.Self, prj)
	p.self = *prj
	p.timestamp = pf.Timestamp
	v, _ := ParseVersion(pf.Version)
	p.version = v
}

func EncodePackageFile(p Package) *PackageFile {
	return &PackageFile{
		Self:      *EncodeProjectFile(p.self),
		Timestamp: p.timestamp,
		Version:   p.version.String(),
	}
}

func EncodeProjectIDFile(p ProjectID) *ProjectIDFile {
	return &ProjectIDFile{
		Name:    p.name,
		Version: p.version.String(),
	}
}
func EncodeProjectFile(p Project) *ProjectFile {
	dep := make([]ProjectIDFile, 0, len(p.dependencies))
	for _, d := range p.dependencies {
		dep = append(dep, *EncodeProjectIDFile(d))
	}

	return &ProjectFile{
		Name:         p.name,
		Dependencies: dep,
		License:      p.license.FullName,
	}
}
func EncodeLocalRepositoryFile(p LocalRepository) *LocalRepositoryFile {
	dep := make([]RemoteFile, 0, len(p.remotes))
	for _, d := range p.remotes {
		dep = append(dep, *EncodeRemoteFile(d))
	}

	return &LocalRepositoryFile{
		Remotes: dep,
	}
}
