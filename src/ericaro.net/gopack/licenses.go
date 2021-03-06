package gopack

import (
	"errors"
	"fmt"
)

//License is a simple fullname alias collection of licenses
type License struct {
	FullName, Alias string
}

//IsValid return true if the license belongs to one of the preset license
func (l License) IsValid() bool {
	r, err := Licenses.Get(l.FullName)
	return err == nil && r.Alias == l.Alias
}
//IsOSS return true if the license belongs to one of the preset OSS license
func (l License) IsOSS() bool {
	return l.IsValid( ) && l.Alias != "OCS"
}

//LicenseSet is a slice of license (used to create a singleton of licenses)
type LicenseSet []License

//The singleton of Valid licenses
var (
	Licenses LicenseSet = ([]License{
		License{"Apache License 2.0", "ASF"},
		License{"Eclipse Public License 1.0", "EPL"},
		License{"GNU GPL v2", "GPL2"},
		License{"GNU GPL v3", "GPL3"},
		License{"GNU Lesser GPL", "LGPL"},
		License{"MIT License", "MIT"},
		License{"Mozilla Public License 1.1", "MPL"},
		License{"New BSD License", "BSD"},
		License{"Other Open Source", "OOS"},
		License{"Other Closed Source", "OCS"},
	})
)
//String return a pretty print version of the license list
func (set LicenseSet) String() (licenses string) {
	licenses = ""
	for _, l := range set {
		licenses += fmt.Sprintf("%s\n", l.FullName)
	}
	return licenses
}

//Get return a license by its fullname
func (licenses LicenseSet) Get(fullname string) (lic *License, err error) {
	for i := range licenses {
		if licenses[i].FullName == fullname {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license %s", fullname))
}

//GetAlias return a license by its alias
func (licenses LicenseSet) GetAlias(alias string) (lic *License, err error) {
	for i := range licenses {
		if licenses[i].Alias == alias {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license's alias %s", alias))
}
