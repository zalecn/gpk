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

var serverAddrFlag *string 

var Serve = Command{
	Name:      `serve`,
	Alias:     `serve`,
	UsageLine: `ADDR`,
	Short:     `Serve local repository as an http server`,
	Long: `Serve local repository as an http remote repository so that others can get latest updates, or push new releases.
  ADDR usually ':8080' `,
	RequireProject: false, // false if we add the options to set which the local repo
	FlagInit: func(Serve *Command) {
	serverAddrFlag  = Serve.Flag.String("s", ":8080", "Serve the current local repository as a remote one for others to use.")
	},
	Run: func(Serve *Command) {

		// run the go build command for local src, and with the appropriate gopath

		server := HttpServer{
			Local: *Serve.Repository,
		}
		fmt.Printf("starting server %s\n", *serverAddrFlag)
		server.Start(*serverAddrFlag)

	},
}

//var deployAddrFlag *string = Push.Flag.String("to", "central", "deploy to a specific remote repository.")
//var pushRecursiveFlag *bool = Push.Flag.Bool("r", false, "Also pushes package's dependencies.")
var Push = Command{
	Name:      `push`,
	Alias:     `push`,
	UsageLine: `REMOTE PACKAGE VERSION`,
	Short:     `Push a project in a remote repository`,
	Long: `Push a project in a remote repository
  REMOTE  a remote name in the remote list
  PACKAGE a package available in the local repository (use search to list them)
  VERSION a semantic version of the PACKAGE available in the local repository`,
	RequireProject: false,
	Run: func(Push *Command) {
		rem := Push.Flag.Arg(0)
		remote, err := Push.Repository.Remote(rem)
		if err != nil {
			ErrorStyle.Printf("Unknown Remote %s.\n    \u21b3 %s\n", rem, err)

			fmt.Printf("Available remotes are:\n")
			for _, r := range Push.Repository.Remotes() {
				u := r.Path()
				fmt.Printf("    %-40s %s\n", r.Name(), u.String())
			}
			return
		}

		version, err := semver.ParseVersion(Push.Flag.Arg(2))
		if err != nil {
			ErrorStyle.Printf("Invalid Version \"%s\".\n    \u21b3 %s\n", Push.Flag.Arg(2), err)
			return
		}

		// now look for the real package in the local repo
		pkg, err := Push.Repository.FindPackage(*NewProjectID(Push.Flag.Arg(1), version))
		if err != nil {
			ErrorStyle.Printf("Cannot find Package %s %s in Local Repository %s.\n    \u21b3 %s\n", Push.Flag.Arg(1), Push.Flag.Arg(2), Push.Repository.Root(), err)
			// TODO as soon as I've got some search capability display similar results
			return
		}

		// build its ID 
		tm := pkg.Timestamp()
		pid := protocol.PID{
			Name:      Push.Flag.Arg(1),
			Version:   version,
			Token:     remote.Token(),
			Timestamp: &tm,
		}

		// read it in memory (tar.gz)
		buf := new(bytes.Buffer)
		pkg.Pack(buf)

		// and finally push the buffer
		err = remote.Push(pid, buf)
		if err != nil {
			ErrorStyle.Printf("Error from the remote while pushing.\n    \u21b3 %s\n", err)
			// TODO as soon as I've got some search capability display similar results
			return
		}
		SuccessStyle.Printf("Success\n")
	},
}

var Get = Command{
	Name:           `goget`,
	Alias:          `gg`,
	UsageLine:      `<goget package>`,
	Short:          `Run go get a package and install it`,
	Long:           `Run go get a package and install it`,
	Run:            func(Get *Command) { panic("not yet implemented") },
	RequireProject: false,
}
