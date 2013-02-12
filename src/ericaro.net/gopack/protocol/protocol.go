//protocol defines the basic remote protocols between a client and a remote repository
// it defines a Client interface that all remotes must implement
package protocol

import (
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"
	"log"
)

const ( // codes operations
	FETCH  = "fetch"
	PUSH   = "push"
	PUSH_EXEC  = "pushx"
	SEARCH = "search"
)

//ProtocolError is an error, but adds an error code. This module provides several "standard" errors
type ProtocolError struct {
	Message string
	Code    int
}

//Error part of the error interface.
func (p *ProtocolError) Error() string { return p.Message }

//Standard errors
var ( 
	StatusForbidden         = &ProtocolError{"Forbidden Operation", http.StatusForbidden}
	StatusIdentityMismatch  = &ProtocolError{"Mismatch between Identity Declared and Received", http.StatusExpectationFailed}
	StatusCannotOverwrite   = &ProtocolError{"Cannot Overwrite a Package", http.StatusConflict}
	StatusMissingDependency = &ProtocolError{"Missing Dependency", http.StatusNotAcceptable}
)

// convert any error into a suitable error code. it uses http.StatusInternalServerError if this is not a protocol error
func ErrorCode(err error) int {
	switch e := err.(type) {
	case *ProtocolError:
		return e.Code
	}
	return http.StatusInternalServerError
}

//Server is an interface that a server should implement 
type Server interface {
	//Receive will process the package.
	// you can use the pid to perform some quick checks before reading the package in r
	//r is a reader to a tar.gzed stream containing the package and the .gpk
	Receive(pid PID, r io.ReadCloser) error
	
	ReceiveExecutables(pid PID, r io.ReadCloser) error
	
	//Serve is expected to find the package and write it down to the the writer interface.
	// w must be a tar.gzed stream containing all the package structure, and a .gpk file
	Serve(pid PID, w io.Writer) error
	//Search actually perform the query and return a list of PID found
	Search(query string, start int) ([]PID, error)
	// The handlers make use of a debugf function.	
	Debugf(format string, args ...interface{})
}

// not commenting below, I'm not very happy with that. I can reused this code properly, it is not the good way to do it I guess.

func Handle(p string, s Server) { HandleMux(p, s, http.DefaultServeMux) }

func HandleMux(p string, s Server, mux *http.ServeMux) {
	mux.HandleFunc(path.Join(p, PUSH), func(w http.ResponseWriter, r *http.Request) {
		servePush(s, w, r)
	})
	mux.HandleFunc(path.Join(p, PUSH_EXEC), func(w http.ResponseWriter, r *http.Request) {
		serveBuilt(s, w, r)
	})
	mux.HandleFunc(path.Join(p, FETCH), func(w http.ResponseWriter, r *http.Request) {
		serveFetch(s, w, r)
	})
	mux.HandleFunc(path.Join(p, SEARCH), func(w http.ResponseWriter, r *http.Request) {
		serveSearch(s, w, r)
	})
}

//Receive HandlerFunc that s
func servePush(s Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // on the push URL only POST method are supported
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		log.Printf("%s not a POST request. %s instead", PUSH, r.Method)
		return
	}

	// identify the package
	vals := r.URL.Query()
	pid, err := FromParameter(&vals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s invalid parameters. %s", PUSH, err)
		return
	}
	err = s.Receive(*pid, r.Body) // create and fill the blob
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		log.Printf("%s Receive Error. %s", PUSH, err)
	}
	// can pass the reason as body response)
	//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//	fmt.Fprintln(w, error)

}

//Receive HandlerFunc that s
func serveBuilt(s Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { // on the built URL only POST method are supported
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		log.Printf("%s not a POST request. %s instead", PUSH_EXEC, r.Method)
		return
	}

	// identify the package
	vals := r.URL.Query()
	pid, err := FromParameter(&vals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s invalid parameters. %s", PUSH, err)
		return
	}
	err = s.ReceiveExecutables(*pid, r.Body) // create and fill the blob
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		log.Printf("%s Receive Error. %s", PUSH, err)
	}
	// can pass the reason as body response)
	//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//	fmt.Fprintln(w, error)

}

func serveFetch(s Server, w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	pid, err := FromParameter(&vals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s invalid parameters. %s", FETCH, err)
	}
	if pid.Name == "" {
		http.NotFound(w, r)
		return
	}
	err = s.Serve(*pid, w)
	if err != nil {
		http.Error(w, err.Error(), ErrorCode(err))
		log.Printf("%s Serve Error. %s", PUSH, err)
	}
	return
}

func serveSearch(s Server, w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	start, _ := strconv.Atoi(r.FormValue("start"))
	results, err := s.Search(query, start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s Search Error. %s", SEARCH, err)

	} else {
		json.NewEncoder(w).Encode(results)
	}
}
