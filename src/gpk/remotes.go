package main

import (
	. "ericaro.net/gpk"
	"ericaro.net/gpk/httpserver"
	"fmt"
)

func init() {
	Reg(
		&Serve,
		&Push,
		&Get,
	)

}

var Serve = Command{
	Name:           `serve`,
	Alias:          `serve`,
	UsageLine:      `<local addr>`,
	Short:          `Serve local repository as an http server`,
	Long:           `Serve local repository as a remote repository so that others can get latest updates, or push new releases.`,
	call:           func(c *Command) { c.Serve() },
	RequireProject: false, // false if we add the options to set which the local repo
}

// TODO move around the remote tool chain
var Push = Command{
	Name:           `push`,
	Alias:          `push`,
	UsageLine:      `<remote> <version>`,
	Short:          `Push a project in a remote repository`,
	Long:           `Push a project in a remote repository`,
	call:           func(c *Command) { c.Push() },
	RequireProject: true,
}

var Get = Command{
	Name:           `goget`,
	Alias:          `gg`,
	UsageLine:      `<goget package>`,
	Short:          `Run go get a package and install it`,
	Long:           `Run go get a package and install it`,
	call:           func(c *Command) { c.Get() },
	RequireProject: false,
}

var serverAddrFlag *string = Serve.Flag.String("s", ":8080", "Serve the current local repository as a remote one for others to use.")

func (c *Command) Serve() {

	// run the go build command for local src, and with the appropriate gopath

	server := httpserver.StandaloneBackendServer{
		Local: *c.Repository,
	}
	fmt.Printf("starting server %s\n", *serverAddrFlag)
	server.Start(*serverAddrFlag)

}

//var deployAddrFlag *string = Push.Flag.String("to", "central", "deploy to a specific remote repository.")
func (c *Command) Push() {

	rem := c.Flag.Arg(0)
	remote, err := c.Repository.Remote(rem)
	if err != nil {
		ErrorStyle.Printf("Unknown remote %s.\n", rem)

		fmt.Printf("Available remotes are:\n")
		for _, r := range c.Repository.Remotes() {
			u := r.Path()
			fmt.Printf("    %-40s %s\n", r.Name(), u.String())
		}
		return
	}

	version, err := ParseVersion(c.Flag.Arg(1))
	if err != nil {
		ErrorStyle.Printf("Invalid Version: %s\n", err)
		return
	}

	fmt.Printf("Installing ...\n")
	// ? really ? 
	pkg := c.Repository.InstallProject(c.Project, version) // ensure that the project is installed first in the local repo

	u := remote.Path()
	fmt.Printf("Pushing to %s\n", u.String())
	remote.UploadPackage(pkg)
}

func (c *Command) Get() {
	panic("not yet implemented")
	//	for _, p := range c.Flag.Args() {
	//		//TODO fill it
	//		//c.Repository.GoGetInstall(p)
	//	}
}
