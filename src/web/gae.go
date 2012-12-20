package httpserver

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"bytes"
	"ericaro.net/gpk"
	"ericaro.net/gpk/httpserver"
	"io"
	"net/http"
	"time"
)

//a Package is a pure, in-memory representation of a Package
type GaePackage struct {
	Timestamp   time.Time
	ContentBlob appengine.BlobKey // tar.gz of the content

}

type GaeBackendServer struct {
	appengine.Context
}

func gaereceive(w http.ResponseWriter, r *http.Request) {
	be := &GaeBackendServer{
		Context: appengine.NewContext(r),
	}
	httpserver.Receive(be, w, r)
}
func gaesend(w http.ResponseWriter, r *http.Request) {
	be := &GaeBackendServer{
		Context: appengine.NewContext(r),
	}
	httpserver.Send(be, w, r)
}

func gaenewer(w http.ResponseWriter, r *http.Request) {
	be := &GaeBackendServer{
		Context: appengine.NewContext(r),
	}
	httpserver.Newer(be, w, r)
}

//func (s *GaeBackendServer) Debugf(format string, args ...interface{}) {
//	
//}

//Contains return true if the server contains the ProjectID
func (s *GaeBackendServer) Receive(id gpk.ProjectID, timestamp time.Time, w http.ResponseWriter, r *http.Request) (err error) {

	c := s.Context
	writer, err := blobstore.Create(s.Context, "application/x-gzip")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err != nil {
		return
	}
	_, err = io.Copy(buf, r.Body) // store it in memory
	if err != nil {
		return
	}

	pack, err := gpk.ReadPackageInPackage(buf)
	if err != nil {
		return
	}
	_, err = io.Copy(writer, buf) // write the blob back
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	blobkey, err := writer.Key()
	if err != nil {
		return
	}

	// build the entity

	p := &GaePackage{
		Timestamp : pack.Timestamp(),
		ContentBlob: blobkey,
	}

	_, err = datastore.Put(c, datastore.NewKey(c, "GaePackage", id.ID(), 0, nil), p)
	return

}

func (s *GaeBackendServer) Send(id gpk.ProjectID, w http.ResponseWriter, r *http.Request) {
	c := s.Context
	p := new(GaePackage)
	err := datastore.Get(c, datastore.NewKey(c, "GaePackage", id.ID(), 0, nil), p)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	blobstore.Send(w, p.ContentBlob)
	return
}

func (s *GaeBackendServer) Newer(id gpk.ProjectID, timestamp time.Time, w http.ResponseWriter, r *http.Request) {
	c := s.Context
	p := new(GaePackage)
	err := datastore.Get(c, datastore.NewKey(c, "GaePackage", id.ID(), 0, nil), p)

	if err != nil || p == nil || !p.Timestamp.After(timestamp) {
		http.NotFound(w, r)
		return
	}
}

func (s *GaeBackendServer) DeleteBlob(k appengine.BlobKey) (err error) {
	return blobstore.Delete(s.Context, k)

}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func init() {
	http.HandleFunc("/p/dl", gaesend)
	http.HandleFunc("/p/ul", gaereceive)
	http.HandleFunc("/p/nl", gaenewer)
}
