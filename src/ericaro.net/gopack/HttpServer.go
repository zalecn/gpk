package gopack

import (
	"ericaro.net/gopack/protocol"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//HttpServer serve a local repository as a remote
type HttpServer struct {
	Local  LocalRepository // handles the real operations
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
func (s *HttpServer) Receive(pid protocol.PID, r io.ReadCloser) (err error) {
	pak, err := s.Local.Install(r)
	if err != nil {
		return
	}
	log.Printf("RECEIVING %s %s %s INTO %s", pak.Name(), pak.Version().String(), pak.License(), pak.InstallDir())
	return
}

//Contains return true if the server contains the ProjectID
func (s *HttpServer) ReceiveExecutables(pid protocol.PID, r io.ReadCloser) (err error) {
	pak, err := s.Local.InstallAppend(r)
	if err != nil {
		return
	}
	log.Printf("RECEIVING EXECUTABLES %s %s %s INTO %s", pak.Name(), pak.Version().String(), pak.License(), pak.InstallDir())
	return
}
func (s *HttpServer) Serve(pid protocol.PID, w io.Writer) (err error) {
	//func (s *StandaloneBackendServer) Send(id gopack.ProjectID, w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVING %s %s", pid.Name, pid.Version.String())
	id := *NewProjectID(pid.Name, pid.Version)
	p, err := s.Local.FindPackage(id)

	if err != nil {
		return
	}
	if pid.Timestamp != nil && p.Timestamp().Before(*pid.Timestamp) {
		return protocol.StatusNotNewPackage
	}

	if *pid.Executables {
		p.PackExecutables(w)
	} else {
		p.Pack(w)
	}
	return
}

func (s *HttpServer) Get(pid protocol.PID, goos, goarch, name string, w io.Writer) (err error) {
	//func (s *StandaloneBackendServer) Send(id gopack.ProjectID, w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVING Exec %s %s %s %s %s", pid.Name, pid.Version.String(), goos, goarch, name)
	id := *NewProjectID(pid.Name, pid.Version)
	p, err := s.Local.FindPackage(id)
	if err != nil {
		return
	}
	dst := filepath.Join(p.InstallDir(), "bin", goos+"_"+goarch, name)
	f, err := os.Open(dst)
	if err != nil {
		return
	}
	defer f.Close()
	n, err := io.Copy(w, f)
	log.Printf("Sent  %d b", n)
	return
}

func (s *HttpServer) List(pid protocol.PID, goos, goarch string, w io.Writer) (list []string, err error) {
	//func (s *StandaloneBackendServer) Send(id gopack.ProjectID, w http.ResponseWriter, r *http.Request) {
	log.Printf("LIST Exec %s %s %s %s", pid.Name, pid.Version.String(), goos, goarch)
	id := *NewProjectID(pid.Name, pid.Version)
	p, err := s.Local.FindPackage(id)
	if err != nil {
		return
	}

	dst := filepath.Join(p.InstallDir(), "bin", goos+"_"+goarch)
	f, err := os.Open(dst)
	if err != nil {
		return
	}
	defer f.Close()
	subdir, err := f.Readdir(-1)
	list = make([]string, 0)
	for _, fi := range subdir {
		if !fi.IsDir() {
			list = append(list, fmt.Sprintf(`/get?n=%s&v=%s&goos=%s&goarch=%s&exe=%s`, pid.Name, pid.Version.String(), goos, goarch, fi.Name()))
		}
	}
	return
}

func (s *HttpServer) Search(query string, start int) (pids []protocol.PID, err error) {
	pids = s.Local.Search(query, start)
	return
}
