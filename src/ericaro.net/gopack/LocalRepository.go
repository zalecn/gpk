package gopack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GpkrepositoryFile = ".gpkrepository"
)

//RemoteRepository is any code that can act as a remote. usually a project has a chain of remote repository where to look
type RemoteRepository interface {
	CheckPackageCanPush(p *Package) (newer bool, err error)
	CheckPackageUpdate(p *Package) (newer bool, err error)
	ReadPackage(p ProjectID) (r io.Reader, err error)
	UploadPackage(p *Package) (err error)
	SearchPackage(search string) ([]string)
	Name() string
	Path() url.URL
	// TODO provide some "reader" from the remote, so local can copy it down
}

type RemoteConstructor func(name string, u url.URL) RemoteRepository

var RemoteRepositoryFactory map[string]RemoteConstructor // factory
func RegisterRemoteRepositoryFactory(urlprotocol string, xtor RemoteConstructor) {
	if RemoteRepositoryFactory == nil {
		RemoteRepositoryFactory = make(map[string]RemoteConstructor)
	}
	if _, ok := RemoteRepositoryFactory[urlprotocol]; ok {
		panic("double remote repository definition for " + urlprotocol + "\n")
	}
	RemoteRepositoryFactory[urlprotocol] = xtor
}

func NewRemoteRepository(name string, u url.URL) RemoteRepository {
	//fmt.Printf("new remote %s %v. scheme factory = %s\n", name, u.String(), RemoteRepositoryFactory[u.Scheme])
	return RemoteRepositoryFactory[u.Scheme](name, u)
}

type LocalRepository struct {
	root    string // absolute path to the repo, this must be a filesystem writable path.
	remotes []RemoteRepository
}

func (p LocalRepository) Write() (err error) {
	dst := filepath.Join(p.root, GpkrepositoryFile)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	err = EncodeLocalRepository(f, p)
	return err
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
		root:    root,
		remotes: []RemoteRepository{Central},
	}

	f, err2 := os.Open(filepath.Join(root, GpkrepositoryFile))
	if err2 != nil {
		return
	}
	defer f.Close()
	err = DecodeLocalRepository(f, r)
	return
}

func (r *LocalRepository) Remotes() []RemoteRepository {
	return r.remotes
}

func (r *LocalRepository) Remote(name string) (remote RemoteRepository, err error) {
	for _, re := range r.remotes {
		if strings.EqualFold(name, re.Name()) {
			return re, nil
		}
	}
	return nil, errors.New("Missing remote")
}

func (p *LocalRepository) RemoteAdd(remote RemoteRepository) (err error) {
	for _, r := range p.remotes {
		if strings.EqualFold(remote.Name(), r.Name()) {
			return errors.New(fmt.Sprintf("A Remote called %s already exists", remote.Name()))
		}
	}
	p.remotes = append(p.remotes, remote)
	return
}

func (p *LocalRepository) RemoteRemove(name string) (err error) {
	for i, r := range p.remotes {
		if strings.EqualFold(name, r.Name()) {
			tmp := make([]RemoteRepository, 0, len(p.remotes))
			if i > 0 {
				tmp = append(tmp, p.remotes[0:i]...)
			}
			if i+1 < len(p.remotes) {
				tmp = append(tmp, p.remotes[i+1:]...)
			}
			p.remotes = tmp
			return
		}
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
	//	fmt.Printf("Installing to %s %s %s\n", r.root, p.Path(), dst)
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

// search for package starting with name, and return them
func (r *LocalRepository) SearchPackage(search string) (result []string) {
	sp := filepath.Join(r.root, search)
	sd := filepath.Dir(sp)
	base := filepath.Base(sp)
	
	results := make([]string, 100)
	i:=0
	handler := func(srcpath string) {
		//fmt.Printf("found  %s %d\n", srcpath, i)
		results[i], _= filepath.Rel(r.root, srcpath)
		i++
	}
	PackageWalker(sd, base, handler)
	//fmt.Printf("found  %s\n", results[:i])
	return results[:i]
}

func (r *LocalRepository) ResolvePackageDependencies(p *Package, offline, update bool) (dependencies []*Package, err error) {
	return r.ResolveDependencies(&p.self, offline, update)
}

//FindProjectDependencies lookup recursively for all project dependencies
func (r *LocalRepository) ResolveDependencies(p *Project, offline, update bool) (dependencies []*Package, err error) {
	depMap := make(map[ProjectID]*Package)
	dependencies = make([]*Package, 0, 10)

	err = r.findProjectDependencies(p, r.remotes, depMap, &dependencies, offline, update)
	if err != nil {
		return
	}
	return
}

//findProjectDependencies private recursive version
func (r *LocalRepository) findProjectDependencies(p *Project, remotes []RemoteRepository, dependencies map[ProjectID]*Package, dependencyList *[]*Package, offline, update bool) (err error) {
	for _, d := range p.dependencies {
		if dependencies[d] == nil { // it's a new dependencies
			prj, err := r.FindPackage(d)
			if !offline {
				if err != nil { // missing dependency in local repo, search remote

					prj, err = remoteHandler(remotes, func(remote RemoteRepository, suc chan *Package, fail chan error) (p *Package, err error) {
						prj, err = r.DownloadPackage(remote, d)
						fail <- err
						if err == nil {
							suc <- prj
						}
						return
					})

				} else if update {
					// try to get a newer version into prjnew
					prjnew, err := remoteHandler(remotes, func(remote RemoteRepository, suc chan *Package, fail chan error) (p *Package, err error) {
						if newer, _ := remote.CheckPackageUpdate(prj); newer {
							prj, err = r.DownloadPackage(remote, d)
							fail <- err
							if err == nil {
								suc <- prj
							}
						}
						return
					})
					// prjnew contains a newer version (downloaded though)
					if err != nil {
						return err // failed to download
					}
					prj = prjnew
				}
			}
			if prj == nil {
				return errors.New(fmt.Sprintf("Missing dependency: %v\n", d))
			}
			dependencies[d] = prj
			*dependencyList = append(*dependencyList, prj)
			err = r.findProjectDependencies(&prj.self, remotes, dependencies, dependencyList, offline, update)
			if err != nil {
				return err
			}
		}
	}
	return
}

//remoteHandler process a function for every remote
func remoteHandler(remotes []RemoteRepository, handler func(r RemoteRepository, success chan *Package, failure chan error) (p *Package, err error)) (p *Package, err error) {
	success := make(chan *Package)
	failure := make(chan error)

	for _, rem := range remotes {
		go handler(rem, success, failure)
	}
	// now collect results
	for _ = range remotes {
		select {
		case p = <-success:
			err = nil
			close(success)
			break
		case err = <-failure:
			// nothing to do // the loop will end after remotes queries
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

func (r LocalRepository) Root() string {
	return r.root
}

// add listing capacities (list current version for a given package) 

// resolve dependencies? include local search and repo tree search, not local info)
