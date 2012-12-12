package got

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
	"path/filepath"
)

const (
	DefaultRepository = ".gotepository"
	Release           = "Release"
	Snapshot          = "Snapshot"
)

//a Repository is a directory where dependencies are stored
// they are splitted into releases, and snapshots
type Repository struct {
	Root string // absolute path to the repository root
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
		Root: root,
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

func (r *Repository) findProject(mode string, p ProjectReference, offline, update bool) (prj *Project, err error) {
	// TODO hande offline and update here
	relative := p.Path()
	abs := filepath.Join(r.Root, mode, relative, GorFile)
	//log.Printf("Looking for %v into %v", p, abs)
	if exists(abs) {
		prj, err = ReadProjectFile(abs)
	} else {
		err =  errors.New(fmt.Sprintf("Missing dependency %v",p ))
	}
	return // nil possible
}

//FindProject lookup for the project reference in the local repository
// if the repository is in snapshot mode it looks for a snapshot version first.
// if it fails or if is not in snapshot mode it looks for a release version.
func (r *Repository) FindProject(p ProjectReference, searchSnapshot, offline, update bool) (prj *Project, err error) {
	
	prj, err = r.findProject(Release, p, offline, false)// never search updates on the release
		
	if searchSnapshot && err != nil {// search for snapshot if and only if needed
		prj, err = r.findProject(Snapshot, p, offline, update)
	}

	if prj == nil {
		err =  errors.New(fmt.Sprintf("Missing dependency %v.\nCaused by:%v",p, err ))
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

//InstallProject install the project into this repository
func (r *Repository) InstallProject(p *Project, v Version, snapshotMode bool) {
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
	WriteProjectFile(filepath.Join(dst, GorFile), p)

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
