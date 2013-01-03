package gopack

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"os"
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

//ReadProject read project from the current dir
func ReadProject() (p *Project, err error) {
	p = &Project{}
	path, err := os.Getwd()
	if err != nil{
	return
	}
	for ; !FileExists(filepath.Join(path, GpkFile)) && path != "/"; path =  filepath.Dir(path){}
	if path != "/" { 
		err = JsonReadFile(filepath.Join(path, GpkFile), p)
	}
	p.workingDir, _ = filepath.Abs(filepath.Dir(GpkFile)) // TODO try a up dir lookup
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
			ref = &r
			is = append(is, i)
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
	// last bit of slice
	p.dependencies = dep
	return
}

//UnmarshalJSON part of the json protocol
func (p *Project) UnmarshalJSON(data []byte) (err error) {
	type ProjectFile struct { // TODO append a version number to make it possible to handle "format upgrade"
		Name         string
		Dependencies []ProjectID
		License      string // one of the value in the restricted list
	}
	var pf ProjectFile
	json.Unmarshal(data, &pf)

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
		Name         string
		Dependencies []ProjectID
		License      string // one of the value in the restricted list
	}
	pf := ProjectFile{
		Name:         p.name,
		Dependencies: p.dependencies,
		License:      p.license.FullName,
	}
	return json.Marshal(pf)
}
