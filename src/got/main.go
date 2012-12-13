package main

import (
	"flag"
	"fmt"
	"got.ericaro.net/got"
	"os"
)

// here are got flags not specific ones

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var hostFlag *string  = flag.String("host", got.GotCentral, "Set the host for the central server")



var Commands map[string]*Command = make(map[string]*Command)

func Reg(commands ...*Command) {
	for _,c:= range commands {
		Commands[c.Name] = c
	}
}


func main() {

	flag.Parse() // Scan the main arguments list
	if *versionFlag {
		fmt.Println("Version:", got.GotVersion)
		return
	}
	

	cmdName := flag.Arg(0)
	cmd := Commands[cmdName]
	if cmd == nil {
		fmt.Printf("Unknown command %v. Available commands are:\n", cmdName)
		for k, c := range Commands {
			fmt.Printf("  got %-10s %s\n", k, c.Short)
		}
		return
	}

	r, err := got.NewDefaultRepository()
	handleError(err)
	r.ServerHost = *hostFlag
	
	cmd.Repository = r
	
	if cmd.RequireProject {
		p, err := got.ReadProject()
		handleError(err)
		cmd.Project = p
	}
	err = cmd.Flag.Parse(flag.Args()[1:])
	handleError(err)
	cmd.Run()
	

}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type Commander interface{Run()}

type Command struct {
	
	call                    func(c *Command)
	Name,UsageLine, Short, Long string
	Flag                  flag.FlagSet
	RequireProject         bool
	Project                *got.Project
	Repository             *got.Repository
}

func (c *Command) Run() {
	c.call(c)
}

