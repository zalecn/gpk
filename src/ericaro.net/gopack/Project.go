package gopack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	//"strings"
)

const (
	GpkFileVersion = "1.0.0"
)

//Project is a Go Project, plus some metadata:
// a workingDir that must be layouted as a standard go project (a src, bin, pkg directory)
// a unique name
// a list of dependency references (name, version)
// and a license for the source code. This is required because we cannot move around licenses if we aren't allowed to.
type Project struct {
	workingDir   string      // transient workding directory aboslute path
	name         string      // package name
	dependencies []ProjectID // contains the current project's dependencies
	license      License     // one of the predefined licenses
	// TO be added build time , and test dependencies
}

//ReadProject read project from the current dir, or parent's one (recursively)
func ReadProject() (p *Project, err error) {
	p = &Project{}
	path, err := os.Getwd()
	if err != nil {
		return
	}
	for ; !FileExists(filepath.Join(path, GpkFile)) && path != "/"; path = filepath.Dir(path) {
	}
	if path != "/" {

		gpk := filepath.Join(path, GpkFile)
		err = JsonReadFile(gpk, p)
		p.workingDir, _ = filepath.Abs(filepath.Dir(gpk))
	} else { // return a new empty prj set using the current dir
		p.workingDir, _ = filepath.Abs(filepath.Dir(GpkFile))
		err = errors.New("no .gpk file found in the current working directory hierarchy.")
	}
	return
}

//Write down the project into the WorkingDir
func (p *Project) Write() (err error) {
	dst := filepath.Join(p.workingDir, GpkFile)
	err = JsonWriteFile(dst, p)
	return
}

//WorkingDir the directory containing the project
func (p *Project) WorkingDir() string {
	return p.workingDir
}

//Name the project unique name, it must be the package name.
func (p *Project) Name() string {
	return p.name
}

//License the project license for source code
func (p *Project) License() License {
	return p.license
}

func (p *Project) SetWorkingDir(pwd string) {
	p.workingDir = pwd
}
func (p *Project) SetName(name string) {
	p.name = name
}

func (p *Project) SetLicense(license License) {
	if _, err := Licenses.Get(license.FullName); err != nil {
		panic(err)
	}
	p.license = license
}

//Dependencies is a slice of ProjectID used by this project. Caveat this is not the whole dependency tree, just the root dependencies
func (p *Project) Dependencies() []ProjectID {
	return p.dependencies[:]
}

//AppendDependency append some root dependencies
func (p *Project) AppendDependency(ref ProjectID) (rem *ProjectID) {
	//BUG: check that the new dependencies does not "replace" existing one (there shall be only one dependency per package name
	rem = p.RemoveDependency(ref.Name()) // first remove it
	p.dependencies = append(p.dependencies, ref)
	return
}

//RemoveDependency removes the dependency by name, and return the removed reference
func (p *Project) RemoveDependency(name string) (ref *ProjectID) {
	src := p.dependencies
	// first compute the dependencies to be removed (yes accidentally there might be more than one
	is := make([]int, 0, len(src))
	for i, r := range src {
		if r.Name() == name {
			is = append(is, i)
			ref = NewProjectID(name, r.Version())
		}
	}
	length := len(is)
	if length == 0 { // nothing to do
		return nil
	}
	// now apply the removal, unfortunately, I don't how to make it easier

	// I create a new slice of project id
	dep := make([]ProjectID, 0, len(src)-length)
	// and copy all but the removed
	if is[0] > 0 {
		dep = append(dep, src[0:is[0]]...)
	}
	for j := 0; j < length-1; j++ {
		s, e := is[j]+1, is[j+1]
		dep = append(dep, src[s:e]...)
	}
	// oops forget to finish the stuf

	s, e := is[length-1]+1, len(src)
	if s < len(src) { // the last removed, is not the last in the src , copy the trailing stuff
		dep = append(dep, src[s:e]...)
	}

	// last bit of slice
	p.dependencies = dep
	return
}

//ScanProjectSrc recursively walk into src directory  and fire callbacks to dirHandler, and fileHandler
// dirHandler is called on every directory dst beeing a join between the path passed at first, and the relative path to the current directory, and src and absolute one. 
// fileHandler is called on every source file. For now source files are just .go files
// dst is just a path that is used as root for the dst path in the handler.
// for instance, if you scan a dir in you prj src/foo/bar and you initially passed a path "toto" then, handlers will be called with toto/src/foo/bar 
func (p *Project) ScanProjectSrc(dst string, dirHandler, srcHandler func(dst, src string) error) error {
	src := filepath.Join(p.WorkingDir(), "src")
	dst = filepath.Join(dst, "src")
	return scanProjectSrc(dst, src, dirHandler, srcHandler)
}

