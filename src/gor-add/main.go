package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	fmt.Println("Reading local information")
	p,_ := gor.ReadProject() // read or create
	
	for _, v :=range flag.Args() {
		fmt.Printf("arg%v\n", v)
		p.AppendDependency(  gor.ParseProjectReference(v) )
	}
	
	gor.WriteProject(p)
	
	
	
}
