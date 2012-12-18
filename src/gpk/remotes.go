package main

import (
	. "ericaro.net/gpk"
	"ericaro.net/gpk/httpserver"
	"fmt"
	"net/url"
)

func init() {
	Reg(
		&Serve,
		&Deploy,
		&Get,
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


// TODO move around the remote tool chain
var Deploy = Command{
	Name:           `deploy`,
	Alias:          `d`,
	UsageLine:      `<version>`,
	Short:          `Deploy the current project in the remote repository`,
	Long:           `Deploy the current project in the remote repository`,
	call:           func(c *Command) { c.Deploy() },
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



var deployAddrFlag *string = Deploy.Flag.String("to", "http://gpk.ericaro.net", "deploy to a specific remote repository.")
func (c *Command) Deploy() {

	version, _ := ParseVersion(c.Flag.Arg(0))
	fmt.Printf("Installing %v : %v\n", c.Project, version)
	pkg := c.Repository.InstallProject(c.Project, version) // ensure that the project is installed first in the local repo
	
	// then tries to upload it
	u,_ := url.Parse(*deployAddrFlag)
	
	
	remote := NewRemoteRepository(*u)
	fmt.Printf("Uploading %v : %v to %v\n", c.Project, version, u)
	remote.UploadPackage(pkg)
}

func (c *Command) Get() {
		panic("not yet implemented")
//	for _, p := range c.Flag.Args() {
//		//TODO fill it
//		//c.Repository.GoGetInstall(p)
//	}
}
