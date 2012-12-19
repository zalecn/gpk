package main

import (
	. "ericaro.net/gpk"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

const (
	Cmd               = "gpk"
	GopackageVersion  = "0.0.1" //?
	DefaultRepository = ".gpkrepository"
)

var (
	GopackageCentral, _ = url.Parse("http://gpk.ericaro.net")
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

func init() {
	Reg(
		&Help,
	)
}

var Help = Command{
	Name:      `help`,
	Alias:     `h`,
	UsageLine: `<command>`,
	Short:     `Display more advanced help for the given command`,
	Long:      `Display more advanced help for the given command`,
	call:      func(c *Command) { c.Help() },
}

func (c *Command) Help() {

	if len(c.Flag.Args()) == 0 {
		PrintGlobalUsage()
		return
	}
	cmdName := c.Flag.Arg(0)
	cmd, ok := Commands[cmdName]
	if !ok {
		fmt.Printf("Unknown command %v. Available commands are:\n\n", cmdName)
		PrintGlobalUsage()
		return
	}
	
	fmt.Printf( //TODO beautify with console colors
`gpk %s %s

%s

%s
`, cmd.Name, cmd.UsageLine, cmd.Short, cmd.Long )

}

func PrintGlobalUsage() {
	fmt.Printf("%s [general options] <alias|name> [options]  \n", Cmd)
	fmt.Printf("  %-8s %-10s %s\n", "alias", "name", "description")
	fmt.Printf("  -------------------\n")

	for _, c := range AllCommands {
		fmt.Printf("  %-8s %-10s %s\n", c.Alias, c.Name, c.Short)
	}
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
		PrintGlobalUsage()

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
