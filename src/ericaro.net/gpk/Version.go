package gpk

import (
	"fmt"
)

//var (
//PreReleaseCharRegexp string = `[0-9A-Za-z-](\.[0-9A-Za-z-])*`
//)

//Version is a struct that hold all [http://semver.org/ semantic version] components.
type Version struct {
	major, minor, patch uint32
	preRelease, build   string
}

// TODO add tests and methods to this struct
// TODO add persistence format control ( string back and forth is a good objective)

func (v Version) String() (version string) {
	version = fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)

	if len(v.preRelease) != 0 {
		version += fmt.Sprintf("-%s", v.preRelease)
	}
	if len(v.build) != 0 {
		version += fmt.Sprintf("+%s", v.build)
	}
	return
}
