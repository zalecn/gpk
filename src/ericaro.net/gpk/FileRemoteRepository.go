package gpk

import (
	"bytes"
	"io"
	"net/url"
	"path/filepath"
	"os"
)

func init() {
	file := func(u url.URL) RemoteRepository {
		f, _ := NewFileRemoteRepository(u)
		return f
	}
	RegisterRemoteRepositoryFactory("file", file)
}

//FileRemoteRepository act as a remote repository for a 
type FileRemoteRepository struct {
	repo LocalRepository // contains a local repo
}

func NewFileRemoteRepository(u url.URL) (r *FileRemoteRepository, err error) {
	loc, err := NewLocalRepository(u.Path)
	r = &FileRemoteRepository{
		repo: *loc,
	}
	return

}

func (r *FileRemoteRepository) UploadPackage(p *Package) (err error) {
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

	return
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
