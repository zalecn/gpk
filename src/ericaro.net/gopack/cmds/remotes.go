package cmds

import (
	"bytes"
	. "ericaro.net/gopack"
	"ericaro.net/gopack/protocol"
	"ericaro.net/gopack/semver"
	"fmt"
	"net/url"
)

func init() {
	Reg(
		&Serve,
		&Push,
		&AddRemote,
		&RemoveRemote,
	)

}

var serverAddrFlag *string

var Serve = Command{
	Name:      `serve`,
	Alias:     `serve`,
	UsageLine: `ADDR`,
	Short:     `Serve local repository as an http server`,
	Long: `Serve local repository as an http remote repository
       so that others can get latest updates, or push new releases.
       ADDR usually ':8080' `,
	RequireProject: false, // false if we add the options to set which the local repo
	FlagInit: func(Serve *Command) {
		serverAddrFlag = Serve.Flag.String("s", ":8080", "Serve the current local repository as a remote one for others to use.")
	},
	Run: func(Serve *Command) (err error) {

		// run the go build command for local src, and with the appropriate gopath

		server := HttpServer{
			Local: *Serve.Repository,
		}
		fmt.Printf("starting server %s\n", *serverAddrFlag)
		server.Start(*serverAddrFlag)
		return

	},
}

//var deployAddrFlag *string = Push.Flag.String("to", "central", "deploy to a specific remote repository.")
//var pushRecursiveFlag *bool = Push.Flag.Bool("r", false, "Also pushes package's dependencies.")
var pushExecutables *bool
var Push = Command{
	Name:      `push`,
	Alias:     `push`,
	UsageLine: `REMOTE PACKAGE VERSION`,
	Short:     `Push a project in a remote repository`,
	Long: `Push a project in a remote repository
       REMOTE  a remote name in the remote list
       PACKAGE a package available in the local repository (use search to list them)
       VERSION a semantic version of the PACKAGE available in the local repository
       `,
	RequireProject: false,
	FlagInit: func(Push *Command) {
		pushExecutables = Push.Flag.Bool("x", false, "pushes executables instead.")
	},
	Run: func(Push *Command) (err error) {
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

		if *pushExecutables {
			pkg.PackExecutables(buf) // pack either exec or src
			// and finally push the buffer
			err = remote.PushExecutables(pid, buf) // either exec or src		
		} else {
			pkg.Pack(buf) // pack either exec or src
			// and finally push the buffer
			err = remote.Push(pid, buf) // either exec or src
		}
		if err != nil {
			ErrorStyle.Printf("Error from the remote while pushing.\n    \u21b3 %s\n", err)
			// TODO as soon as I've got some search capability display similar results
			return
		}
		SuccessStyle.Printf("Success\n")
		return
	},
}


//var Get = Command{
//	Name:           `goget`,
//	Alias:          `gg`,
//	UsageLine:      `<goget package>`,
//	Short:          `Run go get a package and install it`,
//	Long:           `Run go get a package and install it`,
//	Run:            func(Get *Command) { panic("not yet implemented") },
//	RequireProject: false,
//}

////////////////////////////////////////////////////////////////////////////////////////

var AddRemote = Command{
	Name:      `radd`,
	Alias:     `r+`,
	Category:  RemoteCategory,
	UsageLine: `NAME URL [TOKEN]`,
	Short:     `Add a remote server.`,
	Long: `Remote server can be used to publish or receive other's code.
       NAME    local alias for this remote
       URL     full URL to the remote server.
               file:// and http:// are actually supported. For http, see 'gpk serve'
       TOKEN   option: a secret TOKEN identification provided by the server to deliver authentication.
               See with your server provider`,
	RequireProject: false,
	Run: func(AddRemote *Command) (err error) {

		if len(AddRemote.Flag.Args()) < 2 || len(AddRemote.Flag.Args()) > 3 {
			ErrorStyle.Printf("Illegal arguments count\n")
			return
		}

		name, remote := AddRemote.Flag.Arg(0), AddRemote.Flag.Arg(1)
		u, err := url.Parse(remote)
		if err != nil {
			ErrorStyle.Printf("Invalid URL passed as a remote Repository.\n    \u21b3 %s\n", err)
			return
		}
		// TOKEN handling
		var token *protocol.Token // nil by default
		if len(AddRemote.Flag.Args()) == 3 {
			token, err = protocol.ParseStdToken(AddRemote.Flag.Arg(2))
			if err != nil {
				ErrorStyle.Printf("Invalid token syntax, please enter a valid token base64- RFC 4648 Encoded array of bytes.\n")
				return
			}
		}

		client, err := protocol.NewClient(name, *u, token)
		if err != nil {
			ErrorStyle.Printf("Failed to create the a client for this remote:\n    \u21b3 %s\n", err)
		}
		stoken := ""
		if token != nil {
			stoken = fmt.Sprintf("%s", token)
		}
		err = AddRemote.Repository.RemoteAdd(client)
		if err != nil {
			ErrorStyle.Printf("%s\n", err)
			return
		}
		SuccessStyle.Printf("       +%s %s %s\n", name, u, stoken)
		AddRemote.Repository.Write()
		return
	},
}

var RemoveRemote = Command{
	Name:           `rremove`,
	Alias:          `r-`,
	Category:       RemoteCategory,
	UsageLine:      `NAME`,
	Short:          `Remove a Remote`,
	Long:           ``,
	RequireProject: false,
	Run: func(RemoveRemote *Command) (err error) {

		if len(RemoveRemote.Flag.Args()) != 1 {
			RemoveRemote.Flag.Usage()
			return
		}

		name := RemoveRemote.Flag.Arg(0)
		ref, err := RemoveRemote.Repository.RemoteRemove(name)
		if err != nil {
			ErrorStyle.Printf("Cannot Remove remote %s\n    \u21b3 %s", name, err)
			return
		}
		if ref != nil {
			u := ref.Path()
			SuccessStyle.Printf("Removed Remote %s %s\n", ref.Name(), u.String())
			RemoveRemote.Repository.Write()
		} else {
			ErrorStyle.Printf("Nothing to Remove\n")
		}
		RemoveRemote.Repository.Write()
		return
	},
}
