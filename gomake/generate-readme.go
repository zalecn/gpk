package main

import (
	. "ericaro.net/gopack"
	"ericaro.net/gopack/cmds"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Rst2Markdown(root, src, dst string) {

	cmd := exec.Command("pandoc", "-f", "rst", "-t", "markdown", "-o", dst, src)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = root // absolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}

}

func Rst2Man(root, src, dst string) {

	cmd := exec.Command("rst2man", src, dst)

	//cmd.Env = BuildEnv(locals)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = root // absolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}

}

func Rst2Html(root, src, dst string) {

	cmd := exec.Command("rst2html", "", src, dst)

	//cmd.Env = BuildEnv(locals)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = root // absolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}

}

func main() {
	
	
	Rst2Markdown(".", "README.rst", "README.md")
	
	
	// generate the summary of all commands
//	os.MkdirAll(filepath.Join("./target"), os.ModeDir|os.ModePerm) // mkdir -p
//
//	filename := filepath.Join("./target/readme.rst")
//	f, err := os.Create(filename)
//	if err != nil {
//		ErrorStyle.Printf("Cannot create file %s. Due to %v\n", filename, err)
//	}
//	defer f.Close()
//	
//	s, err := os.Open("./doc/gpk.rst")
//	defer s.Close()
//	io.Copy(f, s)
//	
//	for _, c := range cmds.AllCommands {
//		CreateCmdSummary(f, c)
//		//		CreateCmdMan(c)
//	}
//	
//
//	fmt.Fprintf(f, "\n")

}

func CreateCmdSummary(f io.Writer, c *cmds.Command) {
	fmt.Fprintf(f,
		`gpk %s %s
    %s
    
    

`, c.Alias, c.UsageLine, c.Short)

}

//
//func CreateCmdMan(c *cmds.Command) {
//	filename := filepath.Join("./target/", "gpk-"+c.Name+".rst") // generate only the start of the command
//	man := filepath.Join("./target/man1", "gpk-"+c.Name+".1")        // generate only the start of the command
//	f, err := os.Create(filename)
//	if err != nil {
//		ErrorStyle.Printf("Cannot create file %s. Due to %v\n", filename, err)
//	}
//	defer f.Close()
//	defer Rst2Man(".", filename, man)
//	fmt.Fprintf(f,
//		`=====================================
//gpk-%s
//=====================================
//%s
//----------------------------------------------------------------------------------------------------
//
//:Author: eric@ericaro.net
//:Date:   20012-12-22
//:Copyright: GPL
//:Version: 0.1
//:Manual section: 1
//:Manual group: compiler
//`, c.Name, c.Short)
//
//	fmt.Fprintf(f, `
//SYNOPSIS
//===========
//
//gpk %s %s
//
//`, c.Name, c.UsageLine)
//	if c.Alias != c.Name {
//		fmt.Fprintf(f, `
//or
//
//gpk %s %s
//`, c.Alias, c.UsageLine)
//	}
//
//	fmt.Fprintf(f, `
//%s
//
//`, c.Long)
//
//	count := 0
//	i := &count
//	counter := func(flag *flag.Flag) { (*i)++ }
//	c.Flag.VisitAll(counter)
//	if *i > 0 {
//
//		fmt.Fprintf(f, `
//OPTIONS
//=======
//
//`)
//		rstFlag := func(flag *flag.Flag) {
//			dash := "-"
//			if len(flag.Name) > 1 {
//				dash = "--"
//			}
//			fmt.Fprintf(f, "%s%s %s  %s\n", dash, flag.Name, flag.DefValue, flag.Usage)
//		}
//		c.Flag.VisitAll(rstFlag)
//	}
//
//	fmt.Fprintln(f)
//}

//I idents a block by n spaces
func I(block string, n int) string {
	return strings.Replace(block, "\n", "\n"+strings.Repeat(" ", n), -1)
}
