package entities

import ()

//a Package is a pure, in-memory representation of a Package
type Package struct {
	Group, Artifact, Root      string
	Major, Minor, Micro, Build int16
	ContentBlob                string // tar.gz of the content
}
