package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	log.Printf("running gor download \n")
	
	
	r, err := http.Get("http://localhost:8080/dl?group=test&artifact=moi&version=master-1.0.0.0")
	if err != nil {
		fmt.Printf("get error %v\n", err)
		return
	}
	f, err := os.Create("toto.tar.gz")
	fmt.Printf("create error %v\n", err)
	fmt.Printf("to be read %d %d\n", r.StatusCode, r.ContentLength)
	j,err:=io.Copy(f, r.Body)
	
	fmt.Printf("read %d %v\n", j, err)
	f.Close()
	r.Body.Close()
	if err != nil {
		fmt.Printf("download error %v\n", err)
		return
	}
	fmt.Printf("done\n" )
	
	
	
	
	
	
}
