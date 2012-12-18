package gpk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"net/url"
)

const (
	GpkFile = ".gpk"
)

//RemoteRepository is any code that can act as a remote. usually a project has a chain of remote repository where to look
type RemoteRepository interface {
	CheckPackageUpdate(p *Package) (newer bool, err error)
	ReadPackage(p ProjectID) (r io.Reader, err error)
	UploadPackage(p *Package) (err error)
	// TODO provide some "reader" from the remote, so local can copy it down
}

type RemoteConstructor func(u url.URL) RemoteRepository

var RemoteRepositoryFactory map[string]RemoteConstructor // factory
func RegisterRemoteRepositoryFactory(urlprotocol string , xtor RemoteConstructor){
	if RemoteRepositoryFactory == nil {
		RemoteRepositoryFactory = make(map[string]RemoteConstructor)
	}
	if _,ok := RemoteRepositoryFactory[urlprotocol]; ok {
		panic("double remote repository definition for "+urlprotocol+"\n")
	}
	RemoteRepositoryFactory[urlprotocol] = xtor
}

func NewRemoteRepository(u url.URL) RemoteRepository {
	return RemoteRepositoryFactory[u.Scheme](u)
}


type LocalRepository struct {
	root string // absolute path to the repo, this must be a filesystem writable path.
}

// features:
// xtor (based on user/project config)
// search project
//install project as package
// computes the local gopath (yes its part of the local repo, gopath always point to the local repo)
// goget compatibility

func NewLocalRepository(root string) (r *LocalRepository, err error) {
	root, err = filepath.Abs(filepath.Clean(root))
	if err != nil {
		return
	}

	r = &LocalRepository{
		root: root,
	}
	return
}

func (r *LocalRepository) FindPackage(p ProjectID) (prj *Package, err error) {
	relative := p.Path()
	abs := filepath.Join(r.root, relative, GpkFile)
	//log.Printf("Looking for %v into %v", p, abs)
	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("Dependency %v is missing.", p))
	} else {
		prj, err = ReadPackageFile(abs)
	}
	return // nil possible
}

//InstallProject install the project into this repository
func (r *LocalRepository) InstallProject(prj *Project, v Version) (p *Package) {

	p = &Package{
		self:      *prj,
		version:   v,
		timestamp: time.Now(),
	}
	// computes the project relative path 
	// computes the absolute path
	dst := filepath.Join(r.root, p.Path())
	fmt.Printf("Installing to %s %s %s\n",r.root, p.Path(),  dst)
	_, err := os.Stat(dst)
	if os.IsExist(err) { // also check for the local policy
		os.RemoveAll(dst)
	}
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p

	//prepare recursive handlers
	dirHandler := func(ldst, lsrc string) (err error) {
		err = os.MkdirAll(ldst, os.ModeDir|os.ModePerm) // mkdir -p
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		_, err = CopyFile(ldst, lsrc)
		return
	}
	//makes the copy
	walkDir(filepath.Join(dst, "src"), filepath.Join(prj.workingDir, "src"), dirHandler, fileHandler)
	p.self.workingDir = dst
	p.Write()
	return
}

//FindProjectDependencies lookup recursively for all project dependencies
func (r *LocalRepository) ResolveDependencies(p *Project, remote RemoteRepository, offline, update bool) (dependencies []*Package, err error) {
	depMap := make(map[ProjectID]*Package)
	err = r.findProjectDependencies(p, remote, depMap, offline, update)
	if err != nil {
		return
	}
	dependencies = make([]*Package, 0, len(depMap))
	for _, v := range depMap {
		dependencies = append(dependencies, v)
	}
	return
}

//findProjectDependencies private recursive version
func (r *LocalRepository) findProjectDependencies(p *Project, remote RemoteRepository, dependencies map[ProjectID]*Package, offline, update bool) (err error) {
	for _, d := range p.dependencies {
		if dependencies[d] == nil { // it's a new dependencies
			prj, err := r.FindPackage(d)
			if !offline {
				if err != nil { // missing dependency in local repo, search remote
					prj, err = r.DownloadPackage(remote, d)
				} else if update {
					if newer, _ := remote.CheckPackageUpdate(prj); newer {
						// there is a new project
						prjnew, err := r.DownloadPackage(remote, d)
						if err != nil {
							return err // failed to download
						}
						prj = prjnew
					}
				}
			}
			if prj == nil {
				return errors.New(fmt.Sprintf("Missing dependency: %v\n", d))
			}
			dependencies[d] = prj
			err = r.findProjectDependencies(&prj.self, remote, dependencies, offline, update)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (r *LocalRepository) DownloadPackage(remote RemoteRepository, p ProjectID) (prj *Package, err error) {
	reader, err := remote.ReadPackage(p)
	if err != nil {
		return nil, err
	}
	return r.Install(reader)
}

func (r *LocalRepository) Install(reader io.Reader) (prj *Package, err error) {
	fmt.Printf("installing ...\n")
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader) // download the tar.gz
	//reader.Close()
	if err != nil {
		return
	}
	mem := bytes.NewReader(buf.Bytes())
	prj, err = ReadPackageInPackage(mem) // foretell the package object from within a buffer
	fmt.Printf("                %v\n", prj)
	if err != nil {
		return
	}
	prj.self.workingDir = filepath.Join(r.root, prj.self.name, prj.version.String())
	fmt.Printf("                              TO %v\n", prj.self.workingDir)
	
	mem = bytes.NewReader(buf.Bytes())
	err = prj.Unpack(mem) // now I know the target I can unpack it.
	
	return

}

func (r *LocalRepository) GoPath(dependencies []*Package) (gopath string, err error) {
	sources := make([]string, 0, len(dependencies))
	for _, pr := range dependencies {
		sources = append(sources, pr.self.workingDir) // here if you are smart you can build a gopath on a snapshot dependency ;-) for real
	}
	gopath = strings.Join(sources, string(os.PathListSeparator))
	return
}

func (r LocalRepository) Root() string{
	return r.root
}

// add listing capacities (list current version for a given package) 

// resolve dependencies? include local search and repo tree search, not local info)
