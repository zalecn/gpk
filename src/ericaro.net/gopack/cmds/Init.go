package cmds

import (
	. "ericaro.net/gopack"
	"fmt"
	"os"
)

func init() {
	Reg(
		&Init,
	)

}

var initNameFlag *string
var initLicenseFlag *string

var Init = Command{
	Name:      `init`,
	Alias:     `!`,
	Category:  InitCategory,
	UsageLine: `-n NAME -l LICENSE`,
	Short:     `Initialize or Edit the current project`,
	Long: `Init the current directory creates or updates the gopack project file, 
       and name the current package NAME, setting the license to LICENSE.
       
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
	Run: func(Init *Command) (err error) {
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
			SuccessStyle.Printf("new name:%s\n", p.Name())
		}

		if *initLicenseFlag != "" {
			var lic *License

			// sorry for that, this is ugly plumbing, I'll come back and fix later
			if l, err := Licenses.Get(*initLicenseFlag); err != nil {
				if l, err = Licenses.GetAlias(*initLicenseFlag); err != nil {
					ErrorStyle.Printf("new license: unknown or unsupported license:%s\n", p.License())
				} else {
					lic = l
				}
			} else {
				lic = l
			}
			// end of sorryness 

			if lic != nil {
				p.SetLicense(*lic)
				SuccessStyle.Printf("new license:\"%s\"\n", p.License().FullName)
			}
		}

		Init.Project = p // in case we implement sequence of commands (in the future)
		p.Write()        // store it  one day I'll implement a lock on this file, right ?
		return
	},
}
