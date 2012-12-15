package main

import (
	"ericaro.net/gopackage"
	"fmt"
	"os"
)

func init() {
	Reg(
		&Status,
		&Add,
		&Remove,
		&Init,
	)

}

var Status = Command{
	Name:           `status`,
	Alias:          `?`,
	UsageLine:      ``,
	Short:          `Prints current directory project status.`,
	Long:           `Prints current directory status`,
	call:           func(c *Command) { c.Status() },
	RequireProject: true,
}

var Add = Command{
	Name:      `add`,
	Alias:     `+`,
	UsageLine: `<name> <version>`,
	Short:     `Add a dependency to this project.`,
	Long: `Dependency is formatted as follow
	name    : any string
	version : is a version syntax
	          <root>-X.X.X.X
	          where:
	          root : is a simple branch name
	          X    : is an unsigned int  
	`,
	call:           func(c *Command) { c.Add() },
	RequireProject: true,
}
var Remove = Command{
	Name:      `remove`,
	Alias:     `-`,
	UsageLine: `<name> <version>`,
	Short:     `Remove dependency from this project.`,
	Long: `Dependencies are formatted as follow
	where   :
	name    : any string
	artifact: is a simple name (usually not hierarchical )
	version : is a version syntax
	          <root>-X.X.X.X
	          where:
	          root : is a simple branch name
	          X    : is an unsigned int  
	`,
	call:           func(c *Command) { c.Remove() },
	RequireProject: true,
}

var Init = Command{
	Name:      `init`,
	Alias:     `!`,
	UsageLine: `<name>`,
	Short:     `Init the current directory as a gopackage project.`,
	Long: `where   :
	name   : any string`,
	call: func(c *Command) { c.Init() },
}

//var statusX *bool = status.Flag.Bool("x",false, "test" )

func (c *Command) Status() {

	fmt.Printf("%v\n", c.Project)
}

func (c *Command) Init() {

	p, err := gopackage.ReadProject()
	if err == nil {
		fmt.Printf("warning: init an existing project. This is fine if you want to edit it\n")
	}
	p.Name = c.Flag.Arg(0)
	c.Project = p // in case we implement sequence of commands (in the future)
	p.Root, err = os.Getwd()
	if err!=nil {
		fmt.Printf("Cannot create the project, there is no current directory. Because %v\n", err)
	}
	gopackage.WriteProjectSrc(c.Project)

}

func (c *Command) Add() {

	if len(c.Flag.Args()) != 2 {
		c.Flag.Usage()
		return
	}
	name, version := c.Flag.Arg(0), c.Flag.Arg(1)
	ref := gopackage.NewProjectReference(name, gopackage.ParseVersionReference(version))
	fmt.Printf("  -> %v\n", ref)
	c.Project.AppendDependency(ref)

	gopackage.WriteProjectSrc(c.Project)
}

func (c *Command) Remove() {

	if len(c.Flag.Args())  != 2 {
		c.Flag.Usage()
		return
	}
	name, version := c.Flag.Arg(0), c.Flag.Arg(1)
	ref := gopackage.NewProjectReference(name, gopackage.ParseVersionReference(version))

	fmt.Printf("  -> %v\n", ref)
	c.Project.RemoveDependency(ref)
	gopackage.WriteProjectSrc(c.Project)
}
