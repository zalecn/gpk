package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"got.ericaro.net/got/cmd"
)

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var releaseFlag *bool = flag.Bool("r", false, "Build using only release dependencies.")
var offlineFlag *bool = flag.Bool("o", true,  "Try to find missing dependencies at http://got.ericaro.net")
var updateFlag *bool = flag.Bool("u", false,  "Look for updated version of dependencies")


func main() {
	flag.Parse() // Scan the arguments list 
	if *versionFlag {
		cmd.PrintVersion()
		return
	}

	p, _ := got.ReadProject() // creates a new empty project with default values
	r, _ := got.NewDefaultRepository()


	// parse dependencies, and build the gopath
	dependencies, err := r.FindProjectDependencies(p, !*releaseFlag, *offlineFlag, *updateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies", err)
		return
	}
	
	// run the go build command for local src, and with the appropriate gopath
	gopath, err := r.GoPath(dependencies)
	//fmt.Printf("gopath=%v\n", gopath)
	goEnv := got.NewGoEnv(gopath)
	goEnv.Install(p.Root)
	
}
