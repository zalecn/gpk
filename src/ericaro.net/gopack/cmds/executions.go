package cmds

//Executions commands

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/gocmd"
	"fmt"
	"os"
	"time"
)

func init() {
	Reg(
		&Compile,
		&Test,
	)

}

// The flag package provides a default help printer via -h switch
var compileAllFlag *bool
var compileOfflineFlag *bool
var compileUpdateFlag *bool
var Compile = Command{
	Name:      `compile`,
	Alias:     `c`,
	Category:  CompileCategory,
	UsageLine: ``,
	Short:     `Compile project`,
	Long: `Computes current project dependencies as a GOPATH variable (accessible through the p Option),
       and then run go install on the project.`,
	RequireProject: true,
	FlagInit: func(Compile *Command) {
		compileAllFlag = Compile.Flag.Bool("a", false, "all. Force rebuilding of packages that are already up-to-date.")
		compileOfflineFlag = Compile.Flag.Bool("o", false, "offline. Try to find missing dependencies at http://gpk.ericaro.net")
		compileUpdateFlag = Compile.Flag.Bool("u", false, "update. Look for updated version of dependencies")
	},
	Run: func(Compile *Command) (err error) {
		// parse dependencies, and build the gopath
		dependencies, err := Compile.Repository.ResolveDependencies(Compile.Project, *compileOfflineFlag, *compileUpdateFlag)
		if err != nil {
			ErrorStyle.Printf("Error Resolving project's dependencies:\n    \u21b3 %v", err)
			return
		}
		// run the go build command for local src, and with the appropriate gopath
		gopath, err := Compile.Repository.GoPath(dependencies)

		goEnv := gocmd.NewGoEnv(gopath)
		err = goEnv.Install(Compile.Project.WorkingDir(), *compileAllFlag) // TODO finalize the effort to wrap all the go install command (even maybe go build)
		// also provide a go run equivalent 
		return
	},
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////:

var testWatchFlag *time.Duration
var Test = Command{
	Name:           `test`,
	Alias:          `t`,
	Category:       CompileCategory,
	UsageLine:      ``,
	Short:          `Run go test`,
	Long:           `Run go test on the whole project.`, // TODO add options to select the package to be executed
	RequireProject: true,
	FlagInit: func(Compile *Command) {
		testWatchFlag = Compile.Flag.Duration("w", 0, "watch. Repeat the command for ever every watched seconds")
	},
	Run: func(Test *Command) (err error) {

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
		// make two cases, either I'm on the root dir, then the package is ./src/... or I'm within the project 
		// and the packages path is ./...

		wd, err := os.Getwd()
		// use cwd as root, unless it is in the workding dir
		var wdi, pwdi os.FileInfo
		wdi, err = os.Stat(wd)
		pwdi, err = os.Stat(Test.Project.WorkingDir())

		goEnv := gocmd.NewGoEnv(gopath)
		var args []string
		if os.SameFile(wdi, pwdi) { // if I'm on the project root dir, sources are in ./src/...
			args = []string{"./src/..."}
		} else { // otherwise assume that you know what you are doing, and just recuse in ./...
			args = []string{"./..."}
		}
		// check for -watch option

		if *testWatchFlag > 0 {

			c := time.Tick(100 * time.Millisecond)
			for {
				next := time.Now()
				for now := range c {
					if now.After(next) {
						next = now.Add(*testWatchFlag)
						NormalStyle.Clear()
						err = goEnv.Test(Test.Project.WorkingDir(), wd, args)
						fmt.Printf("\n")
					}
					fmt.Printf("■")
				}
			}

		}

		return
	},
}
