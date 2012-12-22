package gopack

import (
	"bytes"
	"io"
	"net/url"
	"path/filepath"
	"os"
)

func init() {
	file := func(name string, u url.URL) RemoteRepository {
		f, _ := NewFileRemoteRepository(name, u)
		return f
	}
	RegisterRemoteRepositoryFactory("file", file)
}

//FileRemoteRepository act as a remote repository for a 
type FileRemoteRepository struct {
	repo LocalRepository // contains a local repo
	name string
}

func NewFileRemoteRepository(name string, u url.URL) (r *FileRemoteRepository, err error) {
	loc, err := NewLocalRepository(u.Path)
	r = &FileRemoteRepository{
		repo: *loc,
		name: name,
	}
	return

}

func (r FileRemoteRepository) Name() string {return r.name}
func (r FileRemoteRepository) Path() url.URL {
	return url.URL{
	Scheme: "file",
	Path: r.repo.Root(),
	}
}


func (r *FileRemoteRepository) UploadPackage(pkg *Package) (err error) {
	p:= *pkg
	dst := filepath.Join(r.repo.Root(), p.Path())
	_, err = os.Stat(dst)
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
	walkDir(filepath.Join(dst, "src"), filepath.Join(p.InstallDir(), "src"), dirHandler, fileHandler)
	
	p.self.workingDir = dst
	p.Write()
	return
}

func (r *FileRemoteRepository) CheckPackageCanPush(p *Package) (canpush bool, err error) {
	// cave at p is the local package, I need to check for the same version in this one
	rp, err := r.repo.FindPackage(p.ID())
	if err != nil {
		canpush = true
	} else {
		canpush = rp.timestamp.Before(p.timestamp)
	}
	return
	
}

func (r *FileRemoteRepository) SearchPackage(search string) ([]string){
	return r.repo.SearchPackage(search)
}

func (r *FileRemoteRepository) CheckPackageUpdate(p *Package) (newer bool, err error) {
	// cave at p is the local package, I need to check for the same version in this one

	rp, err := r.repo.FindPackage(p.ID())
	if err != nil {
		newer = false
	} else {
		newer = rp.timestamp.After(p.timestamp)
	}
	return
}

func (r *FileRemoteRepository) ReadPackage(p ProjectID) (reader io.Reader, err error) {
	rp, err := r.repo.FindPackage(p)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	rp.Pack(buf)
	// the package has been built into the buffer
	return buf, nil
}

// TODO provide some "reader" from the remote, so local can copy it down
