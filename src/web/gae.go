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
	"encoding/json"
)

//a Package is a pure, in-memory representation of a Package
type GaePackage struct {
	Timestamp   time.Time
	Package string // package name (for queries)
	Major, Minor, Patch int // version digits
	PreRelease string
	Build string
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
func gaecanpush(w http.ResponseWriter, r *http.Request) {
	be := &GaeBackendServer{
		Context: appengine.NewContext(r),
	}
	httpserver.CanPush(be, w, r)
}

func gaesearch(w http.ResponseWriter, r *http.Request) {
	be := &GaeBackendServer{
		Context: appengine.NewContext(r),
	}
	httpserver.SearchPackage(be, w, r)
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
	v:= pack.Version()
	maj,min, patch := v.Digits()
	p := &GaePackage{
		Package : pack.Name(),
		Timestamp : pack.Timestamp(),
		Major: int(maj),
		Minor: int(min),
		Patch: int(patch),
		PreRelease: v.PreRelease(),
		Build : v.Build(),
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

func (s *GaeBackendServer) CanPush(id gpk.ProjectID, timestamp time.Time, w http.ResponseWriter, r *http.Request) {
	c := s.Context
	p := new(GaePackage)
	err := datastore.Get(c, datastore.NewKey(c, "GaePackage", id.ID(), 0, nil), p)
	var canpush bool
	if err != nil {
		canpush = true
	} else {
		canpush = p.Timestamp.Before(timestamp)
	}
	
	if ! canpush{
		http.NotFound(w, r)
		return
	}
}

func (s *GaeBackendServer) SearchPackage(search string, w http.ResponseWriter, r *http.Request) {
		
	q := datastore.NewQuery("GaePackage").
        Filter("Package >=", search).
        Filter("Package <", search+string('\uFFFD')).
        Order("Package").
        Limit(100)
        
    //count, err := q.Count(s.Context) 
    //s.Context.Debugf("searching %s, %d, %s", search, count, err)
    result:= make([]string, 100)
    i:=0
    for t := q.Run(s.Context); ;i++ {
        x := new( GaePackage )
        _, err := t.Next(x)
        if err == datastore.Done {
                break
        }
        if err != nil {
                http.NotFound(w, r)
                return
        }
        //s.Context.Debugf("found %s", x.Package)
        result[i]= x.Package
    }
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result[:i])
	
}

func (s *GaeBackendServer) DeleteBlob(k appengine.BlobKey) (err error) {
	return blobstore.Delete(s.Context, k)

}



type handlerFunc func(http.ResponseWriter, *http.Request) error

func init() {
	http.HandleFunc("/p/dl", gaesend)
	http.HandleFunc("/p/ul", gaereceive)
	http.HandleFunc("/p/nl", gaenewer)
	http.HandleFunc("/p/cp", gaecanpush)
	http.HandleFunc("/p/qp", gaesearch)
}
