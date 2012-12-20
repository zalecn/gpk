package main

import (
	. "ericaro.net/gpk"
	"fmt"
	"os"
	"net/url"
)

func init() {
	Reg(
		&Status,
		&Add,
		&Remove,
		&AddRemote,
		&RemoveRemote,
		&Init,
	)

}

var Status = Command{
	Name:           `status`,
	Alias:          `?`,
	UsageLine:      ``,
	Short:          `Print status`,
	Long:           `Print status`,
	call:           func(c *Command) { c.Status() },
	RequireProject: true,
}

var Add = Command{
	Name:      `add`,
	Alias:     `+`,
	UsageLine: `<name> <version>`,
	Short:     `Add a dependency.`,
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
var AddRemote = Command{
	Name:      `add-remote`,
	Alias:     `r+`,
	UsageLine: `<name> <url>`,
	Short:     `Add a remote server.`,
	Long: `Remote server can be used to publish or share code snapshots  
	`,
	call:           func(c *Command) { c.AddRemote() },
	RequireProject: true,
}
var Remove = Command{
	Name:      `rem`,
	Alias:     `-`,
	UsageLine: `<name> <version>`,
	Short:     `Remove dependency`,
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
var RemoveRemote = Command{
	Name:           `rem-remote`,
	Alias:          `r-`,
	UsageLine:      `<name>`,
	Short:          `Remove remote synchronization server`,
	Long:           ``,
	call:           func(c *Command) { c.RemoveRemote() },
	RequireProject: true,
}

var Init = Command{
	Name:      `init`,
	Alias:     `!`,
	UsageLine: `<name>`,
	Short:     `Init the current directory as a go package kit project.`,
	Long: `where   :
	name   : any string`,
	call: func(c *Command) { c.Init() },
}

//var statusX *bool = status.Flag.Bool("x",false, "test" )

func (c *Command) Status() {

	TitleStyle.Printf("Project %s:\n\n", c.Project.Name())
	fmt.Printf("    Depends on\n")
	for _, d := range c.Project.Dependencies() {
		fmt.Printf("        %-40s %s\n", d.Name(), d.Version().String())
	}

	fmt.Printf("\n    Synchronizes with\n")
	for _, r := range c.Project.Remotes() {
		u := r.Path()
		fmt.Printf("        %-40s %v\n", r.Name(), u.String())
	}

}

func (c *Command) Init() {

	p, err := ReadProject()
	if err == nil {
		fmt.Printf("warning: init an existing project. This is fine if you want to edit it\n")
	}
	p.SetName(c.Flag.Arg(0))
	pwd, err := os.Getwd()
	p.SetWorkingDir(pwd)

	c.Project = p // in case we implement sequence of commands (in the future)

	if err != nil {
		fmt.Printf("Cannot create the project, there is no current directory. Because %v\n", err)
	}
	p.Write()

}

func (c *Command) Add() {

	if len(c.Flag.Args()) != 2 {
		c.Flag.Usage()
		return
	}
	name, version := c.Flag.Arg(0), c.Flag.Arg(1)
	v, _ := ParseVersion(version)
	ref := NewProjectID(name, v)
	fmt.Printf("  -> %v\n", ref)

	c.Project.AppendDependency(ref)

	c.Project.Write()
}

func (c *Command) AddRemote() {

	if len(c.Flag.Args()) != 2 {
		c.Flag.Usage()
		return
	}
	name, remote := c.Flag.Arg(0), c.Flag.Arg(1)
	u, err := url.Parse(remote)
	if err != nil {
		ErrorStyle.Printf("Invalid URL passed as a remote Repository.\n    Caused by %s\n", err)
		return
	}
	c.Project.RemoteAdd(NewRemoteRepository(name, *u))
	c.Project.Write()
}

func (c *Command) Remove() {

	if len(c.Flag.Args()) != 2 {
		c.Flag.Usage()
		return
	}
	name, version := c.Flag.Arg(0), c.Flag.Arg(1)
	v, _ := ParseVersion(version)
	ref := NewProjectID(name, v)

	fmt.Printf("  -> %v\n", ref)
	c.Project.RemoveDependency(ref)
	c.Project.Write()
}
func (c *Command) RemoveRemote() {

	if len(c.Flag.Args()) < 1 {
		c.Flag.Usage()
		return
	}
	for _, name := range c.Flag.Args() {
		c.Project.RemoteRemove(name)
	}
	c.Project.Write()
}
