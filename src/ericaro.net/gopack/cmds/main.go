package cmds

import (
	. "ericaro.net/gopack"
	"flag"
	"fmt"
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

// We keep a dict AND a list of all available commands, the main command being generic
var Commands map[string]*Command = make(map[string]*Command)
var AllCommands []*Command = make([]*Command, 0)

//Reg is to register a command (or a bunch of them) to be available in the main
func Reg(commands ...*Command) {
	for _, c := range commands { // we append the command in the double map (name, and alias)
		Commands[c.Name] = c
		Commands[c.Alias] = c
		if c.FlagInit != nil {
			c.FlagInit(c)
		}
	}
	AllCommands = append(AllCommands, commands...)
}

// every file containing command shall register it this way
func init() {
	Reg(
		&Help,
	)
}

//Help Command
var Help = Command{
	Name:      `help`,
	Alias:     `h`,
	UsageLine: `[COMMAND]`,
	Short:     `Display help information about COMMAND`,
	Long:      ``, // better nothing than repeat
	Run: func(Help *Command) {

		if len(Help.Flag.Args()) == 0 {
			PrintGlobalUsage()
			return
		}
		cmdName := Help.Flag.Arg(0)
		cmd, ok := Commands[cmdName]
		if !ok {
			ErrorStyle.Printf("Unknown command %v.\n", cmdName)
			PrintGlobalUsage()
			return
		}
		TitleStyle.Printf("\nNAME\n\n")
		fmt.Printf("    gpk %s  - %s\n", cmd.Name, cmd.Short)
		TitleStyle.Printf("\nSYNOPSIS\n\n")
		fmt.Printf("    gpk %s %s\n", cmd.Name, cmd.UsageLine)
		
		TitleStyle.Printf("\nOPTIONS\n\n")
		TitleStyle.Printf("    %-10s  %-20s %s\n", "option", "default", "usage")
		cmd.Flag.VisitAll(printFlag)
		TitleStyle.Printf("\n\nDESCRIPTION\n\n")
		fmt.Print("    "+cmd.Long)
		fmt.Println("\n")
	},
}

func printFlag(f *flag.Flag) {
	fmt.Printf("    -%-10s %-20s %-s\n", f.Name, f.DefValue, f.Usage)

}

func PrintGlobalUsage() {
	TitleStyle.Printf("\n\nNAME\n\n")

	fmt.Printf("  gpk - Gopack is a software project management tool for Golang.\n")
	TitleStyle.Printf("\nSYNOPSIS\n\n")
	fmt.Printf("  %s [general options] <command> [options]  \n", Cmd)
	TitleStyle.Printf("\nOPTIONS\n\n")
	TitleStyle.Printf("    %-10s  %-20s %-s\n", "option", "default", "usage")
	flag.VisitAll(printFlag)
	fmt.Println()
	TitleStyle.Printf("\nCOMMANDS\n\n")
	for _, c := range AllCommands {
		fmt.Printf("  %-8s %-10s %s\n", c.Alias, c.Name, c.Short)
	}
	fmt.Println("\n")
}

//Gopack is the main. the real function main lies outside to create an executable
// It is fairly generic wrt to Commands, it first parses the general commands, then the command
func Gopack() {

	flag.Parse() // Scan the main arguments list
	if *versionFlag {
		fmt.Println("Version:", GopackageVersion)
		return
	}
	if len(flag.Args()) == 0 {
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

	// always allocate a local repository (create one if required)
	r, err := NewDefaultRepository()
	if err != nil {
		ErrorStyle.Printf("Cannot initialize the default repository. %s\n", err)
		return
	}
	cmd.Repository = r

	if cmd.RequireProject { // Commands can require to be executed on a project, in which case we have to load it
		p, err := ReadProject()
		if err != nil {
			ErrorStyle.Printf("Cannot initialize the current project. %s\n", err)
			return
		}
		cmd.Project = p
	}

	// now continue parsing the command's args, using the command flags
	err = cmd.Flag.Parse(flag.Args()[1:])
	if err != nil {
		ErrorStyle.Printf("Cannot parse command line. %s\n", err)
		return
	}
	cmd.Run(cmd) // really execute the command
}

//NewDefaultepository is the factory for a local repo. It tries to find one in the user's home dir. The full policy is defined here. 
func NewDefaultRepository() (r *LocalRepository, err error) {
	u, _ := user.Current()
	path := filepath.Join(u.HomeDir, *localRepositoryFlag)
	path = filepath.Clean(path)
	return NewLocalRepository(path)
}

//Command contains mainly declarative info about a specific command, and pointer to a function in charge of executing the command. 
// as command definitions are static within the code, there is no need to pass the command to the function, it already known it.
type Command struct {
	Run                                 func(c *Command) // the callable to run
	FlagInit                            func(c *Command) // the callable to init flags
	Name, Alias, UsageLine, Short, Long string
	Flag                                flag.FlagSet // the command options to be parsed
	RequireProject                      bool
	Project                             *Project         // the project if required project as true, or a place holder for the command to write the current project
	Repository                          *LocalRepository // the local repository
}
