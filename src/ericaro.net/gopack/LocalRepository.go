package gopack

import (
	"ericaro.net/gopack/protocol"
	. "ericaro.net/gopack/semver"
	"bytes"
	"encoding/json"
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

type LocalRepository struct {
	root    string // absolute path to the repo, this must be a filesystem writable path.
	remotes []protocol.Client
	
}

func (p LocalRepository) Write() (err error) {
	dst := filepath.Join(p.root, GpkrepositoryFile)
	fmt.Printf("writing down the local repo metadata to %s\n", dst)
	err = JsonWriteFile(dst, &p)
	return err
}


func NewLocalRepository(root string) (r *LocalRepository, err error) {
	root, err = filepath.Abs(filepath.Clean(root))
	dst := filepath.Join(root, GpkrepositoryFile)
	if err != nil {
		return
	}

	r = &LocalRepository{
		root:    root,
		remotes: make([]protocol.Client, 0),
		
	}
	err = JsonReadFile(dst, r)
	return
}

func (r *LocalRepository) Remotes() []protocol.Client {
	return r.remotes
}



func (r *LocalRepository) Remote(name string) (remote protocol.Client, err error) {
	for _, re := range r.remotes {
		if strings.EqualFold(name, re.Name()) {
			return re, nil
		}
	}
	return nil, errors.New("Missing remote")
}

func (p *LocalRepository) RemoteAdd(remote protocol.Client) (err error) {
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
			tmp := make([]protocol.Client, 0, len(p.remotes))
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
	fmt.Printf("new package %#v\n", p)
	p.Write()
	return
}

// search for package starting with name, and return them
func (r *LocalRepository) Search(search string, start int) (result []protocol.PID) {
	//fmt.Printf("q: %s start=%d\n", search, start)
	sp := filepath.Join(r.root, search)
	sd := filepath.Dir(sp)
	base := filepath.Base(sp)
	M := 10
	results := make([]protocol.PID, M)
	i := 0
	handler := func(srcpath string) bool {
		if i >= start {
			path, _ := filepath.Rel(r.root, srcpath)
			pack := filepath.Dir(path)
			v, _ := ParseVersion(filepath.Base(path))
			results[i-start] = protocol.PID{
				Name:    pack,
				Version: v,
			}
		}
		i++
		return i-start < M
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
func (r *LocalRepository) findProjectDependencies(p *Project, remotes []protocol.Client, dependencies map[ProjectID]*Package, dependencyList *[]*Package, offline, update bool) (err error) {
	for _, d := range p.dependencies {
		if dependencies[d] == nil { // it's a new dependencies
			prj, err := r.FindPackage(d)
			if !offline {
				if err != nil { // missing dependency in local repo, search remote

					prj, err = remoteHandler(remotes, func(remote protocol.Client, suc chan *Package, fail chan error) (p *Package, err error) {
						prj, err = r.DownloadPackage(remote, d)
						fail <- err
						if err == nil {
							suc <- prj
						}
						return
					})

				} else if update {
					// try to get a newer version into prjnew
					prjnew, err := remoteHandler(remotes, func(remote protocol.Client, suc chan *Package, fail chan error) (p *Package, err error) {
						prj, err = r.DownloadPackage(remote, d) // always try to download updates, if there is no update it fails fast
						fail <- err
						if err == nil {
							suc <- prj
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
func remoteHandler(remotes []protocol.Client, handler func(r protocol.Client, success chan *Package, failure chan error) (p *Package, err error)) (p *Package, err error) {
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

func (r *LocalRepository) DownloadPackage(remote protocol.Client, p ProjectID) (prj *Package, err error) {

	reader, err := remote.Fetch(protocol.PID{
		Name:    p.Name(),
		Version: p.Version(),
	})

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

func (p *LocalRepository) UnmarshalJSON(data []byte) (err error) {
	type RemoteFile struct {
		Name string
		Url  string
		Token string
	}

	type LocalRepositoryFile struct {
		Remotes []RemoteFile
	}
	var pf LocalRepositoryFile
	json.Unmarshal(data, &pf)

	for _, r := range pf.Remotes {
		ur, err := url.Parse(r.Url)
		if err != nil {return err}
		token, err := protocol.ParseStdToken(r.Token)
		if err != nil {return err}
		client, err := protocol.NewClient(r.Name, *ur, token )
		if err != nil {continue}
		p.RemoteAdd(client)
	}
	return
}

func (p *LocalRepository) MarshalJSON() ([]byte, error) {
	type RemoteFile struct {
		Name string
		Url  string
		Token string
	}

	type LocalRepositoryFile struct {
		Remotes []RemoteFile
	}
	
	pf := LocalRepositoryFile{
		Remotes:  make([]RemoteFile, len(p.remotes) ),
	}
	for i:= range p.remotes {
		pr := p.remotes[i]
		u := pr.Path()
		pf.Remotes[i] = RemoteFile{
			Name : pr.Name(),
			Url : u.String(),
		}
		tok := pr.Token()
		if tok != nil {
			pf.Remotes[i].Token = tok.FormatStd()
		}
	}
	return json.Marshal(pf)
}

// add listing capacities (list current version for a given package) 

// resolve dependencies? include local search and repo tree search, not local info)
