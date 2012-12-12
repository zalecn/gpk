package main

import (
	"flag"
	"got.ericaro.net/got"
	"fmt"
	"log"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", got.GOR_VERSION)
		return
	}

	log.Printf("running got package \n")
	p,_ := got.ReadProject() // read //TODO assert I'm on a real project
	
	log.Printf("packaging project %v:%v\n", p.Group, p.Artifact)
	
	err := got.MakeTarget()
	if err != nil {
		panic(err)
	}
	fw, gz, tw, err := got.CreateTarGz(fmt.Sprintf("target/%v-%v.tar.gz", p.Group, p.Artifact))
	if err != nil {
		panic(err)
	}
 	defer fw.Close()
 	defer gz.Close()
 	defer tw.Close()
	p.PackageProject(tw) 
	
	
	
	
	
}
