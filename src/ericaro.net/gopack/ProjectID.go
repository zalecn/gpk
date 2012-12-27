package gopack

import (
	. "ericaro.net/gopack/semver"
	"encoding/json"
	"path/filepath"
	"fmt"
)

type ProjectID struct {
	name    string // any valid package name
	version Version
}

func NewProjectID(name string, version Version) ProjectID {
	return ProjectID{name: name, version: version}
}


func (p *ProjectID) Name() string {
	return p.name
}
func (p *ProjectID) Version() Version {
	return p.version
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
