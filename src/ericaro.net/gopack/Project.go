package gopack

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"fmt"
)

type Project struct {
	workingDir   string      // transient workding directory aboslute path
	name         string      // package name
	dependencies []ProjectID // contains the current project's dependencies
	license      License     // one of the predefined licenses
	// TO be added build time , and test dependencies
}

func ReadProject() (p *Project, err error) {
	p = &Project{}
	err = JsonReadFile(GpkFile, p)
	p.workingDir, _ = filepath.Abs(filepath.Dir(GpkFile))
	return
}

//Write project  info to where it belongs (project holds working dir info)
func (p Project) Write() (err error) {
	dst := filepath.Join(p.workingDir, GpkFile)
	err = JsonWriteFile(dst, p)
	return
}

func (p *Project) WorkingDir() string {
	return p.workingDir
}

func (p *Project) Name() string {
	return p.name
}
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
	
	if l,e:= Licenses.Get(pf.License); e!= nil {
		err = errors.New(fmt.Sprintf(`Illegal license: "%s" was expecting one of: %s`, pf.License, Licenses) )
	} else {
		p.license = *l
	}
	return
}

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
