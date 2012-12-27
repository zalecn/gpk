package protocol

import (
	"encoding/json"
	"ericaro.net/gopack/semver" // todo move this version to another package (standalone semantic version package
	"net/url"
	"time"
	"path/filepath"
)

//PID represent a Project ID through the internet. Can be either passed as parameter to a query, or returned as a list in a search result
type PID struct {
	Name      string
	Version   semver.Version
	Timestamp *time.Time
	Token     *Token // is optional
}


func (p PID) Path() string {
	return filepath.Join(p.Name, p.Version.String())
}

func (pid *PID) InParameter(v *url.Values) {
	// prepare central server query args
	v.Set("n", pid.Name)
	v.Set("v", pid.Version.String())
	if pid.Timestamp != nil {
		v.Set("t", pid.Timestamp.Format(time.ANSIC))
	}
	if pid.Token != nil {
		v.Set("k", pid.Token.Format())
	}
}

//FromParameter fills a PID object from the url.Values, reverse of InParameter Function
func FromParameter(v *url.Values) (pid *PID, err error) {
	pid = &PID{}
	pid.Name = v.Get("n") // todo validate the syntax
	pid.Version, err = semver.ParseVersion(v.Get("v"))
	if err != nil {
		return // this is not an optional parameter
	}

	t, err := time.Parse(time.ANSIC, v.Get("t"))
	k, err := DecodeString(v.Get("k"))

	pid.Timestamp = &t
	pid.Token = k
	return
}

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
		return
	}
	pid.Timestamp = &t
	return
}

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

