package gpk

import (
	"bytes"
	"io"
)

//FileRemoteRepository act as a remote repository for a 
type FileRemoteRepository struct {
	repo LocalRepository // contains a local repo
}





func (r *FileRemoteRepository) 	CheckPackageUpdate(p *Package) (newer bool, err error) {
	// cave at p is the local package, I need to check for the same version in this one
	
	rp, err := r.repo.FindPackage(p.ID() )
	if err != nil {
		newer =false
	} else {
		newer = rp.timestamp.After(p.timestamp)
	}
	return
}


func (r *FileRemoteRepository) 	ReadPackage(p ProjectID) (reader io.Reader, err error) {
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
