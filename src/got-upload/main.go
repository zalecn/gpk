package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
	"log"
	"net/http"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"net/url"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
//var realFlag *bool = flag.Bool("real", false, "Use the real remote server.")
var releaseFlag *bool = flag.Bool("r", false, "Upload in release mode.")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
		return
	}

	log.Printf("running got upload %s\n", flag.Arg(0))
	version := got.ParseVersionReference(flag.Arg(0))
	p,_ := got.ReadProject() // read //TODO assert I'm on a real project
	
	
	buf := new(bytes.Buffer) 	
	gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	tw := tar.NewWriter(gz)
	if err != nil {
		panic(err)
	}
	p.PackageProject(tw) 
	
 	gz.Close()
 	tw.Close()
	
	
	
	
	var client http.Client
	v:= url.Values{}
	v.Set("g", p.Group)
	v.Set("a", p.Artifact)
	v.Set("v", version.String() )
	if *releaseFlag {
		v.Set("r", "true")
	}
	u := url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Scheme: "http",
		Host: got.GorCentral,
		Path: "/p/ul",
		RawQuery: v.Encode(),
	}
//	if ! *realFlag {u.Host = "localhost:8080"}
	
	
	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		fmt.Printf("invalid request %v\n", err)
		return
	}
	req.ContentLength = int64(buf.Len())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("upload error %v\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("upload failed %d: %v\n", resp.StatusCode,  resp.Status)
	}
	
}
