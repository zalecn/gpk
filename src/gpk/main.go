package main

import (
	. "ericaro.net/gpk"
	"flag"
	"fmt"
	"os"
	"os/user"
	"net/url"
	"path/filepath"
)

const (
	Cmd               = "gpk"
	GopackageVersion  = "0.0.1" //?
	DefaultRepository = ".gpkrepository"
)

var (
	GopackageCentral,_  = url.Parse("http://gpk.ericaro.net")
)

// here are gopackage flags not specific ones

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var localRepositoryFlag *string = flag.String("local", DefaultRepository, "path to the local repository to be used by default.")


var Commands map[string]*Command = make(map[string]*Command)
var AllCommands []*Command = make([]*Command, 0)

func Reg(commands ...*Command) {
	for _, c := range commands {
		Commands[c.Name] = c
		Commands[c.Alias] = c
	}
	AllCommands = append(AllCommands, commands...)
}

func main() {

	flag.Parse() // Scan the main arguments list
	if *versionFlag {
		fmt.Println("Version:", GopackageVersion)
		return
	}

	cmdName := flag.Arg(0)
	cmd, ok := Commands[cmdName]
	if !ok {
		fmt.Printf("Unknown command %v. Available commands are:\n\n", cmdName)

		fmt.Printf("%s [general options] <alias|name> [options]  \n", Cmd)
		fmt.Printf("  %-8s %-10s %s\n", "alias", "name", "description")
		fmt.Printf("  -------------------\n")

		for _, c := range AllCommands {
			fmt.Printf("  %-8s %-10s %s\n", c.Alias, c.Name, c.Short)
		}
		return
	}

	r, err := NewDefaultRepository()
	
	handleError(err)
	
	cmd.Repository = r

	if cmd.RequireProject {
		p, err := ReadProject()
		handleError(err)
		cmd.Project = p
	}
	err = cmd.Flag.Parse(flag.Args()[1:])
	handleError(err)
	cmd.Run()

}

func NewDefaultRepository() (r *LocalRepository, err error) {
	u, _ := user.Current()
	path := filepath.Join(u.HomeDir, *localRepositoryFlag)
	path = filepath.Clean(path)
	fmt.Printf("local repo = %s\n", path)
	return NewLocalRepository(path)
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type Commander interface {
	Run()
}

type Command struct {
	call                                func(c *Command)
	Name, Alias, UsageLine, Short, Long string
	Flag                                flag.FlagSet
	RequireProject                      bool
	Project                             *Project
	Repository                          *LocalRepository
}

func (c *Command) Run() {
	c.call(c)
}

/* here collect cmd use case
gpk : display available version in local repo for the current repo


*/
