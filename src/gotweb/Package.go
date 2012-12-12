package gotweb

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"got.ericaro.net/got"
	"io"
	"net/http"
	"strconv"
)

//a Package is a pure, in-memory representation of a Package
type Package struct {
	Group, Artifact, Root, Version string
	ContentBlob                    appengine.BlobKey // tar.gz of the content
	Release                        bool              // whether its a release or not
}

type PackageServer struct {
	appengine.Context
}

func New(r *http.Request) *PackageServer {
	return &PackageServer{
		Context: appengine.NewContext(r),
	}
}

func (s *PackageServer) Put(key got.ProjectReference, p *Package) (k *datastore.Key, err error) {
	c := s.Context
	k, err = datastore.Put(c, datastore.NewKey(c, "Package", key.String(), 0, nil), p)
	return
}

func (s *PackageServer) Get(key got.ProjectReference) (p *Package, err error) {
	c := s.Context
	p = new(Package)
	err = datastore.Get(c, datastore.NewKey(c, "Package", key.String(), 0, nil), p)
	return
}

func (s *PackageServer) DeleteBlob(k appengine.BlobKey) (err error) {
	return blobstore.Delete(s.Context, k)

}
func (s *PackageServer) CreateBlob(r io.Reader) (k appengine.BlobKey, err error) {

	w, err := blobstore.Create(s.Context, "application/x-gzip")
	if err != nil {
		return k, err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return k, err
	}
	err = w.Close()
	if err != nil {
		return k, err
	}
	return w.Key()
}

//receive handles the artifact upload
func receive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}
	s := New(r)
	s.Debugf("processing receive\n")

	// identify the package
	vals := r.URL.Query()
	group := vals.Get("g")    // todo validate the syntax
	artifact := vals.Get("a") // todo validate the syntax
	version := got.ParseVersionReference(vals.Get("v"))
	release, _ := strconv.ParseBool(vals.Get("r"))

	s.Debugf("upload release=%t %v:%v:%v \n", release, group, artifact, version)

	// try to get the package if it already exists
	pr := got.NewProjectReference(group, artifact, version)
	p, err := s.Get(pr)
	if p != nil {
		s.Debugf("project already exists \n")
	}
	if p != nil && p.Release {
		http.Error(w, "Artifact already exists in Release Mode.", http.StatusMethodNotAllowed)
		return
	}
	// it's ok to create it

	blob, err := s.CreateBlob(r.Body) // create and fill the blob
	if p == nil {
		p = &Package{
			Group:    group,
			Artifact: artifact,
			Root:     version.Root,
			Version:  version.Parts,
			Release:  release,
		}
	} else {
		// delete the previous blob
		p.Release= release
		s.DeleteBlob(p.ContentBlob)
	}
	// update the new blob
	p.ContentBlob = blob
	_, err = s.Put(pr, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Debugf("stored %v\n", pr)
}

func serve(w http.ResponseWriter, r *http.Request) {
	s := New(r)
	group := r.FormValue("g")
	artifact := r.FormValue("a")
	version := got.ParseVersionReference(r.FormValue("v"))
	if group == "" || artifact == "" {
		http.NotFound(w, r)
		return
	}
	pr := got.NewProjectReference(group, artifact, version)
	s.Debugf("processing serve\n")
	p, err := s.Get(pr)
	if err == datastore.ErrNoSuchEntity {
		http.NotFound(w, r)
		return
	}
	s.Debugf("serving %#v\n", p)
	blobstore.Send(w, p.ContentBlob)
	return
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func init() {
	http.HandleFunc("/p/dl", serve)
	http.HandleFunc("/p/ul", receive)
}
