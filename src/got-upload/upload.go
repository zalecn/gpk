package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
	"log"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var releaseFlag *bool = flag.Bool("r", false, "Upload in release mode.")
var hostFlag *string = flag.String("host", got.GotCentral,  "Set the host for the central server")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
		return
	}
	

	log.Printf("running got upload %s\n", flag.Arg(0))
	version := got.ParseVersion(flag.Arg(0))
	p,_ := got.ReadProject() // read //TODO assert I'm on a real project
	r, err := got.NewDefaultRepository()
	if err!= nil {
		fmt.Println(err)
		return
	}
	r.ServerHost = *hostFlag
	
	// prepare the project identity
	snapshot := !*releaseFlag
	p.Snapshot = &snapshot
	p.Version = &version
	
	
	err = r.UploadProject(p)
	if err != nil {
		fmt.Println(err)
		return
	}	
		
}
