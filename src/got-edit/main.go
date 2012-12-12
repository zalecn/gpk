package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"got.ericaro.net/got/cmd"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

var artifactFlag *string = flag.String("a", "", "set this artifact's name.")
var groupFlag *string = flag.String("g", "", "set this group's name.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		cmd.PrintVersion()
		return
	}

	p, err := got.ReadProject() // read or create
	done := false
	if err != nil {
		done = true // meaning that the project is new, the previous one has failed to be read
	}

	if *groupFlag != "" {
		fmt.Printf("group    <- %v\n", *groupFlag)
		p.Group = *groupFlag
		done = true
	}
	if *artifactFlag != "" {
		fmt.Printf("artifact <- %v\n", *artifactFlag)
		p.Artifact = *artifactFlag
		done = true
	}
	if done {
		got.WriteProject(p) // store only if needed
		fmt.Println(p)
	} else {
		flag.Usage()
	}
}
