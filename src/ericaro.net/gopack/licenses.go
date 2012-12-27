package gopack

import (
	"errors"
	"fmt"
)

type LicenseSet []License
type License struct {
	FullName, Alias string
}

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

func (set LicenseSet) String() (licenses string) {
	licenses = ""
	for _, l := range set {
		licenses += fmt.Sprintf("%s\n", l.FullName)
	}
	return licenses
}

func (licenses LicenseSet) Get(fullname string) (lic *License, err error) {
	for i := range licenses {
		if licenses[i].FullName == fullname {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license %s", fullname))
}
func (licenses LicenseSet) GetAlias(alias string) (lic *License, err error) {
	for i := range licenses {
		if licenses[i].Alias == alias {
			lic = &licenses[i]
			return
		}
	}
	return nil, errors.New(fmt.Sprintf("Unknown or unsupported license's alias %s", alias))
}
