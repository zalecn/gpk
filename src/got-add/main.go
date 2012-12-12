package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"got.ericaro.net/got/cmd"
	"os"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

const (
	Usage = `Usage of got-add:
got add <group>:<artifact>:<version name>-<version 4 digits>

add the dependencies to this project dependency list.
`
)

func main() {
	flag.Parse() // Scan the arguments list
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, Usage)
		flag.PrintDefaults()
	}
	if *versionFlag {
		cmd.PrintVersion()
		return
	}

	p, _ := got.ReadProject() // read or create
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}
	for _, v := range flag.Args() {
		fmt.Printf("  -> %v\n", v)
		p.AppendDependency(got.ParseProjectReference(v))
	}

	got.WriteProject(p)
	fmt.Println(p)

}
