package gopack

import (
	"encoding/json"
	. "ericaro.net/gopack/semver"
	"fmt"
	"path/filepath"
)

//ProjectID is a simple symbolic reference to a Package, made of a (name, version)
type ProjectID struct {
	name    string // any valid package name
	version Version
}

//NewProjectID creates a new ProjectID
func NewProjectID(name string, version Version) *ProjectID {
	return &ProjectID{name: name, version: version}
}

//Name the name of the package this ProjectID references
func (p *ProjectID) Name() string {
	return p.name
}

//Version the version of the package this ProjectID references
func (p *ProjectID) Version() Version {
	return p.version
}

//Path converts this project reference into the path it should have in the repository layout
func (d ProjectID) Path() string {
	return filepath.Join(d.name, d.version.String())
}

//Equals compare two PprojectID equals if they are supposed to reference the same Package
func (this ProjectID) Equals(that ProjectID) bool {
	return this.id() == that.id()
}

func (d ProjectID) id() string {
	return fmt.Sprintf("%s %s", d.name, d.version.String())
}

//String returns a simple " " separated representation of the pid pair
func (d ProjectID) String() string {
	return d.id()
}

//UnmarshalJSON part of the json protocol
func (pid *ProjectID) UnmarshalJSON(data []byte) (err error) {
	type ProjectIDFile struct {
		Name, Version string
	}
	var pf ProjectIDFile
	json.Unmarshal(data, &pf)
	pid.name = pf.Name
	pid.version, err = ParseVersion(pf.Version)
	return
}

//MarshalJSON part of the json protocol
func (p *ProjectID) MarshalJSON() ([]byte, error) {
	type ProjectIDFile struct {
		Name, Version string
	}
	pf := ProjectIDFile{
		Name:    p.name,
		Version: p.version.String(),
	}
	return json.Marshal(pf)
}
