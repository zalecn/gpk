package gor

import (
	"fmt"
	"time"
	"path/filepath"
	"strings"
	"strconv"
)

// versions are NOT linear ones. there are three kinds of versions:
//	semi linear ones (1, 2, 3 (should I accept sub numbers ? yes ) (there is a total order, and always the possibility to insert a new version after another one
// 	snapshot ones (i.e linear in time)
// branch ones  non linear history version (there is no way to tell if one is after the other, and that doesn't make sense anyway)

type Version struct {
	Root      string
	Parts     [4]uint8 // successive version derived from the root, the default root is ""
	Timestamp time.Time
	// should I allow classifiers ? no way
}

// TODO provide ordering functions for versions

func NewVersion() *Version {
	return &Version{
		Root:      "",
		Parts:     [4]uint8{0, 0, 0, 0},
		Timestamp: time.Now(),
	}
}


func ParseVersion(version string) Version {
	v := NewVersion()
	parts := strings.Split(version, "-")
	v.Root = parts[0]
	if len(parts) == 2 {
		sparts := strings.Split( parts[1] , "." )
		v.Parts[0]= atoi(sparts[0])  
		v.Parts[1]= atoi(sparts[1])  
		v.Parts[2]= atoi(sparts[2])  
		v.Parts[3]= atoi(sparts[3])  
	}
	return *v
}

func atoi(s string) uint8{
	i,_ := strconv.ParseUint(s, 10, 8)
	return uint8(i)
}


func (v *Version) String() string {
	return fmt.Sprintf("%s-%d.%d.%d.%d", v.Root, v.Parts[0], v.Parts[1], v.Parts[2], v.Parts[3])
}

func (v *Version) Path() string {
	return filepath.Join(v.Root, fmt.Sprintf("%d.%d.%d.%d", v.Parts[0], v.Parts[1], v.Parts[2], v.Parts[3] ) )
}