// recursive impl of eponym function
func scanProjectSrc(dst, src string, dirHandler, srcHandler func(dst, src string) error) error {
	if dirHandler != nil {
		dirHandler(dst, src)
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	subdir, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range subdir {
		switch {

		case fi.IsDir():
			ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
			err = scanProjectSrc(ndst, nsrc, dirHandler, srcHandler)
			if err != nil {
				return err
			}
		default: //case strings.HasSuffix(fi.Name(), ".go"):
			ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
			err := srcHandler(ndst, nsrc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//ScanBinPlatforms scans the bin directory for binaries organized by platform (and yes it put the current platform in the right place)
func (p *Project) ScanBinPlatforms(dst string, srcHandler func(dst, src string) error) error {

	src := filepath.Join(p.WorkingDir(), "bin")
	dst = filepath.Join(dst, "bin")
	// scan bin as the current platform

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	files, err := file.Readdir(-1)
	if err != nil {
		return err
	}
	// first scan the current bin as if it was in a bin/{platform} one
	platform := runtime.GOOS + "_" + runtime.GOARCH
	localdst := filepath.Join(dst, platform)
	err = scanBinPlatform(localdst, src, files, srcHandler) // scan real files and add them in the bin/{current_platform}
	if err != nil {
		return err
	}

	// then rescan only subdirs of bin

	for _, fi := range files {
		if fi.IsDir() {
			// this is a platform actually
			ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
			nfile, err := os.Open(nsrc)
			nfiles, err := nfile.Readdir(-1)
			if err != nil {
				return err
			}
			scanBinPlatform(ndst, nsrc, nfiles, srcHandler)
		}
		if err != nil {
			return err
		}
	}
	return nil

}

// recursive impl of eponym function
func scanBinPlatform(dst, src string, files []os.FileInfo, srcHandler func(dst, src string) error) error {
	for _, fi := range files {
		if !fi.IsDir() {
			ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
			err := srcHandler(ndst, nsrc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//scanPkg recursive scans the pkg directory for .a organized by platform
func scanPkg(dst, src string, srcHandler func(dst, src string) error) error {

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	files, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range files {
		ndst, nsrc := filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name())
		if fi.IsDir() {
			err = scanPkg(ndst, nsrc, srcHandler)
			if err != nil {
				return err
			}
		} else {
			err := srcHandler(ndst, nsrc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil // no error so far
}

//ScanPkg scans the pkg directory for .a organized by platform
// convenient method to package the pkg in a tar.gz and send them to a server
// but we need a way to store the pkg per golang version as they might not be compatible.
// nevertheless golang is not ready for that (there is a bug that make the distribution of binaries only impossible).
// this is a back code for later.
func (p *Project) ScanPkg(dst string, srcHandler func(dst, src string) error) error {

	src := filepath.Join(p.WorkingDir(), "pkg")
	dst = filepath.Join(dst, "pkg")
	return scanPkg(dst, src, srcHandler) // scan real files and add them in the bin/{current_platform}

}

//Packages returns the list of packages present in the src dir, and eventually containing tests
func (p *Project) Packages() []string {
	// srcDir is where to scan for sources imports
	srcDir := filepath.Join(p.WorkingDir(), "src")
	dirs, _ := ScanPackages(srcDir)
	packages := make([]string, len(dirs))
	for i, d := range dirs {
		p, err := filepath.Rel(srcDir, d)
		packages[i] = p
		if err != nil {
			panic(err) //really unexpected
		}
	}
	return packages
}

//UnmarshalJSON part of the json protocol
func (p *Project) UnmarshalJSON(data []byte) (err error) {
	type ProjectFile struct { // TODO append a version number to make it possible to handle "format upgrade"
		FormatVersion string
		Name          string
		Dependencies  []ProjectID
		License       string // one of the value in the restricted list
	}
	var pf ProjectFile
	json.Unmarshal(data, &pf)
	if pf.FormatVersion != GpkFileVersion {
		log.Printf("Warning: Unknown format version \"%s\"", pf.FormatVersion)
	}

	p.name = pf.Name
	p.dependencies = pf.Dependencies

	if l, e := Licenses.Get(pf.License); e != nil {
		err = errors.New(fmt.Sprintf(`Illegal license: "%s" was expecting one of: %s`, pf.License, Licenses))
	} else {
		p.license = *l
	}
	return
}

//MarshalJSON part of the json protocol
func (p *Project) MarshalJSON() ([]byte, error) {
	type ProjectFile struct { // TODO append a version number to make it possible to handle "format upgrade"
		FormatVersion string
		Name          string
		Dependencies  []ProjectID
		License       string // one of the value in the restricted list
	}
	pf := ProjectFile{
		FormatVersion: GpkFileVersion,
		Name:          p.name,
		Dependencies:  p.dependencies,
		License:       p.license.FullName,
	}
	return json.Marshal(pf)
}
