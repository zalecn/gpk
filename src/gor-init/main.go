package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var artifactFlag *string = flag.String("a", "", "set this artifact's name.")
var groupFlag *string = flag.String("g", "", "set this group's name.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	fmt.Println("Reading local information")
	p,_ := gor.ReadProject() // read or create
	
	if *artifactFlag != "" {
		p.Artifact = *artifactFlag
	}
	if *groupFlag != "" {
		p.Group = *groupFlag
	}
	
	gor.WriteProject(p)
	
	
	
}
