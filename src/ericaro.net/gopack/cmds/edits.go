package cmds

import (
	. "ericaro.net/gopack"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

func init() {
	Reg(
		&Init,
		&Status,
		&Add,
		&Remove,
		&List,
		&Search,
		&AddRemote,
		&RemoveRemote,
	)

}

var Init = Command{
	Name:      `init`,
	Alias:     `!`,
	UsageLine: `NAME`,
	Short:     `Initialize current directory.`,
	Long: `In the current directory creates or updates the gopack project file, and name the current package NAME
`,
	call: func(c *Command) { c.Init() },
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







var Status = Command{
	Name:           `status`,
	Alias:          `?`,
	UsageLine:      ``,
	Short:          `Print status`,
	Long:           ``,
	call:           func(c *Command) { c.Status() },
	RequireProject: true,
}

var Add = Command{
	Name:      `add`,
	Alias:     `+`,
	UsageLine: `<name> <version>`,
	Short:     `Add a dependency.`,
	Long: `Dependency is formatted as follow

name
    any string
version
    a version definition with Semantic Version syntax  

`,
	call:           func(c *Command) { c.Add() },
	RequireProject: true,
}
var List = Command{
	Name:           `list`,
	Alias:          `l`,
	UsageLine:      ``,
	Short:          `List Dependencies.`,
	Long:           `List all dependencies in a dependency order. Meaning that for any given package, all its dependencies are listed before it.`,
	call:           func(c *Command) { c.List() },
	RequireProject: true,
}
var Search = Command{
	Name:           `search`,
	Alias:          `s`,
	UsageLine:      `<search>`,
	Short:          `Search Packages.`,
	Long:           `Search Packages in the local repository that starts with the query`,
	call:           func(c *Command) { c.Search() },
	RequireProject: false,
}
var AddRemote = Command{
	Name:      `add-remote`,
	Alias:     `r+`,
	UsageLine: `<name> <url>`,
	Short:     `Add a remote server.`,
	Long: `Remote server can be used to publish or share code snapshots`,
	call:           func(c *Command) { c.AddRemote() },
	RequireProject: false,
}
var Remove = Command{
	Name:      `rem`,
	Alias:     `-`,
	UsageLine: `<name> <version>`,
	Short:     `Remove dependency`,
	Long: `Dependencies are formatted as follow:

name
    any string
artifact
    is a simple name (usually not hierarchical )
version
    is a version definition in semantic version syntax.
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
	RequireProject: false,
}


//var statusX *bool = status.Flag.Bool("x",false, "test" )

func (c *Command) Status() {

	TitleStyle.Printf("Project %s:\n\n", c.Project.Name())
	fmt.Printf("    Depends on\n")
	for _, d := range c.Project.Dependencies() {
		fmt.Printf("        %-40s %s\n", d.Name(), d.Version().String())
	}

	fmt.Printf("\n    Synchronizes with\n")
	for _, r := range c.Repository.Remotes() {
		u := r.Path()
		fmt.Printf("        %-40s %v\n", r.Name(), u.String())
	}

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

var listSnapshotOnlyFlag *bool = List.Flag.Bool("s", false, "snapshot only. List only snapshot dependencies")

func (c *Command) List() {
	if *listSnapshotOnlyFlag {
		TitleStyle.Printf("Snaphost Dependencies for %s:\n", c.Project.Name())
	} else {
		TitleStyle.Printf("All Dependencies for %s:\n", c.Project.Name())
	}
	dependencies, err := c.Repository.ResolveDependencies(c.Project, true, false)
	if err != nil {
		ErrorStyle.Printf("Cannot resolve dependencies. %s\n", err)
	}
	all := !*listSnapshotOnlyFlag
	for _, d := range dependencies {
		if all || d.Version().IsSnapshot() {
			fmt.Printf("        %-40s %s\n", d.Name(), d.Version().String())
		}
	}
}

var searchRemoteFlag *string = Search.Flag.String("r", "", "Search in the remote <remote> instead")

func (c *Command) Search() {
	search := c.Flag.Arg(0)
	var result []string
	if *searchRemoteFlag != "" {
		rem := *searchRemoteFlag
		remote, err := c.Repository.Remote(rem)
		if err != nil {
			ErrorStyle.Printf("Unknown remote %s.\n", rem)

			fmt.Printf("Available remotes are:\n")
			for _, r := range c.Repository.Remotes() {
				u := r.Path()
				fmt.Printf("    %-40s %s\n", r.Name(), u.String())
			}
			return
		}
		result = remote.SearchPackage(search)
	} else {
		result = c.Repository.SearchPackage(search)
	}

	pkg := ""
	for _, v := range result {
		if v == pkg { // do not print if not new
			v = ""
		} else {
			pkg = v
		}

		fmt.Printf("    %-40s %s\n",
			v,
			filepath.Base(v))
	}

}

func (c *Command) AddRemote() {

	if len(c.Flag.Args()) != 2 {
		fmt.Printf("Illegal arguments\n")
		return
	}
	name, remote := c.Flag.Arg(0), c.Flag.Arg(1)
	u, err := url.Parse(remote)
	if err != nil {
		ErrorStyle.Printf("Invalid URL passed as a remote Repository.\n    Caused by %s\n", err)
		return
	}
	c.Repository.RemoteAdd(NewRemoteRepository(name, *u))
	c.Repository.Write()
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
		c.Repository.RemoteRemove(name)
	}
	c.Repository.Write()
}