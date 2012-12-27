package gopack

import (
	"bytes"
	"ericaro.net/gopack/protocol"
	"io"
	"net/url"
	"fmt"
	//"path/filepath"
)


func init() { // register this as a handler for file:/// url scheme
	protocol.RegisterClient("file", NewFileClient)
}

//FileClient act as a remote repository for a 
type FileClient struct {
	repo LocalRepository // contains a local repo
	protocol.BaseClient
}

func NewFileClient(name string, u url.URL, token *protocol.Token) (r protocol.Client, err error) {
	loc, err := NewLocalRepository(u.Path)
	if err != nil {
		err = fmt.Errorf("Invalid Remote Repository \"%s\"\n    ↳ %s is not a valid repository.\n    ↳ %s", name, u.String(), err)
	}
	r = &FileClient{
		repo: *loc,
		BaseClient: *protocol.NewBaseClient(name, u, token),
	}
	return

}

func (c *FileClient) Push(pid protocol.PID, r io.Reader) (err error) {
	//dst := filepath.Join(c.repo.Root(), pid.Path())
	c.repo.Install(r)
	return
}

func (r *FileClient) Search(query string, start int) (result []protocol.PID) {
	return r.repo.Search(query, start)
}

func (r *FileClient) CheckPackageUpdate(p *Package) (newer bool, err error) {
	// cave at p is the local package, I need to check for the same version in this one

	rp, err := r.repo.FindPackage(p.ID())
	if err != nil {
		newer = false
	} else {
		newer = rp.timestamp.After(p.timestamp)
	}
	return
}

func (c *FileClient) Fetch(pid protocol.PID) (r io.ReadCloser, err error) {
	//ReadPackage(p ProjectID) (reader io.Reader, err error) {
	p := NewProjectID(pid.Name, pid.Version)
	rp, err := c.repo.FindPackage(p)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	rp.Pack(buf)
	
	// the package has been built into the buffer
	return Closeable{buf}, nil
}

type Closeable struct{
		*bytes.Buffer
	}
	func (c Closeable) Close() error{return nil}

// TODO provide some "reader" from the remote, so local can copy it down
