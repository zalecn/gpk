//package semver contains an as close as possible implementation of http://semver.org
//
//Semver norm does not require digits to be used (they require it for  "normal" versions )
// so we use this hole to define "snapshot" version:
// version 0.0.0 are considered snapshots. the digits can be skipped, so does the prelease dash ("-")
// this way "master" is a suitable semver, that fully qualifies to 0.0.0-master
package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// regexp elements for the syntax
const (
	digits = `(\d+)?(?:\.(\d+)(?:\.(\d+))?)?`
	sub    = `[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*`
	all    = `%s(?:\-?(%s))?(?:\+(%s))?`
)

//regexp for a semver
var (
	SemVer, _ = regexp.Compile(fmt.Sprintf(all, digits, sub, sub))
)

//Version is a struct that hold all [http://semver.org/ semantic version] components.
type Version struct {
	major, minor, patch uint32
	pre, build          string // for now pre and build part are not parsed (splitted into . ) 
}

//NewVersion creates a new standard semver
func NewVersion(major, minor, patch uint32, pre, build string) *Version {
	return &Version{
		major:major,
		minor: minor,
		patch: patch,
		pre: pre,
		build: build,
	}
}

//String pretty prints the version.

func (v Version) String() (version string) {
	version = fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.major == 0 && v.minor == 0 && v.patch == 0 {
		version = fmt.Sprintf("%s", v.pre)
	} else {
		if len(v.pre) != 0 {
			version += fmt.Sprintf("-%s", v.pre)
		}
	}
	if len(v.build) != 0 {
		version += fmt.Sprintf("+%s", v.build)
	}
	return
}
//Digits return the three digits of the semver
func (v Version) Digits() (major, minor, patch uint32) {
	return v.major, v.minor, v.patch
}
//PreRelease returns the Pre Release part unparsed
func (v Version) PreRelease() string {
	return v.pre
}
//Build return the Build part unparsed
func (v Version) Build() string {
	return v.build
}
//IsSnapshot return true if the three digits are equals to 0
func (v Version) IsSnapshot() bool {
	return v.major == 0 && v.minor == 0 && v.patch == 0
}

func atoi(s string) uint32 {
	i, _ := strconv.ParseUint(s, 10, 8)
	return uint32(i)
}

//ParseVersion reads the version. It supports ommiting the digits (that are considered 0.0.0) and the leading "-"
func ParseVersion(v string) (version Version, err error) {
	v = strings.Trim(v, " {}[]\"'")
	parts := SemVer.FindStringSubmatch(v)
	//fmt.Printf("%24s -> %v\n",v, parts[1:])
	version = Version{
		major: atoi(parts[1]),
		minor: atoi(parts[2]),
		patch: atoi(parts[3]),
		pre:   parts[4],
		build: parts[6],
	}
	return
}
