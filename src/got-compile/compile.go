package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"os"
	"strings"
)

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var releaseFlag *bool = flag.Bool("r", false, "Build using only release dependencies.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
	}

	p, _ := got.ReadProject() // creates a new empty project with default values
	// now the project has been parsed

	r, _ := got.NewDefaultRepository()

	// parse dependencies, and build the gopath
	dependencies := r.FindProjectDependencies(p, !*releaseFlag)
	// run the go build command for local src, and with the appropriate gopath

	sources := make([]string, 0, len(dependencies))
	sources = append(sources,p.Root)
	for _, pr := range dependencies {
		sources = append(sources,pr.Root)
	}
	gopath := strings.Join(sources, string(os.PathListSeparator))
	fmt.Printf("gopath=%v\n", gopath)
	
	goEnv := got.NewGoEnv(gopath)
	goEnv.Install(p.Root)
	
}
