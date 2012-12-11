package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
	"log"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var releaseFlag *bool = flag.Bool("r", false, "Install as a Release.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	log.Printf("running gor install \n")
	p,_ := gor.ReadProject() // read //TODO assert I'm on a real project
	
	log.Printf("installing project %v:%v\n", p.Group, p.Artifact)
	version := gor.ParseVersion(flag.Arg(0) )
	log.Printf("version  %v\n", version)
	
	r,_ := gor.NewDefaultRepository() // build the defaut rep
	
	r.InstallProject(p, version ,!*releaseFlag)
	
	
	
}
