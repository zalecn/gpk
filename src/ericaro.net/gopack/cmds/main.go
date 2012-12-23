package cmds

import (
	. "ericaro.net/gopack"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	Cmd               = "gpk"
	GopackageVersion  = "0.0.1" //?
	DefaultRepository = ".gpkrepository"
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
	UsageLine: `[COMMAND]`,
	Short:     `Display help information about COMMAND`,
	Long:      ``, // better nothing than repeat
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


	

	fmt.Printf("\nusage: %s\n", TitleStyle.Sprintf("gpk %s %s",cmd.Name, cmd.UsageLine) )
	fmt.Printf("    %s\n", ShortStyle.Sprintf("%s", cmd.Short) )
	fmt.Printf("where:\n")
	TitleStyle.Printf("    %-10s %-20s %s\n", "option", "default", "usage")
	cmd.Flag.VisitAll(printFlag)
	fmt.Print(cmd.Long)
}

func printFlag( f *flag.Flag ) {
	fmt.Printf("    -%-10s %-20s %s\n", f.Name, f.DefValue, f.Usage)
	
}


func PrintGlobalUsage() {
	TitleStyle.Printf("\nGopack is a software project management tool for Golang.\n") 
	fmt.Printf("\nusage: ")
	TitleStyle.Printf("%s [general options] <command> [options]  \n", Cmd)

	fmt.Printf("Where general options are:\n")
	TitleStyle.Printf("    %-10s %-20s %s\n", "option", "default", "usage")
	flag.VisitAll(printFlag)
	fmt.Println()
	fmt.Printf("Where <command> are:\n")

	fmt.Printf("  %-8s %-10s %s\n", "alias", "name", "description")
	fmt.Printf("  -------------------\n")

	for _, c := range AllCommands {
		fmt.Printf("  %-8s %-10s %s\n", c.Alias, c.Name, c.Short)
	}
}


func Gopack() {

	
	flag.Parse() // Scan the main arguments list
	if *versionFlag {
		fmt.Println("Version:", GopackageVersion)
		return
	}
	if len(flag.Args() ) == 0 {
		PrintGlobalUsage()
		return
	}
	cmdName := flag.Arg(0)
	cmd, ok := Commands[cmdName]
	if !ok {
		ErrorStyle.Printf("Unknown command %v\n", cmdName)
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
