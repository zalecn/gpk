package cmds

import (
	. "ericaro.net/gopack"
	"fmt"
)

func init() {
	Reg(
		&Status,
	)

}
var statusNameFlag *bool
var Status = Command{
	Name:           `status`,
	Alias:          `?`,
	Category:       InitCategory,
	UsageLine:      ``,
	Short:          `Print status`,
	Long:           `Display current information about the current project and the current local repository`,
	RequireProject: true,
	FlagInit: func(Imports *Command) {
		statusNameFlag = Imports.Flag.Bool("n", false, "name. Print only the name")
		

	},
	Run: func(Status *Command)  (err error){
		
		if *statusNameFlag{
			fmt.Print(Status.Project.Name() )
			return
		}
		
		TitleStyle.Printf("    Name        : %s\n", Status.Project.Name())
		SuccessStyle.Printf("    License     : %s\n", Status.Project.License().FullName)
		dep := Status.Project.Dependencies()
		if len(dep) == 0 {
			SuccessStyle.Printf("    Dependencies: <empty>\n")
		} else {
			SuccessStyle.Printf("    Dependencies:\n")
			for _, d := range dep {
				SuccessStyle.Printf("        %-40s %s\n", d.Name(), d.Version().String())
			}
		}

		rem := Status.Repository.Remotes()
		if len(rem) == 0 {
			SuccessStyle.Printf("    Remotes     : <empty>\n")
		} else {
			SuccessStyle.Printf("    Remotes     :\n")
			for _, r := range rem {
				u := r.Path()
				tr := "" // str repr of the token
				t := r.Token()
				if t != nil { // applies only if not nul
					tr = t.FormatStd()
				}

				SuccessStyle.Printf("        %-40s %-40s %s\n", r.Name(), u.String(), tr)
			}
		}
		return
	},
}
