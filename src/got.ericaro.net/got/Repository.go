package got

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultRepository = ".gotepository"
	Release           = "Release"
	Snapshot          = "Snapshot"
)

//a Repository is a directory where dependencies are stored
// they are splitted into releases, and snapshots
type Repository struct {
	Root       string // absolute path to the repository root
	ServerHost string
}

func NewDefaultRepository() (r *Repository, err error) {
	u, _ := user.Current()
	return NewRepository(filepath.Join(u.HomeDir, DefaultRepository))
}
func NewRepository(root string) (r *Repository, err error) {
	root, err = filepath.Abs(filepath.Clean(root))

	if err != nil {
		return
	}

	r = &Repository{
		Root:       root,
		ServerHost: GotCentral,
	}
	return
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (r *Repository) findLocalProject(mode string, p ProjectReference) (prj *Project, err error) {
	// TODO hande offline and update here
	
	relative := p.Path()
	abs := filepath.Join(r.Root, mode, relative, GotFile)
	//log.Printf("Looking for %v into %v", p, abs)
	if exists(abs) {
		prj, err = ReadProjectFile(abs)
	} else {
		err = errors.New(fmt.Sprintf("Dependency %v is missing.", p))
	}
	return // nil possible
}

//FindProject lookup for the project reference in the local repository
// if the repository is in snapshot mode it looks for a snapshot version first.
// if it fails or if is not in snapshot mode it looks for a release version.
func (r *Repository) FindProject(p ProjectReference, searchSnapshot, offline, update bool) (prj *Project, err error) {
	log.Printf("lookup for %v ", p)
	log.Printf("    in Release\n")
	prj, err = r.findLocalProject(Release, p)
	if searchSnapshot && prj == nil { // search for snapshot if and only if needed
		log.Printf("    in Snapshot ...")
		prj, err = r.findLocalProject(Snapshot, p)
	}
	if !offline { // go for the central server
		if prj == nil {
			log.Printf("    in Remote ...")
			prj, err = r.DownloadProject(p, searchSnapshot) // check for it on the web
		} else if update { // prj has been found, but I want to check for update
			// there are some sub cases where I DO  want to check for update
			if searchSnapshot && *prj.Snapshot {
				log.Printf("  check Update ...")
				if newer, _ := r.CheckNewerProject(prj); newer {
					// there is a new project
					log.Printf("  dl Update %v ...", p)
					prjnew, err := r.DownloadProject(p, searchSnapshot)
					if err != nil {
						return nil, err // failed to download
					}
					prj = prjnew
				}
			}
		}
	}

	if prj == nil {
		err = errors.New(fmt.Sprintf("Missing dependency %v.\nCaused by:%v", p, err))
	} else {
	log.Printf(" found\n")
	}
	return
}

//FindProjectDependencies lookup recursively for all project dependencies
func (r *Repository) FindProjectDependencies(p *Project, searchSnapshot, offline, update bool) (dependencies []*Project, err error) {
	depMap := make(map[ProjectReference]*Project)
	err = r.findProjectDependencies(p, depMap, searchSnapshot, offline, update)
	if err != nil {
		return
	}
	dependencies = make([]*Project, 0, len(depMap))
	for _, v := range depMap {
		dependencies = append(dependencies, v)
	}
	return
}

//findProjectDependencies private recursive version
func (r *Repository) findProjectDependencies(p *Project, dependencies map[ProjectReference]*Project, searchSnapshot, offline, update bool) (err error) {
	for _, d := range p.Dependencies {
		if dependencies[d] == nil { // it's a new dependencies
			prj, err := r.FindProject(d, searchSnapshot, offline, update)
			if err != nil {
				// missing dependency
				return err
			}
			dependencies[d] = prj
			err = r.findProjectDependencies(prj, dependencies, searchSnapshot, offline, update)
			if err != nil {
				return err
			}
		}
	}
	return
}

//UploadProject upload a project to the central server.
// the optional parameter snapshot, and version must be set
func (r *Repository) UploadProject(p *Project) (err error) {

	if p.Snapshot == nil || p.Version == nil {
		return errors.New("Illegal project state. It must be fully qualified: Snaphost and Version field must be defined")
	}

	// package it in memory
	buf := new(bytes.Buffer)
	p.PackageProject(buf)
	// prepare central server query args

	v := url.Values{}
	v.Set("g", p.Group)
	v.Set("a", p.Artifact)
	v.Set("v", p.Version.String())
	if *p.Snapshot {
		v.Set("r", "false")
	} else {
		v.Set("r", "true")
	}
	v.Set("t", p.Version.Timestamp.Format(time.ANSIC))

	//query url
	u := url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Scheme:   "http",
		Host:     r.ServerHost,
		Path:     "/p/ul",
		RawQuery: v.Encode(),
	}
	var client http.Client
	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		fmt.Printf("invalid request %v\n", err)
		return
	}
	req.ContentLength = int64(buf.Len())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New(fmt.Sprintf("http upload failed %d: %v", resp.StatusCode, resp.Status))
	}
	return
}

