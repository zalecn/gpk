package cmds

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/protocol"
	. "ericaro.net/gopack/semver"

	"fmt"
	"net/url"
	"os"
)

func init() {
	Reg(
		&Init,
		&Status,
		&Add,
		&Remove,
		&List,
		&Search,
		&AddRemote,
		&RemoveRemote,
	)

}

var initNameFlag *string
var initLicenseFlag *string

var Init = Command{
	Name:      `init`,
	Alias:     `!`,
	UsageLine: `-n NAME -l LICENSE`,
	Short:     `Initialize or edit the current project.`,
	Long: `Init the current directory creates or updates the gopack project file, and name the current package NAME, setting the license to LICENSE.

  License allowed strings are either alias or fullname in one of the licenses below:
       alias   fullname
       ASF     Apache License 2.0
       EPL     Eclipse Public License 1.0
       GPL2    GNU GPL v2
       GPL3    GNU GPL v3
       LGPL    GNU Lesser GPL
       MIT     MIT License
       MPL     Mozilla Public License 1.1
       BSD     New BSD License
       OOS     Other Open Source
       OCS     Other Closed Source
`,
	RequireProject: false,
	FlagInit: func(Init *Command) {
		initNameFlag = Init.Flag.String("n", "", "sets the project name")
		initLicenseFlag = Init.Flag.String("l", "", "sets the project's license.")
	},
	Run: func(Init *Command) {
		// init does not require a project => I need to parse it myself and ignore failure
		p, err := ReadProject()
		if err == nil {
			fmt.Printf("warning: init an existing project. This is fine if you wanted to edit it\n")
		}
		pwd, err := os.Getwd()
		if err != nil {
			ErrorStyle.Printf("Cannot create the project, there is no current directory. Because %v\n", err)
			return
		}
		p.SetWorkingDir(pwd) // a project need to always to where it is. 

		// processing edits
		if *initNameFlag != "" {
			p.SetName(*initNameFlag)
			fmt.Printf("new name:%s\n", p.Name())
		}

		if *initLicenseFlag != "" {
			var lic *License

			// sorry for that, this is ugly plumbing, I'll come back and fix later
			if l, err := Licenses.Get(*initLicenseFlag); err != nil {
				if l, err = Licenses.GetAlias(*initLicenseFlag); err != nil {
					fmt.Printf("new license: unknown or unsupported license:%s\n", p.License())
				} else {
					lic = l
				}
			} else {
				lic = l
			}
			// end of sorryness 

			if lic != nil {
				p.SetLicense(*lic)
				fmt.Printf("new license:\"%s\"\n", p.License().FullName)
			}
		}

		Init.Project = p // in case we implement sequence of commands (in the future)
		p.Write()        // store it  one day I'll implement a lock on this file, right ?

	},
}

var Status = Command{
	Name:           `status`,
	Alias:          `?`,
	UsageLine:      ``,
	Short:          `Print status`,
	Long:           `Display current information about the current project and the current local repository`,
	RequireProject: true,
	Run: func(Status *Command) {

		TitleStyle.Printf("    Name        :\n", Status.Project.Name())
		fmt.Printf("    License     : %s\n", Status.Project.License().FullName)
		dep := Status.Project.Dependencies()
		if len(dep) == 0 {
			fmt.Printf("    Dependencies: <empty>\n")
		} else {
			fmt.Printf("    Dependencies:\n")
			for _, d := range dep {
				fmt.Printf("        %-40s %s\n", d.Name(), d.Version().String())
			}
		}

		rem := Status.Repository.Remotes()
		if len(rem) == 0 {
			fmt.Printf("    Remotes     : <empty>\n")
		} else {
			fmt.Printf("    Remotes     :\n")
			for _, r := range rem {
				u := r.Path()
				tr := "" // str repr of the token
				t := r.Token()
				if t != nil { // applies only if not nul
					tr = t.FormatStd()
				}

				fmt.Printf("        %-40s %-40s %s\n", r.Name(), u.String(), tr)
			}
		}

	},
}

