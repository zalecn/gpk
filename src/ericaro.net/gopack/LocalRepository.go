package gopack

import (
	"bytes"
	"encoding/json"
	"ericaro.net/gopack/protocol"
	. "ericaro.net/gopack/semver"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GpkrepositoryFile        = ".gpkrepository"
	GpkRepositoryFileVersion = "1.0.0"
)

//LocalRepository centralize operations around a directory (root), and a slice of remotes
type LocalRepository struct {
	root    string // absolute path to the repo, this must be a filesystem writable path.
	remotes []protocol.Client
}

//Write persists the LocalRepository information into it (as a .gpkrepository file)
func (p LocalRepository) Write() (err error) {
	dst := filepath.Join(p.root, GpkrepositoryFile)
	err = JsonWriteFile(dst, &p)
	return err
}

//NewLocalRepository creates a new Empty LocalRepository on the current dir, does not overwrite the actual contect, but read it
func NewLocalRepository(root string) (r *LocalRepository, err error) {
	root, err = filepath.Abs(filepath.Clean(root))
	if err != nil {
		return
	}
	
	
	_, err = os.Stat(root)
	if os.IsNotExist(err) { 
		os.MkdirAll(root, os.ModeDir|os.ModePerm) // mkdir -p
		err = nil
	}
	

	dst := filepath.Join(root, GpkrepositoryFile)
	r = &LocalRepository{
		root:    root,
		remotes: make([]protocol.Client, 0),
	}
	JsonReadFile(dst, r) // those errors are escaped
	return
}

//Remotes is to get the current slice of remotes
func (r *LocalRepository) Remotes() []protocol.Client {
	return r.remotes
}

//Remote to get a remote by name
func (r *LocalRepository) Remote(name string) (remote protocol.Client, err error) {
	for _, re := range r.remotes {
		if strings.EqualFold(name, re.Name()) {
			return re, nil
		}
	}
	return nil, errors.New("Missing remote")
}

//RemoteAdd append a remote to the list, refuse to append if there is a remote with that name already
func (p *LocalRepository) RemoteAdd(remote protocol.Client) (err error) {
	for _, r := range p.remotes[:]  { // operate on a copy of the remotes
		if strings.EqualFold(remote.Name(), r.Name()) {
			p.RemoteRemove(remote.Name() )
			u := r.Path()
			SuccessStyle.Printf("       -%s %s\n", remote.Name(), u.String())
			//return errors.New(fmt.Sprintf("A Remote called %s already exists", remote.Name()))
		}
	}
	p.remotes = append(p.remotes, remote)
	return
}

//RemoteRemove remove a remote from the list. Cannot fail. If there is no such remote it exit silently.
func (p *LocalRepository) RemoteRemove(name string) (ref protocol.Client, err error) {
	for i, r := range p.remotes {
		if strings.EqualFold(name, r.Name()) {
			ref = r
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

//FindPackage Look for the package identified by its PID within this local repository
func (r *LocalRepository) FindPackage(p ProjectID) (prj *Package, err error) {
	relative := p.Path()
	abs := filepath.Join(r.root, relative, GpkFile)
	//log.Printf("Looking for %v into %v", p, abs)
	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("Package %s %s is missing.", p.Name(), p.Version().String()))
	} else {
		prj, err = ReadPackageFile(abs)
	}
	return // nil possible
}

//InstallProject Creates a Package for this project, and the provided version. Copy the project content into this local repository
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
	if !os.IsNotExist(err) { // also check for the local policy
		os.RemoveAll(dst)
	}
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p

	//prepare recursive handlers
	dirHandler := func(ldst, lsrc string) (err error) {
		err = os.MkdirAll(ldst, os.ModeDir|os.ModePerm) // mkdir -p
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		p := filepath.Dir(ldst)
		os.MkdirAll(p, os.ModeDir|os.ModePerm) // mkdir -p
		_, err = CopyFile(ldst, lsrc)
		return
	}
	//makes the copy
	p.self.ScanProjectSrc(dst, dirHandler, fileHandler)
	p.self.ScanBinPlatforms(dst, fileHandler)
	//walkDir(filepath.Join(dst, "src"), filepath.Join(prj.workingDir, "src"), dirHandler, fileHandler)
	p.self.workingDir = dst
	p.Write()
	return
}

