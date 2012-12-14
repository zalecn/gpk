package main

import (
	"ericaro.net/gopackage"
	"fmt"
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
	UsageLine: `<dependencies>...`,
	Short:     `Add dependencies to this project.`,
	Long: `Dependencies are formatted as follow
	<group>:<name>:<version>
	where   :
	group   : is a simple name ( usually DNS name that identify a group )
	artifact: is a simple name (usually not hierarchical )
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
	UsageLine: `<dependencies>...`,
	Short:     `Remove dependencies to this project.`,
	Long: `Dependencies are formatted as follow
	<group>:<name>:<version>
	where   :
	group   : is a simple name ( usually DNS name that identify a group )
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
	UsageLine: `group:name`,
	Short:     `Init the current directory as a gopackage project.`,
	Long: `where   :
	group   : is a simple name ( usually DNS name that identify a group )
	artifact: is a simple name (usually not hierarchical )`,
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
	c.Project = p // in case 
	gopackage.WriteProject(c.Project)

}

func (c *Command) Add() {

	if len(c.Flag.Args()) == 0 {
		c.Flag.Usage()
		return
	}
	for _, v := range c.Flag.Args() {
		ref, err := gopackage.ParseProjectReference(v)
		if err != nil {
			fmt.Printf("  -X-> %v : fail to parse reference: %v\n", v, err)
		} else {
			fmt.Printf("  -> %v\n", v)
			c.Project.AppendDependency(ref)
		}

	}
	gopackage.WriteProject(c.Project)
}

func (c *Command) Remove() {

	if len(c.Flag.Args()) == 0 {
		c.Flag.Usage()
		return
	}
	for _, v := range c.Flag.Args() {
		ref, err := gopackage.ParseProjectReference(v)
		if err != nil {
			fmt.Printf("  -X-> %v : fail to parse reference: %v\n", v, err)
		} else {
			fmt.Printf("  -> %v\n", v)
			c.Project.RemoveDependency(ref)
		}
	}
	gopackage.WriteProject(c.Project)
}
