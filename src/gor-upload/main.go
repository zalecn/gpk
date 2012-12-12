package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
	"log"
	"net/http"
	"os"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	log.Printf("running gor upload \n")
	p,_ := gor.ReadProject() // read //TODO assert I'm on a real project
	
	err := gor.MakeTarget()
	if err != nil {
		panic(err)
	}
	targz := fmt.Sprintf("target/%v-%v.tar.gz", p.Group, p.Artifact)
	fw, gz, tw, err := gor.CreateTarGz(targz)
	if err != nil {
		panic(err)
	}
 	defer fw.Close()
 	defer gz.Close()
 	defer tw.Close()
	p.PackageProject(tw) 
	
	var client http.Client
	f,err := os.Open(targz)
	defer f.Close()
	stat , err := f.Stat()
	req, err := http.NewRequest("POST", "http://localhost:8080/ul?group=test&artifact=moi&version=master-1.0.0.0", f)
	req.ContentLength = stat.Size()
	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("upload error %v", err)
	}
	
	
	
	
	
	
	
}
