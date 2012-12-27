package gopack

import (
	. "ericaro.net/gopack/semver"
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"path/filepath"
	"time"
	"errors"
	"fmt"
)

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

func (p *Package) Timestamp() time.Time {
	return p.timestamp
}

//Write package  info to where it belongs (package holds working dir info)
func (p Package) Write() (err error) {
	dst := filepath.Join(p.self.workingDir, GpkFile)
	err = JsonWriteFile(dst, p)
	return 
}


func (p *Package) InstallDir() string {
	return p.self.workingDir
}

func (p *Package) Name() string {
	return p.self.name
}
func (p *Package) Version() Version {
	return p.version
}

//Path converts this project reference into the path it should have in the repository layout
func (p *Package) Path() string {
	return filepath.Join(p.self.name, p.version.String())
}

func (p *Package) ID() ProjectID {
	return ProjectID{
		name:    p.self.name,
		version: p.version,
	}
}

//ReadProjectTar reads the .gopackage file within the tar in memory. It does not set the Root
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

//Untar reads the .gopackage file within the tar in memory. It does not set the Root
func (p *Package) Unpack(in io.Reader) (err error) {
	return  Unpack(p.self.workingDir, in)
	
}

//PackageProject into a tar writer
func (p *Package) Pack(w io.Writer) (err error) {
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
	walkDir("src", filepath.Join(p.self.workingDir, "src"), dirHandler, fileHandler)
	// copy the package .gpk
	TarFile(filepath.Join("", GpkFile), filepath.Join(p.self.workingDir, GpkFile), tw)
	// or rewrite it (and edit it on the fly ?
	//	buf := new(bytes.Buffer)
	//	json.NewEncoder(buf).Encode(p)
	//	TarBuff(filepath.Join("/", GpkFile), buf, tw)

	return
}

func (p *Package) UnmarshalJSON(data []byte) (err error) {
	type PackageFile struct {
		Self      Project
		Version   string
		Timestamp time.Time
	}
	var pf PackageFile
	json.Unmarshal(data, &pf)

	p.self = pf.Self
	p.timestamp = pf.Timestamp
	v, _ := ParseVersion(pf.Version)
	p.version = v
	return
}

func (p *Package) MarshalJSON() ([]byte, error) {
	type PackageFile struct {
		Self      Project
		Version   string
		Timestamp time.Time
	}
	pf := PackageFile{
		Self:         p.self,
		Timestamp: p.timestamp,
		Version:      p.version.String(),
	}
	return json.Marshal(pf)
}