//Search for package starting with name, and return them
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
	return results[:i]
}

//ResolvePackageDependencies recursively scan a package's dependency ProjectID and tries to resolve every ProjectID into a Package object.
//ProjectID can be seen as just a pointer, or a reference, whereas Package can be seen as real content
// In the process of resolving PID -> Package there are two options:
// offline: if offline does not use any remote to look for missing dependencies
// update:  if online use remotes to search for a newest version of the package.
// Note: update option does not make sens for "released" version, as they are read only, and cannot be updated.
// but we consider that Snapshots version can be. Server might allow to push and override a snapshot version. In which case the update command make sense.  
func (r *LocalRepository) ResolvePackageDependencies(p *Package, offline, update bool) (dependencies []*Package, err error) {
	return r.ResolveDependencies(&p.self, offline, update)
}

//ResolveDependencies lookup recursively for all project dependencies
func (r *LocalRepository) ResolveDependencies(p *Project, offline, update bool) (dependencies []*Package, err error) {
	depMap := make(map[ProjectID]*Package)
	dependencies = make([]*Package, 0, 10)
	err = r.findProjectDependencies(p, r.remotes, depMap, &dependencies, offline, update)
	return
}

//findProjectDependencies private recursive version
func (r *LocalRepository) findProjectDependencies(p *Project, remotes []protocol.Client, dependencies map[ProjectID]*Package, dependencyList *[]*Package, offline, update bool) (err error) {
	for _, d := range p.dependencies {
		if dependencies[d] == nil { // it's a new dependencies
			prj, err := r.FindPackage(d)
			if !offline {
				if err != nil { // missing dependency in local repo, search remote
					log.Printf("Trying to download %s from remotes", d)
					prj, err = remoteHandler(remotes, func(remote protocol.Client, suc chan *Package, fail chan error) (p *Package, err error) {
						rprj, err := r.downloadPackage(remote, d)
						if err != nil {
							fail <- err
						} else {
							suc <- rprj
						}
						return
					})

				} else if update {
					if d.Version().IsSnapshot() {
						// try to get a newer version into prjnew
						log.Printf("Trying to download a newer version for %s", d)
						prjnew, err := remoteHandler(remotes, func(remote protocol.Client, suc chan *Package, fail chan error) (p *Package, err error) {
							rprj, err := r.downloadPackage(remote, d) // always try to download updates, if there is no update it fails fast
							if err != nil {
								fail <- err
							} else {
								suc <- rprj
							}
							return
						})
						// prjnew contains a newer version (downloaded though)
						if err != nil {
							log.Printf("Failed to download new Version for %s", d)
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
// it expect that
func remoteHandler(remotes []protocol.Client, handler func(r protocol.Client, success chan *Package, failure chan error) (p *Package, err error)) (pak *Package, err error) {
	success := make(chan *Package)
	failure := make(chan error)

	for _, rem := range remotes {
		go handler(rem, success, failure)
	}
	// now collect results
	for _ = range remotes {
		select {
		case pak = <-success:
			err = nil
			//close(success) // this is forbidden but I would like a way to leave with the first result asap
		case err = <-failure:
			// nothing to do // the loop will end after remotes queries
		}
	}
	return
}

//downloadPackage fetch the package, and install it in the local repository
func (r *LocalRepository) downloadPackage(remote protocol.Client, p ProjectID) (prj *Package, err error) {
	log.Printf("Downloading %s from %s", p, remote.Name())
	reader, err := remote.Fetch(protocol.PID{
		Name:    p.Name(),
		Version: p.Version(),
	})
	if err != nil {
		return nil, err
	}
	prj, err = r.Install(reader)
	return
}

//Install read a package in the reader (a tar.gzed stream, with a package .gpk inside and the project content)
// find a suitable place for it ( name/version ) and replace the content
func (r *LocalRepository) Install(reader io.Reader) (prj *Package, err error) {
	return r.install(true, reader)
}

func (r *LocalRepository) InstallAppend(reader io.Reader) (prj *Package, err error) {
	return r.install(false, reader)
}

func (r *LocalRepository) install(clean bool, reader io.Reader) (prj *Package, err error) {
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader) // download the tar.gz
	//reader.Close()
	if err != nil {
		return
	}
	mem := bytes.NewReader(buf.Bytes())
	prj, err = ReadPackageInPackage(mem) // foretell the package object from within a buffer
	if err != nil {
		//log.Printf("Cannot read package", err)
		return
	}
	dst := filepath.Join(r.root, prj.self.name, prj.version.String())
	_, err = os.Stat(dst)
	if clean && !os.IsNotExist(err) { 
		os.RemoveAll(dst)
	}
	err = os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	if err != nil {
		log.Printf("Cannot install package %s", err)
		return
	}

	prj.self.workingDir = dst
	mem = bytes.NewReader(buf.Bytes())
	err = prj.Unpack(mem) // now I know the target I can unpack it.
	//	if err != nil {
	//		log.Printf("Installed with err ", err)
	//	}
	return

}

//GoPath computes a GOPATH string based on a slice of Packages (use os.PathListSeparator as separator)
func (r *LocalRepository) GoPath(dependencies []*Package) (gopath string, err error) {
	sources := make([]string, 0, len(dependencies))
	for _, pr := range dependencies {
		sources = append(sources, pr.self.workingDir) // here if you are smart you can build a gopath on a snapshot dependency ;-) for real
	}
	gopath = strings.Join(sources, string(os.PathListSeparator))
	return
}

//Root returns this repo roots
func (r LocalRepository) Root() string {
	return r.root
}

//UnmarshalJSON is part of the json protocol to make this object read/writable in json
func (p *LocalRepository) UnmarshalJSON(data []byte) (err error) {
	type RemoteFile struct {
		Name  string
		Url   string
		Token string
	}

	type LocalRepositoryFile struct {
		FormatVersion string
		Remotes       []RemoteFile
	}
	var pf LocalRepositoryFile
	json.Unmarshal(data, &pf)
	if pf.FormatVersion != GpkRepositoryFileVersion {
		log.Printf("Warning: Unknown format version \"%s\"", pf.FormatVersion)
	}
	for _, r := range pf.Remotes {
		ur, err := url.Parse(r.Url)
		if err != nil {
			return err
		}
		token, err := protocol.ParseStdToken(r.Token)
		if err != nil {
			return err
		}
		client, err := protocol.NewClient(r.Name, *ur, token)
		if err != nil {
			continue
		}
		p.RemoteAdd(client)
	}
	return
}

func (p *LocalRepository) MarshalJSON() ([]byte, error) {
	type RemoteFile struct {
		Name  string
		Url   string
		Token string
	}

	type LocalRepositoryFile struct {
		FormatVersion string
		Remotes       []RemoteFile
	}

	pf := LocalRepositoryFile{
		FormatVersion: GpkRepositoryFileVersion,
		Remotes:       make([]RemoteFile, len(p.remotes)),
	}
	for i := range p.remotes {
		pr := p.remotes[i]
		u := pr.Path()
		pf.Remotes[i] = RemoteFile{
			Name: pr.Name(),
			Url:  u.String(),
		}
		tok := pr.Token()
		if tok != nil {
			pf.Remotes[i].Token = tok.FormatStd()
		}
	}
	return json.Marshal(pf)
}
