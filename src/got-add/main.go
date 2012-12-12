package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
	"os"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")


func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GOR_VERSION)
		return
	}

	os.Setenv("toto", "titi")
	p,_ := got.ReadProject() // read or create
	
	for _, v :=range flag.Args() {
		fmt.Printf("add %v\n", v)
		p.AppendDependency(  got.ParseProjectReference(v) )
	}
	got.WriteProject(p)
	fmt.Printf("Status %v\n", p )
	
	
	
}
