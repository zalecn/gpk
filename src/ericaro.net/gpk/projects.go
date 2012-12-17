package gpk

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type ProjectID struct {
	name    string // any valid package name
	version Version
}

type Project struct {
	workingDir   string      `json:"-"` // transient workding directory aboslute path
	name         string      // package name
	dependencies []ProjectID // contains the current project's dependencies
	// TO be added build time , and test dependencies
}
type Package struct {
	self    Project
	version Version

	// more to come, like sha1,signature, snapshot/release
	// add also go1 , i.e the target go runtime.

}

//ReadPackageFile local info from the specified gopackage file
func ReadPackageFile(gpkPath string) (p *Package, err error) {
	p = &Package{}
	f, err := os.Open(gpkPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(p)
	p.self.workingDir, _ = filepath.Abs(path.Dir(gpkPath))
	return
}

//ReadProjectFile local info from the specified gopackage file
func ReadProjectFile(gpkPath string) (p *Project, err error) {
	p = &Project{}
	f, err := os.Open(gpkPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(p)
	p.workingDir, _ = filepath.Abs(path.Dir(gpkPath))
	return
}

//Write package  info to where it belongs (package holds working dir info)
func (p *Package) Write() (err error) {
	dst := filepath.Join(p.self.workingDir, GpkFile)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(p)
	return err
}

//Write project  info to where it belongs (project holds working dir info)
func (p *Project) Write() (err error) {
	dst := filepath.Join(p.workingDir, GpkFile)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(p)
	return err
}

//Path converts this project reference into the path it should have in the repository layout
func (d ProjectID) Path() string {
	return filepath.Join(d.name, d.version.String())
}

//Path converts this project reference into the path it should have in the repository layout
func (p *Package) Path() string {
	return filepath.Join(p.self.name, p.version.String())
}

//ReadProjectTar reads the .gopackage file within the tar in memory. It does not set the Root
func ReadPackageInPackage(in io.Reader) (p *Package, err error) {
	gz, err := gzip.NewReader(in)
	if err != nil {
		return
	}
	tr := tar.NewReader(in)
	defer gz.Close()
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				err = errors.New(fmt.Sprintf("Invalid package format, %v is missing", GpkFile))
			}
		}
		if hdr.Name == GpkFile {
			p = &Package{}
			err = json.NewDecoder(tr).Decode(p)
			break
		}
	}
	return
}

//Untar reads the .gopackage file within the tar in memory. It does not set the Root
func (p *Package) Unpack(in io.Reader) (err error) {
	gz, err := gzip.NewReader(in)
	tr := tar.NewReader(in)
	defer gz.Close()
	dst := p.self.workingDir
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			break
		}
		// make the target file
		ndst := filepath.Join(dst, hdr.Name)
		df, err := os.Create(ndst)
		if err != nil {
			break
		}
		io.Copy(df, tr)
		df.Close()
	}
	return
}
