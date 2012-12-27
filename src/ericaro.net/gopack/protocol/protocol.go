//protocol defines the basic remote protocols between a client and a remote repository
// it defines a Client interface, and an http implementation. There is also a directory based implementation available in the localrepository package
package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
)

const (
	FETCH  = "fetch"
	PUSH   = "push"
	SEARCH = "search"
)

type ProtocolError int

var (
	StatusForbidden = ProtocolError(http.StatusForbidden)
	StatusOK        = ProtocolError(http.StatusOK)
)

func (p *ProtocolError) Error() string {
	switch *p {
	case StatusForbidden:
		return "Forbidden Operation"
	}
	return fmt.Sprintf("Unknown error code: %d", p)
}

type Server interface {
	Receive(pid PID, r io.ReadCloser) (ProtocolError, error)
	Serve(pid PID, w io.Writer) (ProtocolError, error)
	Search(query string, start int) ([]PID, error)
	Debugf(format string, args ...interface{})
}

func Handle(p string, s Server) { HandleMux(p, s, http.DefaultServeMux) }

func HandleMux(p string, s Server, mux *http.ServeMux) {
	mux.HandleFunc(path.Join(p, PUSH), func(w http.ResponseWriter, r *http.Request) {
		servePush(s, w, r)
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
		return
	}

	// identify the package
	vals := r.URL.Query()
	pid, err := FromParameter(&vals)
	if err != nil {
		s.Debugf("Error %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pe, err := s.Receive(*pid, r.Body) // create and fill the blob
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(int(pe))
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
	}
	if pid.Name == "" {
		http.NotFound(w, r)
		return
	}

	pe, err := s.Serve(*pid, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(int(pe))
	}
	return
}

func serveSearch(s Server, w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	start, _ := strconv.Atoi(r.FormValue("start"))
	results, err := s.Search(query, start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {
		json.NewEncoder(w).Encode(results)
	}
}
