//Package gopack contains the main objects for Gopack that is a software dependency management tool for Go.
package gopack

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	. "ericaro.net/gopack/semver"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

const (
	GpkFile = ".gpk"
)

//Package represent a Go project packaged.
// it includes a reference to a Project content , a version and the timestamp when it was created (for traceability and snapshot management)
type Package struct {
	self      Project
	version   Version
	timestamp time.Time

	// more to come, like sha1,signature, snapshot/release
	// add also go1 , i.e the target go runtime.

}

//ReadPackageFile local info from the specified gopackage file
func ReadPackageFile(gpkPath string) (p *Package, err error) {
	p = &Package{}
	err = JsonReadFile(gpkPath, p)
	p.self.workingDir, _ = filepath.Abs(filepath.Dir(gpkPath))
	return
}

//Timestamp return the timestamp, i.e the date the package was created
func (p *Package) Timestamp() time.Time {
	return p.timestamp
}

//Write down package info into this package InstallDir
func (p *Package) Write() (err error) {
	dst := filepath.Join(p.self.workingDir, GpkFile)
	err = JsonWriteFile(dst, p)
	return
}

//InstallDir is the place where the package is installed
func (p *Package) InstallDir() string {
	return p.self.workingDir
}

//Name The package name. As in referenced in dependency management, and in go sources.
func (p *Package) Name() string {
	return p.self.name
}

//License the license driving this package.
func (p *Package) License() License {
	return p.self.License()
}

//Dependencies return the list of dependencies declared in this package's project
func (p *Package) Dependencies() []ProjectID {
	return p.self.Dependencies()
}

//Version this package semantic version
func (p *Package) Version() Version {
	return p.version
}

//Path converts this package into the path it should have in the repository layout
func (p *Package) Path() string {
	return filepath.Join(p.self.name, p.version.String())
}

//ID computes the ProjectID of this package, the way it should be referenced to.
func (p *Package) ID() ProjectID {
	return ProjectID{
		name:    p.self.name,
		version: p.version,
	}
}

//ReadPackageInPackage reads the .gpk file within the tar in memory. It does not set the Root
func ReadPackageInPackage(in io.Reader) (p *Package, err error) {
	//fmt.Printf("Parsing in memory package\n")
	gz, err := gzip.NewReader(in)
	if err != nil {
		return
	}
	tr := tar.NewReader(gz)
	defer gz.Close()

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				err = errors.New(fmt.Sprintf("Invalid package format, %v is missing", GpkFile))
			}
			break
		}
		//fmt.Printf("hdr %v\n", hdr )
		if hdr.Name == GpkFile {
			p = &Package{}
			err = json.NewDecoder(tr).Decode(p)
			return p, err
		}
	}
	return
}

//Unpack a Package in a reader (a tar.gzed stream) into its InstallDir directory
func (p *Package) Unpack(in io.Reader) (err error) {
	return Unpack(p.self.workingDir, in)

}

//Pack copy the current Package into a Writer. It with write it down in tar.gzed format
func (p *Package) Pack(w io.Writer) (err error) {
	return p.packType(PACK_SRC, w)
}
//Pack copy the current Package exec into a Writer. It with write it down in tar.gzed format
func (p *Package) PackExecutables(w io.Writer) (err error) {
	return p.packType(PACK_EXEC, w)
}
const (
	PACK_SRC  = iota
	PACK_EXEC = iota

//PACK_PKG = iota

)

func (p *Package) packType(typ int, w io.Writer) (err error) {
	gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return
	}

	tw := tar.NewWriter(gz)
	defer gz.Close()
	defer tw.Close()

	//prepare recursive handlers
	dirHandler := func(ldst, lsrc string) (err error) {
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		err = TarFile(ldst, lsrc, tw)
		return
	}
	// same remark as the "install" function
	switch typ {

	case PACK_SRC:
		p.self.ScanProjectSrc("", dirHandler, fileHandler)
	case PACK_EXEC:
		p.self.ScanBinPlatforms("", fileHandler)
	}
	//walkDir("src", filepath.Join(p.self.workingDir, "src"), dirHandler, fileHandler)
	// copy the package .gpk
	TarFile(filepath.Join("", GpkFile), filepath.Join(p.self.workingDir, GpkFile), tw)
	return
}

//UnmarshalJSON part of the json protocol
func (p *Package) UnmarshalJSON(data []byte) (err error) {
	type PackageFile struct {
		Self      *Project
		Version   string
		Timestamp time.Time
	}
	var pf PackageFile
	json.Unmarshal(data, &pf)

	p.self = *pf.Self
	p.timestamp = pf.Timestamp
	v, _ := ParseVersion(pf.Version)
	p.version = v
	return
}

//MarshalJSON part of the json protocol
func (p *Package) MarshalJSON() ([]byte, error) {
	type PackageFile struct {
		Self      *Project
		Version   string
		Timestamp time.Time
	}
	pf := PackageFile{
		Self:      &p.self,
		Timestamp: p.timestamp,
		Version:   p.version.String(),
	}
	return json.Marshal(pf)
}
