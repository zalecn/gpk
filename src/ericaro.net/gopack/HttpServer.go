package gopack

import (
	"ericaro.net/gopack/protocol"
	"io"
	"log"
	"net/http"
)
//HttpServer serve a local repository as a remote
type HttpServer struct {
	Local LocalRepository // handles the real operations
	server http.Server
}


//Start starts an http server at the addr provided.
func (s *HttpServer) Start(addr string) {
	mux := http.NewServeMux()
	protocol.HandleMux("/", s, mux) 
	s.server = http.Server{
		Addr:    addr,
		Handler: mux,
	}
	s.server.ListenAndServe()
}

func (s *HttpServer) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

//Contains return true if the server contains the ProjectID
func (s *HttpServer) Receive(pid protocol.PID, r io.ReadCloser) (err error)            {
	_, err = s.Local.Install(r)
	return
}

func (s *HttpServer) Serve(pid protocol.PID, w io.Writer) (err error)                    {
//func (s *StandaloneBackendServer) Send(id gopack.ProjectID, w http.ResponseWriter, r *http.Request) {
	id := *NewProjectID(pid.Name,pid.Version)
	p, err := s.Local.FindPackage(id)
	if err != nil {
		return
	}
	p.Pack(w)
	return
}

func (s *HttpServer) Search(query string, start int) (pids []protocol.PID, err error) {
	pids = s.Local.Search(query, start)
	return
}
