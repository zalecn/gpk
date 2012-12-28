package cmds

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/gocmd"
	"ericaro.net/gopack/semver"
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

var pathListFlag *bool
var Path = Command{
	Name:      `path`,
	Alias:     `p`,
	UsageLine: ``,
	Short:     `Compute project's GOPATH`,
	Long: `Resolve current project dependencies and print the GOPATH variable.
    
    ` + "Tip: type:\n    alias GP='export GOPATH=`gpk path`'\n\n    to get a simple automatic GOPATH setting.",
	RequireProject: true,
	FlagInit: func(Path *Command) {
		pathListFlag = Path.Flag.Bool("l", false, "list. Pretty Print the list.")
	},
	Run: func(Path *Command) {

		// parse dependencies, and build the gopath
		dependencies, err := Path.Repository.ResolveDependencies(Path.Project, true, false) // path does not update the dependencies
		if err != nil {
			ErrorStyle.Printf("Error Resolving project's dependencies:\n    \u21b3 %v", err)
			return
		}

		if *pathListFlag {
			fmt.Printf("    %s dependencies are:\n\n", ShortStyle.Sprintf("%s", Path.Project.Name()))
			// run the go build command for local src, and with the appropriate gopath
			TitleStyle.Printf("    %-50s -> %s\n", "Dependency", "Path")
			for _, d := range dependencies {
				fmt.Printf("    %-50s -> %s\n", d.ID().String(), d.InstallDir())
			}
			fmt.Println()
		} else {
			gopath, _ := Path.Repository.GoPath(dependencies)
			fmt.Print(gocmd.Join(Path.Project.WorkingDir(), gopath)) // there is no line break because this way we can use it to initialize a local GOPATH var
		}
	},
}

// The flag package provides a default help printer via -h switch
var compileAllFlag *bool
var compileOfflineFlag *bool
var compileUpdateFlag *bool
var Compile = Command{
	Name:      `compile`,
	Alias:     `x`,
	UsageLine: ``,
	Short:     `Compile project`,
	Long: `Computes current project dependencies as a GOPATH variable (accessible through the p Option),
    and then run go install on the project.`,
	RequireProject: true,
	FlagInit: func(Compile *Command) {
		compileAllFlag = Compile.Flag.Bool("a", false, "all. Go standard build option. Force rebuilding of packages that are already up-to-date.")
		compileOfflineFlag = Compile.Flag.Bool("o", false, "offline. Try to find missing dependencies at http://gpk.ericaro.net")
		compileUpdateFlag = Compile.Flag.Bool("u", false, "update. Look for updated version of dependencies")
	},
	Run: func(Compile *Command) {
		// parse dependencies, and build the gopath
		dependencies, err := Compile.Repository.ResolveDependencies(Compile.Project, *compileOfflineFlag, *compileUpdateFlag)
		if err != nil {
			ErrorStyle.Printf("Error Resolving project's dependencies:\n    \u21b3 %v", err)
			return
		}
		// run the go build command for local src, and with the appropriate gopath
		gopath, err := Compile.Repository.GoPath(dependencies)

		goEnv := gocmd.NewGoEnv(gopath)
		goEnv.Install(Compile.Project.WorkingDir(), *compileAllFlag) // TODO finalize the effort to wrap all the go install command (even maybe go build)
		// also provide a go run equivalent 
	},
}

var Install = Command{
	Name:      `install`,
	Alias:     `i`,
	UsageLine: `VERSION`,
	Short:     `Install into the local repository`,
	Long: `Install the current project sources in the local repository.
    
    VERSION is a semantic version to identify this specific project version. See http://semver.org for more details about semantic version
`,
	RequireProject: true,
	Run: func(Install *Command) {
		version, err := semver.ParseVersion(Install.Flag.Arg(0))
		if err != nil {
			ErrorStyle.Printf("Syntax error on Version %s\n", Install.Flag.Arg(0))
			return
		}
		Install.Repository.InstallProject(Install.Project, version)
	},
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////:

var Test = Command{
	Name:           `test`,
	Alias:          `t`,
	UsageLine:      ``,
	Short:          `Run go test`,
	Long:           `Compute current project dependencies as a GOPATH variable, and then run go test on the whole project.`,
	RequireProject: true,
	Run: func(Test *Command) {

		// parse dependencies, and build the gopath
		dependencies, err := Test.Repository.ResolveDependencies(Test.Project, true, false)
		if err != nil {
			ErrorStyle.Printf("Error Resolving project's dependencies:\n    \u21b3 %v", err)
			return
		}

		// run the go build command for local src, and with the appropriate gopath
		gopath, err := Test.Repository.GoPath(dependencies)
		if err != nil {
			ErrorStyle.Printf("Invalid dependency:\n    \u21b3 %v", err)
			return
		}

		goEnv := gocmd.NewGoEnv(gopath)
		goEnv.Test(Test.Project.WorkingDir())
	},
}
