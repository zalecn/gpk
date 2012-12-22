package cmds

import (
	. "ericaro.net/gpk"
	"ericaro.net/gpk/gocmd"
	"fmt"
)

func init() {
	Reg(
		&Path,
		&Compile,
		&Install,
		&Test,
	)

}

var Path = Command{
	Name:      `path`,
	Alias:     `p`,
	UsageLine: ``,
	Short:     `Compute project's GOPATH`,
	Long: `Resolve current project dependencies and print the GOPATH variable.
Used with -l it provides a pretty print of the list.
    
` + "Tip: type \n\n        alias GP='export GOPATH=`gpk path`'\n\n    this is a simple way to automatically export your GOPATH.",
	call:           func(c *Command) { c.Path() },
	RequireProject: true,
}

var Compile = Command{
	Name:      `compile`,
	Alias:     `x`,
	UsageLine: ``,
	Short:     `Compile project`,
	Long: `Computes current project dependencies as a GOPATH variable (accessible through the p Option),
and then run go install on the project.`,
	call:           func(c *Command) { c.Compile() },
	RequireProject: true,
}

var Install = Command{
	Name:      `install`,
	Alias:     `i`,
	UsageLine: `<version>`,
	Short:     `Install into the local repository`,
	Long: `Make the current project sources available in the local repository under the ID <package>:<version>.
	
<version> denotes a semantic version to identify this specific project version.
    
Installing to the local repository  will replace any existing code in the local repository. 
This is a "snapshot" behavior: installing v1 will replace existing v1.
    
If you are interested in handling "read only" version (released version) please consider using "gpk deploy" instead.  
    `,
	call:           func(c *Command) { c.Install() },
	RequireProject: true,
}

var Test = Command{
	Name:      `test`,
	Alias:     `t`,
	UsageLine: ``,
	Short:     `Run go test`,
	Long: `Compute current project dependencies as a GOPATH variable (accessible through the p Option),
and then run go test on the whole project.`,
	call:           func(c *Command) { c.Test() },
	RequireProject: true,
}

// The flag package provides a default help printer via -h switch
var compileOfflineFlag *bool = Compile.Flag.Bool("o", false, "offline. Try to find missing dependencies at http://gpk.ericaro.net")
var compileUpdateFlag *bool = Compile.Flag.Bool("u", false, "update. Look for updated version of dependencies")

func (c *Command) Compile() {

	// parse dependencies, and build the gopath
	// todo remote should be read from the project
	dependencies, err := c.Repository.ResolveDependencies(c.Project, *compileOfflineFlag, *compileUpdateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies: %v", err)
		return
	}

	// run the go build command for local src, and with the appropriate gopath
	gopath, err := c.Repository.GoPath(dependencies)

	goEnv := gocmd.NewGoEnv(gopath)
	goEnv.Install(c.Project.WorkingDir())

}

var pathListFlag *bool = Path.Flag.Bool("l", false, "list. Print the list is a readable way.")

func (c *Command) Path() {

	// parse dependencies, and build the gopath
	// todo remote should be read from the project
	dependencies, err := c.Repository.ResolveDependencies(c.Project, *compileOfflineFlag, *compileUpdateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies: %v", err)
		return
	}
	if *pathListFlag {
		fmt.Printf("    %s dependencies are:\n\n", ShortStyle.Sprintf("%s",c.Project.Name() ) )
		// run the go build command for local src, and with the appropriate gopath
		TitleStyle.Printf("    %-50s -> %s\n", "Dependency", "Path")
		for _, d := range dependencies {
			fmt.Printf("    %-50s -> %s\n", d.ID().String(), d.InstallDir())
		}
		fmt.Println()
	} else {
		gopath, _ := c.Repository.GoPath(dependencies)
		fmt.Print(gocmd.Join(c.Project.WorkingDir(), gopath))
	}
}

func (c *Command) Install() {
	version, _ := ParseVersion(c.Flag.Arg(0))
	c.Repository.InstallProject(c.Project, version)
}

func (c *Command) Test() {

	// parse dependencies, and build the gopath
	dependencies, err := c.Repository.ResolveDependencies(c.Project, *compileOfflineFlag, *compileUpdateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies: %v", err)
		return
	}

	// run the go build command for local src, and with the appropriate gopath
	gopath, err := c.Repository.GoPath(dependencies)
	goEnv := gocmd.NewGoEnv(gopath)
	goEnv.Test(c.Project.WorkingDir())
}
