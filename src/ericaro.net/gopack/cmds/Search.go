package cmds

import (. "ericaro.net/gopack"
	"ericaro.net/gopack/protocol"
	"fmt"
	)

func init() {
	Reg(
		&Search,
	)

}

var searchRemoteFlag *string
var Search = Command{
	Name:           `search`,
	Alias:          `s`,
	UsageLine:      `QUERY`,
	Short:          `Search Packages .`, //TODO add the search import capability
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
