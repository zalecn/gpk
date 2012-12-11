package gor

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	
)

const (
	DefaultRepository = ".gorepository"
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

func (r *Repository) findProject(mode string, p ProjectReference) (prj *Project) {
	relative := p.Path()
	abs := filepath.Join(r.Root, mode, relative, GorFile)
	log.Printf("Looking for %v into %v", p, abs)
	if exists(abs) {
		prj, _ = ReadProjectFile(abs)
	}
	return prj // nil possible
}

//FindProject lookup for the project reference in the local repository
// if the repository is in snapshot mode it looks for a snapshot version first.
// if it fails or if is not in snapshot mode it looks for a release version.
func (r *Repository) FindProject(p ProjectReference, searchSnapshot bool) (prj *Project) {
	if searchSnapshot {
		prj = r.findProject(Snapshot, p)
	}

	if prj == nil {
		prj = r.findProject(Release, p)
	}
	return
}

//FindProjectDependencies lookup recursively for all project dependencies
func (r *Repository) FindProjectDependencies(p *Project, searchSnapshot bool) (dependencies []*Project) {
	depMap := make(map[ProjectReference]*Project)
	r.findProjectDependencies(p, depMap, searchSnapshot)

	dependencies = make([]*Project, 0, len(depMap))
	for _, v := range depMap {
		dependencies = append(dependencies, v)
	}
	return
}

//findProjectDependencies private recursive version
func (r *Repository) findProjectDependencies(p *Project, dependencies map[ProjectReference]*Project, searchSnapshot bool) {
	for _, d := range p.Dependencies {
		if dependencies[d] == nil { // new dependencies
			prj := r.FindProject(d, searchSnapshot)
			if prj == nil {
				log.Fatalf("missing dependency %v \n", d)
			}
			dependencies[d] = prj
			r.findProjectDependencies(prj, dependencies, searchSnapshot)
		}
	}
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
	dirHandler := func(ldst, lsrc string)(err error){
		err = os.MkdirAll(ldst, os.ModeDir|os.ModePerm) // mkdir -p
		return
	}
	fileHandler := func(ldst, lsrc string) (err error) {
		_, err = CopyFile(ldst, lsrc)
		return
	}
	
	walkDir(filepath.Join(dst, "src"), filepath.Join(p.Root, "src") , dirHandler, fileHandler)
	WriteProjectFile(filepath.Join(dst, GorFile), p)

	if !snapshotMode { // if in release, then remove all the snapshots 
		altDir := filepath.Join(r.Root, Snapshot, prjPath)
		if exists(altDir) { // uninstall from alt
			os.RemoveAll(altDir)
		}
	}
}

