package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
	"log"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var releaseFlag *bool = flag.Bool("r", false, "Install as a Release.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
		return
	}

	log.Printf("running got install \n")
	p,_ := got.ReadProject() // read //TODO assert I'm on a real project
	
	log.Printf("installing project %v:%v\n", p.Group, p.Artifact)
	version := got.ParseVersion(flag.Arg(0) )
	log.Printf("version  %v\n", version)
	
	r,_ := got.NewDefaultRepository() // build the defaut rep
	
	r.InstallProject(p, version ,!*releaseFlag)
	
	
	
}
