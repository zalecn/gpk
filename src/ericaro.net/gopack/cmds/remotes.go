package cmds

import (
	"bytes"
	. "ericaro.net/gopack"
	"ericaro.net/gopack/protocol"
	"ericaro.net/gopack/semver"
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
	UsageLine:      `<remote> <package> <version>`,
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

	server := HttpServer{
		Local: *c.Repository,
	}
	fmt.Printf("starting server %s\n", *serverAddrFlag)
	server.Start(*serverAddrFlag)

}

//var deployAddrFlag *string = Push.Flag.String("to", "central", "deploy to a specific remote repository.")
//var pushRecursiveFlag *bool = Push.Flag.Bool("r", false, "Also pushes package's dependencies.")

func (c *Command) Push() {
	//fmt.Printf("helloooooooooooooooooooo\n\n")
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

	version, err := semver.ParseVersion(c.Flag.Arg(2))
	if err != nil {
		ErrorStyle.Printf("Invalid Version: %s\n", err)
		return
	}

	pkg, err := c.Repository.FindPackage(NewProjectID(c.Flag.Arg(1), version))
	if err != nil {
		ErrorStyle.Printf("Cannot find Package %s %s in Local Repository %s. Due to %s\n", c.Flag.Arg(1), c.Flag.Arg(2), c.Repository.Root(), err)
		// TODO as soon as I've got some search capability display similar results
		return
	}
	tm := pkg.Timestamp()
	pid := protocol.PID{
		Name:    c.Flag.Arg(1),
		Version: version,
		Token: remote.Token(),
		Timestamp: &tm,
	}
	// parse locally and fill a buffer
	buf := new(bytes.Buffer)
	pkg.Pack(buf)
	// push the buffer
	err = remote.Push(pid, buf)
	if err != nil {
		ErrorStyle.Printf("Remote Error %s\n", err)
		// TODO as soon as I've got some search capability display similar results
		return
	}
	SuccessStyle.Printf("Success\n")
}

func (c *Command) Get() {
	panic("not yet implemented")
	//	for _, p := range c.Flag.Args() {
	//		//TODO fill it
	//		//c.Repository.GoGetInstall(p)
	//	}
}
