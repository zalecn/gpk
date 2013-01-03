package protocol

import (
	"encoding/json"
	"ericaro.net/gopack/semver" // todo move this version to another package (standalone semantic version package
	"net/url"
	"path/filepath"
	"time"
)

//PID represent a Project ID through the internet. Can be either passed as parameter to a query, or returned as a list in a search result
// contains references to package name, version, the package timestamp (used in update procotol), and a Token
type PID struct {
	Name      string
	Version   semver.Version
	Timestamp *time.Time
	Token     *Token // is optional
}

//Path computes the relative path to the expected package (usually <name> / <version> )
func (p PID) Path() string {
	return filepath.Join(p.Name, p.Version.String())
}

//InParameter encode this PID in an url Values
func (pid *PID) InParameter(v *url.Values) {
	// prepare central server query args
	v.Set("n", pid.Name)
	v.Set("v", pid.Version.String())
	if pid.Timestamp != nil {
		v.Set("t", pid.Timestamp.Format(time.ANSIC))
	}
	if pid.Token != nil {
		v.Set("k", pid.Token.FormatURL())
	}
}

//FromParameter decode a pid from an url.Values
func FromParameter(v *url.Values) (pid *PID, err error) {
	pid = &PID{}
	pid.Name = v.Get("n") // todo validate the syntax
	pid.Version, err = semver.ParseVersion(v.Get("v"))
	if err != nil {
		return // this is not an optional parameter
	}

	t, err := time.Parse(time.ANSIC, v.Get("t"))
	k, err := ParseURLToken(v.Get("k"))
	
	pid.Timestamp = &t
	pid.Token = k
	return
}

//UnmarshalJSON part on the json protocol to make PID json-marshallable
func (pid *PID) UnmarshalJSON(data []byte) (err error) {
	type Pidfile struct {
		Name      string
		Version   string
		Timestamp string
	}
	var pf Pidfile
	json.Unmarshal(data, &pf)
	pid.Name = pf.Name
	pid.Version, err = semver.ParseVersion(pf.Version)
	if err != nil {
		return
	}
	t, err := time.Parse(time.ANSIC, pf.Timestamp)
	if err != nil {
		err = nil
	} else {
		pid.Timestamp = &t
	}
	return
}

//MarshalJSON part on the json protocol to make PID json-marshallable
func (pid *PID) MarshalJSON() ([]byte, error) {
	type Pidfile struct {
		Name      string
		Version   string
		Timestamp string
	}
	pf := Pidfile{
		Name:    pid.Name,
		Version: pid.Version.String(),
	}
	if pid.Timestamp != nil {
		pf.Timestamp = pid.Timestamp.Format(time.ANSIC)
	}
	return json.Marshal(pf)
}
