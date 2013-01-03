package cmds

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/gocmd"
	"fmt"
)

func init() {
	Reg(
		&ListDependencies,
		&ListRemotes,
		&Path,
		&Imports,
	)

}

var importsOfflineFlag *bool
var importsAutofixFlag *bool
var Imports = Command{
	Name:           `list-missing`,
	Alias:          `lm`,
	Category:       DependencyCategory,
	UsageLine:      ``,
	Short:          `Analyse the current directory and report or fix missing dependencies`,
	Long:           ``,
	RequireProject: true,
	FlagInit: func(Imports *Command) {
		importsOfflineFlag = Imports.Flag.Bool("o", false, "offline. Do not Use remotes while looking for dependencies")
		importsAutofixFlag = Imports.Flag.Bool("f", false, "fix auto. Auto fix the current project using default choices")

	},
	Run: func(Imports *Command)  (err error){
		toSave := false
		missing := Imports.Repository.MissingImports(Imports.Project, *importsOfflineFlag)
		missingPack := Imports.Repository.MissingPackages(missing)
		SuccessStyle.Printf("Missing imports (%d), missing packages (%d)\n", len(missing), len(missingPack))
		for _, m := range missingPack {
			found := Imports.Repository.ImportSearch(m)
			if len(found) > 0 {
				SuccessStyle.Printf("Missing packages %-40s -> ☑ %s \n", m, found[0])
				for _, pid := range found[1:] {
					SuccessStyle.Printf("                 %-40s -> ☐ %s \n", "", pid)
				}
				if *importsAutofixFlag {
					Imports.Project.AppendDependency(found[0])
					toSave = true
				}
			}
		}

		if toSave {
			SuccessStyle.Printf("Project Updated\n")
			Imports.Project.Write()
		}
		return
	},
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var ListDependencies = Command{
	Name:           `list-dependencies`,
	Alias:          `ld`,
	Category:       DependencyCategory,
	UsageLine:      ``,
	Short:          `List declared Dependencies.`,
	Long:           `List declared dependencies.`,
	RequireProject: true,
	FlagInit: func(ListDependencies *Command) {
	},
	Run: func(ListDependencies *Command)  (err error){
		TitleStyle.Printf("\nLIST OF DECLARED DEPENDENCIES:\n")
		// TODO print in a suitable way for copy pasting
		dependencies := ListDependencies.Project.Dependencies()
		for _, d := range dependencies {
			SuccessStyle.Printf("        %-40s %s\n", d.Name(), d.Version().String())
		}
		return
	},
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var ListRemotes = Command{
	Name:           `list-remotes`,
	Alias:          `lr`,
	Category:       RemoteCategory,
	UsageLine:      ``,
	Short:          `List Remotes.`,
	Long:           `List declared remotes.`,
	RequireProject: false,
	Run: func(ListRemotes *Command)  (err error){
		TitleStyle.Printf("\nLIST OF REMOTES:\n")
		rem := ListRemotes.Repository.Remotes()
		if len(rem) == 0 {
			SuccessStyle.Printf("       <empty>\n")
		} else {
			for _, r := range rem {
				u := r.Path()
				tr := "" // str repr of the token
				t := r.Token()
				if t != nil { // applies only if not nul
					tr = t.FormatStd()
				}

				SuccessStyle.Printf("       %-8s %-40s %s\n", r.Name(), u.String(), tr)
			}
		}
		return
	},
}

var pathListFlag *bool
var Path = Command{
	Name:      `list-package`,
	Alias:     `lp`,
	Category:  DependencyCategory,
	UsageLine: ``,
	Short:     `List all packages dependencies (recursive)`,
	Long: `Resolve current project dependencies and print result.
       
       ` + "Tip:\n           type:\n           alias GP='export GOPATH=`gpk lp`'\n           to get a simple automatic GOPATH setter.",
	RequireProject: true,
	FlagInit: func(Path *Command) {
		pathListFlag = Path.Flag.Bool("l", false, "list. Pretty Print the list.")
	},
	Run: func(Path *Command)  (err error){

		// parse dependencies, and build the gopath
		dependencies, err := Path.Repository.ResolveDependencies(Path.Project, true, false) // path does not update the dependencies
		if err != nil {
			ErrorStyle.Printf("Error Resolving project's dependencies:\n    \u21b3 %v", err)
			return
		}

		if *pathListFlag {
			TitleStyle.Printf("\nLIST OF PACKAGES:\n")
			// run the go build command for local src, and with the appropriate gopath
			if len(dependencies) > 0 {
				for _, d := range dependencies {
					SuccessStyle.Printf("        %-40s %s\n", d.Name(), d.Version().String())
				}
			}else {
				SuccessStyle.Printf("       <empty>\n")
			}
			fmt.Println()
		} else {
			gopath, _ := Path.Repository.GoPath(dependencies)
			fmt.Print(gocmd.Join(Path.Project.WorkingDir(), gopath)) // there is no line break because this way we can use it to initialize a local GOPATH var
		}
		return
	},
}
