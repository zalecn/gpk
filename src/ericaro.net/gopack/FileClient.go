package gopack

import (
	"bytes"
	"ericaro.net/gopack/protocol"
	"fmt"
	"io"
	"net/url"
	//"path/filepath"
)

func init() { // register this as a handler for file:/// url scheme
	protocol.RegisterClient("file", NewFileClient)
}

//FileClient implements a remote protocol using a local file 
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
		repo:       *loc,
		BaseClient: *protocol.NewBaseClient(name, u, token),
	}
	return

}

func (c *FileClient) Push(pid protocol.PID, r io.Reader) (err error) {
	//dst := filepath.Join(c.repo.Root(), pid.Path())
	_, err = c.repo.Install(r)
	return
}

func (r *FileClient) Search(query string, start int) (result []protocol.PID) {
	return r.repo.Search(query, start)
}

func (c *FileClient) Fetch(pid protocol.PID) (r io.ReadCloser, err error) {
	//ReadPackage(p ProjectID) (reader io.Reader, err error) {
	p := *NewProjectID(pid.Name, pid.Version)
	rp, err := c.repo.FindPackage(p)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	rp.Pack(buf)

	// the package has been built into the buffer
	return closeable{buf}, nil
}


// internal trick to fake a closeable, even if we use a buffer that is not closeable in fact
type closeable struct {
	*bytes.Buffer
}

func (c closeable) Close() error { return nil }
