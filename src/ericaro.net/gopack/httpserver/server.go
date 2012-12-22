package httpserver

import (
	"ericaro.net/gpk"
	"net/http"
	"time"
)

// this server works for as standalone and as gae.
// it delegates a part of the job to the backendServer.

type BackendServer interface {
	Debugf(format string, args ...interface{})
	//Contains return true if the server contains the ProjectID
	Send(id gpk.ProjectID, w http.ResponseWriter, r *http.Request)
	
	CanPush(id gpk.ProjectID,timestamp time.Time,  w http.ResponseWriter, r *http.Request) 
	Newer(id gpk.ProjectID,timestamp time.Time,  w http.ResponseWriter, r *http.Request) 
	Receive(id gpk.ProjectID, timestamp time.Time, w http.ResponseWriter, r *http.Request) (err error)
	SearchPackage(search string, w http.ResponseWriter, r *http.Request)
	
}

//Receive HandlerFunc that s
func Receive(s BackendServer, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}

	// identify the package
	vals := r.URL.Query()
	name := vals.Get("n") // todo validate the syntax
	version, _ := gpk.ParseVersion(vals.Get("v"))
	timestamp, _ := time.Parse(time.ANSIC, vals.Get("t"))

	// try to get the package if it already exists
	pr := gpk.NewProjectID(name, version)
	// it's ok to create it

	err := s.Receive(pr, timestamp, w, r) // create and fill the blob
	if err != nil {
		s.Debugf("Error %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Debugf("Received %v\n", pr)

}

func Newer(s BackendServer, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("n")
	version, _ := gpk.ParseVersion(r.FormValue("v"))
	timestamp, _ := time.Parse(time.ANSIC, r.FormValue("t"))

	if name == "" {
		http.NotFound(w, r)
		return
	}
	pr := gpk.NewProjectID(name, version)
	s.Newer(pr, timestamp,  w, r)

}
func CanPush(s BackendServer, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("n")
	version, _ := gpk.ParseVersion(r.FormValue("v"))
	timestamp, _ := time.Parse(time.ANSIC, r.FormValue("t"))

	if name == "" {
		http.NotFound(w, r)
		return
	}
	pr := gpk.NewProjectID(name, version)
	s.CanPush(pr, timestamp,  w, r)

}

func SearchPackage(s BackendServer, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("n")
	if name == "" {
		http.NotFound(w, r)
		return
	}
	s.SearchPackage(name,  w, r)
}

func Send(s BackendServer, w http.ResponseWriter, r *http.Request) {

	s.Debugf("Receiving %s  %s \n", r.FormValue("n"), r.FormValue("v"))
	name := r.FormValue("n")
	version, _ := gpk.ParseVersion(r.FormValue("v"))
	if name == "" {
		http.NotFound(w, r)
		return
	}
	pr := gpk.NewProjectID(name, version)
	
	s.Send(pr, w, r)
	return
}
