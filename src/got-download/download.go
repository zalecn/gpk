package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
		return
	}

	log.Printf("running got download \n")
	p, _ := got.ParseProjectReference(flag.Arg(0))

	v := url.Values{}
	v.Set("g", p.Group)
	v.Set("a", p.Artifact)
	v.Set("v", p.Version.String())

	u := url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Scheme:   "http",
		Host:     got.GotCentral,
		Path:     "/p/dl",
		RawQuery: v.Encode(),
	}
	r, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("get error %v\n", err)
		return
	}
	file := fmt.Sprintf("target/%v-%v-%v.tar.gz", p.Group, p.Artifact, p.Version)
	f, err := os.Create(file)
	_, err = io.Copy(f, r.Body)
	f.Close()
	r.Body.Close()
	if err != nil {
		fmt.Printf("download error %v\n", err)
		return
	}
}
