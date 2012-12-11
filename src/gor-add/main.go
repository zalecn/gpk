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

	p,_ := gor.ReadProject() // read or create
	
	for _, v :=range flag.Args() {
		p.AppendDependency(  gor.ParseProjectReference(v) )
	}
	gor.WriteProject(p)
	fmt.Printf("Status %v\n", p )
	
	
	
}
