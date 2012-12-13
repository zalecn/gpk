package got

import (
	"fmt"
	"strings"
	"path/filepath"
)

type VersionReference struct {
	Root  string
	Parts string // successive version derived from the root, the default root is ""
}

func (v VersionReference) String() string {
	return fmt.Sprintf("%s-%s", v.Root, v.Parts)
}

func ParseVersionReference(version string) VersionReference {
	v := VersionReference{}
	parts := strings.Split(version, "-")
	v.Root = parts[0]
	if len(parts) == 2 {
		v.Parts = parts[1]
	}
	return v
}

func (vref *VersionReference) Version() (v *Version) {
	ve:= ParseVersion(vref.String())
	return &ve
}

func (v VersionReference) Path() string {
	return filepath.Join(v.Root, v.Parts)
}