//UploadProject upload a project to the central server.
// the optional parameter snapshot, and version must be set
func (r *Repository) CheckNewerProject(p *Project) (yes bool, err error) {

	if p.Version == nil {
		return false, errors.New("Illegal project state. It must be fully qualified: Version field must be defined")
	}
	v := url.Values{}
	v.Set("g", p.Group)
	v.Set("a", p.Artifact)
	v.Set("v", p.Version.String())
	v.Set("t", p.Version.Timestamp.Format(time.ANSIC))

	//query url
	u := url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Scheme:   "http",
		Host:     r.ServerHost,
		Path:     "/p/nl",
		RawQuery: v.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	if resp.StatusCode != 200 || resp.StatusCode != 404 {
		err = errors.New(fmt.Sprintf("http query failed %d: %v", resp.StatusCode, resp.Status))
	}
	return resp.StatusCode == 200, err
}

func (r *Repository) DownloadProject(p ProjectReference, searchSnapshot bool) (prj *Project, err error) {

	// prepare central server query args
	v := url.Values{}
	v.Set("g", p.Group)
	v.Set("a", p.Artifact)
	v.Set("v", p.Version.String())
	if !searchSnapshot {
		v.Set("r", "true")
	}else {
		v.Set("r", "false")
	}
	//query url
	u := url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Scheme:   "http",
		Host:     r.ServerHost,
		Path:     "/p/dl",
		RawQuery: v.Encode(),
	}
	log.Printf("get %v\n",u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body) // download the tar.gz
	resp.Body.Close()
	if err != nil {
		return
	}

	// reads the project in memory
	prj, err = ReadProjectInPackage(buf)
	if err != nil {
		return
	}

	// computes the project relative path 
	prjPath := filepath.Join(prj.Group, prj.Artifact, prj.Version.Path())
	mode := Snapshot
	if prj.Snapshot == nil {
		err = errors.New("Invalid packaged project format: does not define the snapshot attribute")
		return
	}
	if !*prj.Snapshot {
		mode = Release
	}

	prj.Root = filepath.Join(r.Root, mode, prjPath)
	prj.UnpackageProject(buf)
	return

}
func (r *Repository) GoGetInstall(pack string) {
	prj := NewGoGetProjectReference(pack, ParseVersionReference("bigbang-0.0.0.0") )
	fmt.Printf("getting %v as %s \n", pack, prj)
	dst := filepath.Join(r.Root, Snapshot, prj.Group, prj.Artifact, prj.Version.Path() )
	if exists(dst) {
		os.RemoveAll(dst)
	}
	os.MkdirAll(dst, os.ModeDir|os.ModePerm) // mkdir -p
	g := NewGoEnv(dst)
	g.Get(pack)

	// computes the absolute path
		
}

//InstallProject install the project into this repository
func (r *Repository) InstallProject(prj *Project, v Version, snapshotMode bool) {

	p := *prj // copy the prj
	p.Snapshot = &snapshotMode
	p.Version = &v

	var mode string
	switch snapshotMode {
	case false:
		mode = Release
	default:
		mode = Snapshot
	}

	// computes the project relative path 
	prjPath := filepath.Join(p.Group, p.Artifact, v.Path())

	// computes the absolute path
	dst := filepath.Join(r.Root, mode, prjPath)
	if exists(dst) {
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

	walkDir(filepath.Join(dst, "src"), filepath.Join(p.Root, "src"), dirHandler, fileHandler)
	WriteProjectFile(filepath.Join(dst, GotFile), &p)

	if !snapshotMode { // if in release, then remove all the snapshots 
		altDir := filepath.Join(r.Root, Snapshot, prjPath)
		if exists(altDir) { // uninstall from alt
			os.RemoveAll(altDir)
		}
	}
}

func (r *Repository) GoPath(dependencies []*Project) (gopath string, err error) {
	sources := make([]string, 0, len(dependencies))
	for _, pr := range dependencies {
		sources = append(sources, pr.Root)
	}
	gopath = strings.Join(sources, string(os.PathListSeparator))
	return
}
