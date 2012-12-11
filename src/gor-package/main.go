package main

import (
	"flag"
	"go.ericaro.net/gor"
	"fmt"
	"log"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", gor.GOR_VERSION)
		return
	}

	log.Printf("running gor package \n")
	p,_ := gor.ReadProject() // read //TODO assert I'm on a real project
	
	log.Printf("packaging project %v:%v\n", p.Group, p.Artifact)
	
	err := gor.MakeTarget()
	if err != nil {
		panic(err)
	}
	fw, gz, tw, err := gor.CreateTarGz(fmt.Sprintf("target/%v-%v.tar.gz", p.Group, p.Artifact))
	if err != nil {
		panic(err)
	}
 	defer fw.Close()
 	defer gz.Close()
 	defer tw.Close()
	p.PackageProject(tw) 
	
	
	
	
	
}
