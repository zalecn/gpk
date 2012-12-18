package main

import (
	. "ericaro.net/gpk"
	"ericaro.net/gpk/httpserver"
	"fmt"
)

func init() {
	Reg(
		&Serve,
	)

}

var Serve = Command{
	Name:           `serve`,
	Alias:          `serve`,
	UsageLine:      `<local addr>`,
	Short:          `Serve current local repository as a remote one`,
	Long:           `Serve the current local repository as a remote repository so that others can get latest updates, or push new releases.`,
	call:           func(c *Command) { c.Serve() },
	RequireProject: false, // false if we add the options to set which the local repo
}

var serverAddrFlag *string = Compile.Flag.String("s", ":8080", "Serve the current local repository as a remote one for others to use.")

func (c *Command) Serve() {

	// run the go build command for local src, and with the appropriate gopath
	server := httpserver.StandaloneBackendServer{
		Local: *c.Repository,
	}
	fmt.Printf("starting server %s\n", *serverAddrFlag)
	server.Start(*serverAddrFlag)

}
