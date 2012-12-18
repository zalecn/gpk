package main

import (
	. "ericaro.net/gpk"
	"ericaro.net/gpk/gocmd"
	"fmt"
)

func init() {
	Reg(
		&Compile,
		&Install,
		&Test,
		&Deploy,
		&Get,
	)

}

var Compile = Command{
	Name:           `compile`,
	Alias:          `x`,
	UsageLine:      ``,
	Short:          `Compile the current project`,
	Long:           `Computes current project dependencies as a GOPATH variable (accessible through the p Option), and then compile the project`,
	call:           func(c *Command) { c.Compile() },
	RequireProject: true,
}

var Install = Command{
	Name:           `install`,
	Alias:          `i`,
	UsageLine:      `<version>`,
	Short:          `Install the current project in the local repository`,
	Long:           `Install the current project in the local repository`,
	call:           func(c *Command) { c.Install() },
	RequireProject: true,
}

var Test = Command{
	Name:           `test`,
	Alias:          `t`,
	UsageLine:      ``,
	Short:          `Run test on the current project`,
	Long:           `call go test on the current project.`,
	call:           func(c *Command) { c.Test() },
	RequireProject: true,
}

// TODO move around the remote tool chain
var Deploy = Command{
	Name:           `deploy`,
	Alias:          `d`,
	UsageLine:      `<version>`,
	Short:          `Deploy the current project in the remote repository`,
	Long:           `Deploy the current project in the remote repository`,
	call:           func(c *Command) { c.Deploy() },
	RequireProject: true,
}

var Get = Command{
	Name:           `goget`,
	Alias:          `get`,
	UsageLine:      `<goget package>`,
	Short:          `Run go get a package and install it`,
	Long:           `Run go get a package and install it`,
	call:           func(c *Command) { c.Get() },
	RequireProject: false,
}

// The flag package provides a default help printer via -h switch
var compileVersionFlag *bool = Compile.Flag.Bool("v", false, "Print the version number.")
var compileOfflineFlag *bool = Compile.Flag.Bool("o", false, "Try to find missing dependencies at http://gpk.ericaro.net")
var compileUpdateFlag *bool = Compile.Flag.Bool("u", false, "Look for updated version of dependencies")
var compilePathOnlyFlag *bool = Compile.Flag.Bool("p", false, fmt.Sprintf("Does not run the compile, just print the gopath (suitable for scripting for instance: alias GP='export GOPATH=`%s compile -p`' )", Cmd))

func (c *Command) Compile() {

	// parse dependencies, and build the gopath
	// todo remote should be read from the project
	remote, _ := NewHttpRemoteRepository(GopackageCentral)
	dependencies, err := c.Repository.ResolveDependencies(c.Project, remote, *compileOfflineFlag, *compileUpdateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies: %v", err)
		return
	}

	// run the go build command for local src, and with the appropriate gopath
	gopath, err := c.Repository.GoPath(dependencies)

	if *compilePathOnlyFlag {
		fmt.Print(gocmd.Join(c.Project.WorkingDir(), gopath))
		return
	} else {
		goEnv := gocmd.NewGoEnv(gopath)
		goEnv.Install(c.Project.WorkingDir())
	}

}

func (c *Command) Install() {
	version, _ := ParseVersion(c.Flag.Arg(0))
	c.Repository.InstallProject(c.Project, version)
}

func (c *Command) Deploy() {

	//version, _ := ParseVersion(c.Flag.Arg(0))
	panic("not yet implemented")
	// TODO select the remote and then deploy
	//c.Repository.DeployProject(c.Project, version)
}

func (c *Command) Get() {
		panic("not yet implemented")
//	for _, p := range c.Flag.Args() {
//		//TODO fill it
//		//c.Repository.GoGetInstall(p)
//	}
}

func (c *Command) Test() {

	// parse dependencies, and build the gopath
	remote, _ := NewHttpRemoteRepository(GopackageCentral)
	dependencies, err := c.Repository.ResolveDependencies(c.Project, remote, *compileOfflineFlag, *compileUpdateFlag)
	if err != nil {
		fmt.Printf("Error Parsing the project's dependencies: %v", err)
		return
	}

	// run the go build command for local src, and with the appropriate gopath
	gopath, err := c.Repository.GoPath(dependencies)
	goEnv := gocmd.NewGoEnv(gopath)
	goEnv.Test(c.Project.WorkingDir())
}
