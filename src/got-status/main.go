package main

import (
	"fmt"
	"flag"
	"got.ericaro.net/got"
	"got.ericaro.net/got/cmd"
)

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list

	p, _ := got.ReadProject()
	fmt.Println(p)

	if *versionFlag {
		cmd.PrintVersion()
		return
	}
}
