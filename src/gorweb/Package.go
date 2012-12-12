package gorweb

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"go.ericaro.net/gor"
	"io"
	"net/http"
)

//a Package is a pure, in-memory representation of a Package
type Package struct {
	Group, Artifact, Root, Version      string
	ContentBlob                appengine.BlobKey // tar.gz of the content
}



func receive(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}
	c := appengine.NewContext(r)
	c.Debugf("processing receive\n")

	// identify the package
	vals := r.URL.Query()
	group := vals.Get("group")
	artifact:=    vals.Get("artifact")
	version := gor.ParseVersionReference(vals.Get("version"))
	c.Debugf("creating package %v:%v:%v\n", group, artifact, version)
	
	blob, err := CreateBlob(c, r.Body) // create and fill the blob
	p := Package{
		Group:       group,
		Artifact:    artifact,
		Root:        version.Root,
		Version:       version.Parts,
		ContentBlob: blob,
	}
	pr := gor.NewProjectReference(group, artifact, version)
	
	_, err = datastore.Put(c, datastore.NewKey(c, "Package", pr.String(), 0, nil), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Debugf("received\n")
}

func CreateBlob(c appengine.Context, r io.Reader) (k appengine.BlobKey, err error) {
	w, err := blobstore.Create(c, "application/x-gzip")
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

func serve(w http.ResponseWriter, r *http.Request) {
	group := r.FormValue("group")
	artifact := r.FormValue("artifact")
	version := gor.ParseVersionReference(r.FormValue("version"))

	if group == "" || artifact == "" {
		http.NotFound(w, r)
		return
	}
	pr := gor.NewProjectReference(group, artifact, version) 
	var p Package
	c := appengine.NewContext(r)
	c.Debugf("processing serve\n")
	err := datastore.Get(c, datastore.NewKey(c, "Package", pr.String() , 0, nil), &p)
	if err == datastore.ErrNoSuchEntity {
		http.NotFound(w, r)
		return
	}
	c.Debugf("serving %#v\n", p)
	blobstore.Send(w, p.ContentBlob)
//	
//	stat, err := blobstore.Stat(c, p.ContentBlob)
//	w.Header().Set("Content-Type", "application/x-gzip")
//	w.Header().Set("Content-Length", fmt.stat.Size )
//	blob := blobstore.NewReader(c, p.ContentBlob )
//	_, err = io.Copy(w, blob)
//	if err == datastore.ErrNoSuchEntity {
//		http.NotFound(w, r)
//		return
//	}
	return
}


type handlerFunc func(http.ResponseWriter, *http.Request) error

func init() {
	http.HandleFunc("/dl", serve)
	http.HandleFunc("/ul", receive) 
}
