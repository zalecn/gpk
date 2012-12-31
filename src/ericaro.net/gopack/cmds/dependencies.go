package cmds

import (
	. "ericaro.net/gopack"
	. "ericaro.net/gopack/semver"
)

func init() {
	Reg(
		&Add,
		&Remove,
	)
}

var Add = Command{
	Name:      `dadd`,
	Alias:     `d+`,
	Category : DependencyCategory,
	UsageLine: `NAME VERSION`,
	Short:     `Add dependency`,
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
		ref := *NewProjectID(name, v)
		SuccessStyle.Printf("  -> %v\n", ref)

		Add.Project.AppendDependency(ref)

		Add.Project.Write()
	},
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var Remove = Command{
	Name:           `dremove`,
	Alias:          `d-`,
	Category : DependencyCategory,
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
