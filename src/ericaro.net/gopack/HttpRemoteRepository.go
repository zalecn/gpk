package gopack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
)


func init() {
	http := func(name string, u url.URL) RemoteRepository{
		h,_ := NewHttpRemoteRepository(name, u)
		return h
		}
	RegisterRemoteRepositoryFactory("http", http)
	RegisterRemoteRepositoryFactory("https", http)
	
}

// contains a remote repo based on http
type HttpRemoteRepository struct {
	ServerHost url.URL
	name string
}

func NewHttpRemoteRepository(name string, url url.URL) (remote *HttpRemoteRepository, err error){
	return &HttpRemoteRepository{
		ServerHost: url,
		name: name,
		}, nil

} 

func (r HttpRemoteRepository) Name() string {return r.name}
func (r HttpRemoteRepository) Path() url.URL {return r.ServerHost	}


//ReadPackage from this remote repository. Reads the http request's body into a buffer and returns.
func (r *HttpRemoteRepository) ReadPackage(p ProjectID) (reader io.Reader, err error) {

	// prepare central server query args
	v := url.Values{}
	v.Set("n", p.name)
	v.Set("v", p.version.String())

	s := r.ServerHost
	//query url
	u := &url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     "/p/dl", //make it configurable
		RawQuery: v.Encode(),
	}
	u = s.ResolveReference(u)
	
	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body) // download the tar.gz
	resp.Body.Close()
	if err != nil {
		return
	}
	reader = buf
	return
}



func (r *HttpRemoteRepository) CheckPackageCanPush(p *Package) (canpush bool, err error) {
	v := url.Values{}
	v.Set("n", p.self.name)
	v.Set("v", p.version.String())
	v.Set("t", p.timestamp.Format(time.ANSIC)) // ?

	//query url
	u := r.ServerHost.ResolveReference(&url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     "p/cp",
		RawQuery: v.Encode(),
	})
	
	resp, err := http.Get(u.String())
	if err != nil {
		return true, err
	}
	if resp.StatusCode != 200 || resp.StatusCode != 404 {
		err = errors.New(fmt.Sprintf("http query failed %d: %v", resp.StatusCode, resp.Status))
	}
	return resp.StatusCode == 200, err
}

func (r *HttpRemoteRepository) CheckPackageUpdate(p *Package) (newer bool, err error) {
	v := url.Values{}
	v.Set("n", p.self.name)
	v.Set("v", p.version.String())
	v.Set("t", p.timestamp.Format(time.ANSIC)) // ?

	//query url
	u := r.ServerHost.ResolveReference(&url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     "p/nl",
		RawQuery: v.Encode(),
	})
	
	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	if resp.StatusCode != 200 || resp.StatusCode != 404 {
		err = errors.New(fmt.Sprintf("http query failed %d: %v", resp.StatusCode, resp.Status))
	}
	return resp.StatusCode == 200, err

}

//UploadProject upload a project to the central server.
// the optional parameter snapshot, and version must be set
func (r *HttpRemoteRepository) UploadPackage(p *Package) (err error) { // TODO add a token authentication here
	// package it in memory
	
	buf := new(bytes.Buffer)
	p.Pack(buf)
	fmt.Printf("uploading %d\n", buf.Len())
	// prepare central server query args

	// these are the metadata sent to the remote server, so that it does not need to "read" the blob
	v := url.Values{}
	v.Set("n", p.self.name)
	v.Set("v", p.version.String())
	v.Set("t", p.timestamp.Format(time.ANSIC)) //?

	//query url
	u := r.ServerHost.ResolveReference(&url.URL{
		Path:     "p/ul",
		RawQuery: v.Encode(),
	})
	var client http.Client
	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		fmt.Printf("invalid request %v\n", err)
		return
	}
	req.ContentLength = int64(buf.Len())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New(fmt.Sprintf("http upload failed %d: %v", resp.StatusCode, resp.Status))
	}
	fmt.Printf("uploaded")
	return
}

func (r *HttpRemoteRepository) SearchPackage(search string) (result []string){
	v := url.Values{}
	v.Set("n", search)
	//query url
	u := r.ServerHost.ResolveReference(&url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     "p/qp",
		RawQuery: v.Encode(),
	})
	
	resp, err := http.Get(u.String())
	if err != nil {
		return result
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	//fmt.Printf("Received %s\n", result)
	return result
}

