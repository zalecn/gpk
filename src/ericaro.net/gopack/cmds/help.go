package cmds

import (
	. "ericaro.net/gopack"
	"flag"
	"fmt"
)

func init() {
	Reg(
		&Help,
	)
}

//Help Command
var Help = Command{
	Name:      `help`,
	Alias:     `h`,
	Category:  HelpCategory,
	UsageLine: `[COMMAND]`,
	Short:     `Display help information about commands`,
	Long:      `Display general inline help, or help information about COMMAND`, // better nothing than repeat
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
		TitleStyle.Printf("\nNAME\n")
		fmt.Printf("       gpk %s  - %s\n", cmd.Name, cmd.Short)
		TitleStyle.Printf("\nSYNOPSIS\n")
		fmt.Printf("       gpk %s %s\n", cmd.Name, cmd.UsageLine)
		var flags uint
		cmd.Flag.VisitAll(func(f *flag.Flag) {
			flags += 1
		})
		if flags > 0 {
			TitleStyle.Printf("\nOPTIONS\n")
			TitleStyle.PrintTriple("option", "default", "usage")
			cmd.Flag.VisitAll(func(f *flag.Flag) {
				NormalStyle.PrintTriple("-"+f.Name, f.DefValue, f.Usage)
			})
		}
		if len(cmd.Long) > 0 {

			TitleStyle.Printf("\n\nDESCRIPTION\n")
			fmt.Print("       " + cmd.Long)
		}
		fmt.Println("\n")
	},
}
