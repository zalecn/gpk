package main

import (
	"fmt"
	"got.ericaro.net/got"
	"strings"
)

func init() {
	Reg(
		&Status,
		&Add,
		&Init,
	)

}

var Status = Command{
	Name:           `status`,
	UsageLine:      ``,
	Short:          `Prints current directory Got status.`,
	Long:           `Prints current directory status`,
	call:           func(c *Command) { c.Status() },
	RequireProject: true,
}

var Add = Command{
	Name:      `add`,
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

var Init = Command{
	Name:      `init`,
	UsageLine: `group:name`,
	Short:     `init the current directory as a got project.`,
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

	p, err := got.ReadProject()
	if err == nil {
		fmt.Printf("warning: init an existing project. This is fine if you want to edit it")
	}

	parts := strings.Split(c.Flag.Arg(0), ":")
	if len(parts) != 2 {
		fmt.Printf("Invalid Project Reference Format, must be group:name \n")
		return
	}
	p.Group = parts[0]
	p.Artifact = parts[1]

	c.Project = p // in case 
	got.WriteProject(c.Project)

}

func (c *Command) Add() {

	if len(c.Flag.Args()) == 0 {
		c.Flag.Usage()
		return
	}
	for _, v := range c.Flag.Args() {
		ref, err := got.ParseProjectReference(v)
		if err != nil {
			fmt.Printf("  -X-> %v : fail to parse reference: %v\n", v, err)
		} else {
			fmt.Printf("  -> %v\n", v)
			c.Project.AppendDependency(ref)
		}

	}
	got.WriteProject(c.Project)
}
