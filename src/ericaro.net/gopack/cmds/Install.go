package cmds

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/semver"
)

func init() {
	Reg(
		&Install,
	)

}

var Install = Command{
	Name:      `install`,
	Alias:     `i`,
	UsageLine: `VERSION`,
	Short:     `Install into the local repository`,
	Long: `Install the current project sources in the local repository.
       
       VERSION is a semantic version to identify this specific project version.
       See http://semver.org for more details about semantic versions.
`,
	RequireProject: true,
	Run: func(Install *Command)  (err error){
	
		if len(Install.Flag.Args()) !=1 {
			ErrorStyle.Printf("Missing version arguments\n")
			NormalStyle.Printf("       gpk install VERSION\n")
			return InvalidArgumentSize()
			return
		}
		version, err := semver.ParseVersion(Install.Flag.Arg(0))
		if err != nil {
			ErrorStyle.Printf("Syntax error on Version %s\n", Install.Flag.Arg(0))
			return
		}
		Install.Repository.InstallProject(Install.Project, version)
		return
	},
}
