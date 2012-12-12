package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var artifactFlag *string = flag.String("a", "", "set this artifact's name.")
var groupFlag *string = flag.String("g", "", "set this group's name.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GOR_VERSION)
		return
	}

	fmt.Println("Reading local information")
	p,_ := got.ReadProject() // read or create
	
	if *artifactFlag != "" {
		p.Artifact = *artifactFlag
	}
	if *groupFlag != "" {
		p.Group = *groupFlag
	}
	
	got.WriteProject(p)
	
	
	
}
