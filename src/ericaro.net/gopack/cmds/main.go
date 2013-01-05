//Package cmds contains the commands available in gpk executable.
package cmds

import (
	. "ericaro.net/gopack"
	"flag"
	"fmt"
	"math"
	"os/user"
	"path/filepath"
	"sort"
	"os"
	"errors"
	"log"
	"io/ioutil"
)

const (
	Cmd               = "gpk"
	GopackageVersion  = "1.0.0-beta.3" //?
	DefaultRepository = ".gpkrepository"
	

	RemoteCategory = -iota
	DependencyCategory = -iota
	CompileCategory = -iota
	InitCategory = -iota
	HelpCategory = -iota
	
)

// here are gopackage flags not specific ones

var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var localRepositoryFlag *string = flag.String("local", DefaultRepository, "path to the local repository to be used by default.")
var verboseFlag *bool = flag.Bool("verbose", false, "print verbose output.")

// We keep a dict AND a list of all available commands, the main command being generic
var Commands map[string]*Command = make(map[string]*Command)

// handle a sorted set of commands too
type commands []*Command //
var AllCommands commands = make([]*Command, 0)

func (s commands) Len() int      { return len(s) }
func (s commands) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s commands) Less(i, j int) bool {
	ci, cj := s[i], s[j]
	if ci.Category != cj.Category {
		return ci.Category < cj.Category
	}
	return s[i].Name < s[j].Name
}



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


func InvalidArgumentSize() error {
	return errors.New("Invalid Argument Size")
}

func PrintGlobalUsage() {
	TitleStyle.Printf("\n\nNAME\n")

	fmt.Printf("       gpk - Gopack is a software dependency management tool for Golang.\n             It help Managing, Building, and Sharing libraries in Go.\n")
	TitleStyle.Printf("\nSYNOPSIS\n")
	fmt.Printf("       %s [general options] <command> [options]  \n", Cmd)
	TitleStyle.Printf("\nOPTIONS\n")
	TitleStyle.PrintTriple("option", "default", "usage")

	flag.VisitAll(func(f *flag.Flag) {
		NormalStyle.PrintTriple("-"+f.Name, f.DefValue, f.Usage)
	})
	fmt.Println()
	TitleStyle.Printf("\nCOMMANDS\n")
	var category int8 = math.MinInt8
	for _, c := range AllCommands {
		if c.Category != category {
			category = c.Category
			fmt.Println()
		}
		NormalStyle.PrintTriple(c.Alias, c.Name, c.Short)
	}
	
	SuccessStyle.Printf("\n\n       Type 'gpk help [COMMAND]' for more details about a command.")
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
	if ! *verboseFlag {
		log.SetOutput(ioutil.Discard)
	}
	
	
	sort.Sort(AllCommands)
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
	cmd.Flag.Usage = func(){PrintCommandHelp(cmd) }
	// now continue parsing the command's args, using the command flags
	err = cmd.Flag.Parse(flag.Args()[1:])
	if err != nil {
		ErrorStyle.Printf("Cannot parse command line. %s\n", err)
		return
	}
	err = cmd.Run(cmd) // really execute the command
	if err != nil {
		os.Exit( -1 )
	}
}

//NewDefaultepository is the factory for a local repo. It tries to find one in the user's home dir. The full policy is defined here. 
func NewDefaultRepository() (r *LocalRepository, err error) {
	path := *localRepositoryFlag
	if !filepath.IsAbs(path) {
		u, _ := user.Current()
		path = filepath.Join(u.HomeDir, *localRepositoryFlag)
		path = filepath.Clean(path)
	}
	return NewLocalRepository(path)
}

//Command contains mainly declarative info about a specific command, and pointer to a function in charge of executing the command. 
// as command definitions are static within the code, there is no need to pass the command to the function, it already known it.
type Command struct {
	Category                            int8
	Run                                 func(c *Command) error // the callable to run
	FlagInit                            func(c *Command) // the callable to init flags
	Name, Alias, UsageLine, Short, Long string
	Flag                                flag.FlagSet // the command options to be parsed
	RequireProject                      bool
	Project                             *Project         // the project if required project as true, or a place holder for the command to write the current project
	Repository                          *LocalRepository // the local repository
}