var Add = Command{
	Name:      `add`,
	Alias:     `+`,
	UsageLine: `NAME VERSION`,
	Short:     `Add a dependency.`,
	Long: `Add a dependency.

  NAME     dependency package name
  VERSION  a semantic version  
`,
	RequireProject: true,
	Run: func(Add *Command) {

		if len(Add.Flag.Args()) != 2 {
			Add.Flag.Usage()
			return
		}
		name, version := Add.Flag.Arg(0), Add.Flag.Arg(1)
		v, _ := ParseVersion(version)
		ref := NewProjectID(name, v)
		fmt.Printf("  -> %v\n", ref)

		Add.Project.AppendDependency(ref)

		Add.Project.Write()
	},
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var listSnapshotOnlyFlag *bool
var List = Command{
	Name:           `list`,
	Alias:          `l`,
	UsageLine:      ``,
	Short:          `List Dependencies.`,
	Long:           `List all dependencies in a dependency order. Meaning that for any given package, all its dependencies are listed before it.`,
	RequireProject: true,
	FlagInit: func(List *Command) {
		listSnapshotOnlyFlag = List.Flag.Bool("s", false, "snapshot only. List only snapshot dependencies")
	},
	Run: func(List *Command) {
		if *listSnapshotOnlyFlag {
			TitleStyle.Printf("Snapshot Dependencies for %s:\n", List.Project.Name())
		} else {
			TitleStyle.Printf("All Dependencies for %s:\n", List.Project.Name())
		}

		dependencies, err := List.Repository.ResolveDependencies(List.Project, true, false)

		if err != nil {
			ErrorStyle.Printf("Cannot resolve dependencies. %s\n", err)
		}
		all := !*listSnapshotOnlyFlag
		for _, d := range dependencies {
			if all || d.Version().IsSnapshot() {
				SuccessStyle.Printf("        %-40s %s\n", d.Name(), d.Version().String())
			}
		}
	},
}

var searchRemoteFlag *string
var Search = Command{
	Name:           `search`,
	Alias:          `s`,
	UsageLine:      `QUERY`,
	Short:          `Search Packages .`,
	Long:           `Search Packages in the local repository that starts with the QUERY`,
	RequireProject: false,
	FlagInit: func(Search *Command) {
		searchRemoteFlag = Search.Flag.String("r", "", "remote. Search in the remote REMOTE instead")
	},
	Run: func(Search *Command) {

		search := Search.Flag.Arg(0)
		var result []protocol.PID

		if *searchRemoteFlag != "" { // find the remote and use it
			rem := *searchRemoteFlag
			remote, err := Search.Repository.Remote(rem)
			if err != nil {
				ErrorStyle.Printf("Unknown remote %s.\n", rem)
				fmt.Printf("Available remotes are:\n")
				for _, r := range Search.Repository.Remotes() {
					u := r.Path()
					fmt.Printf("    %-40s %s\n", r.Name(), u.String())
				}
				return
			}
			result = remote.Search(search, 0)
		} else {
			result = Search.Repository.Search(search, 0)
		}
		// result contains the acual results every error should have been processed

		pkg := "" //(to avoid printing again and again the package name
		for _, pid := range result {
			currentPackage := pid.Name
			if currentPackage == pkg { // do not print if not new
				currentPackage = ""
			} else {
				pkg = currentPackage // remember it
			}

			SuccessStyle.Printf("    %-40s %s\n", currentPackage, pid.Version.String())
		}
	},
}

////////////////////////////////////////////////////////////////////////////////////////

var AddRemote = Command{
	Name:      `add-remote`,
	Alias:     `r+`,
	UsageLine: `NAME URL [TOKEN]`,
	Short:     `Add a remote server.`,
	Long: `Remote server can be used to publish or receive other's code.
  NAME    local alias for this remote
  URL     full URL to the remote server. file:// and http:// are actually supported. for http, see 'gpk serve'
  TOKEN   option: a secret TOKEN identification provided by the server to deliver authentication. See with your server`,
	RequireProject: false,
	Run: func(AddRemote *Command) {

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
		SuccessStyle.Printf("new remote: %s %s %s\n", name, u, token)
		AddRemote.Repository.RemoteAdd(client)
		AddRemote.Repository.Write()
	},
}

var Remove = Command{
	Name:           `rem`,
	Alias:          `-`,
	UsageLine:      `NAME`,
	Short:          `Remove dependency`,
	Long:           ``,
	RequireProject: true,
	Run: func(Remove *Command) {

		if len(Remove.Flag.Args()) != 1 {
			Remove.Flag.Usage()
			return
		}
		name := Remove.Flag.Arg(0)
		ref := Remove.Project.RemoveDependency(name)
		if ref != nil {
			SuccessStyle.Printf("Removed Dependency %s %s\n", ref.Name(), ref.Version().String())
			Remove.Project.Write()
		} else {
			ErrorStyle.Printf("Nothing to remove %s\n")
		}
	},
}

var RemoveRemote = Command{
	Name:           `rem-remote`,
	Alias:          `r-`,
	UsageLine:      `NAME`,
	Short:          `Remove a Remote`,
	Long:           ``,
	RequireProject: false,
	Run: func(RemoveRemote *Command) {

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
			SuccessStyle.Printf("Removed Remote %s %s\n", ref.Name(), ref.Path())
			Remove.Project.Write()
		} else {
			ErrorStyle.Printf("Nothing to Remove\n")
		}
		RemoveRemote.Repository.Write()
	},
}
