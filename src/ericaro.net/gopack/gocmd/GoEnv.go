package gocmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// reflect some basic go operations

type GoEnv struct {
	gopath string
}

func NewGoEnv(gopath string) *GoEnv {
	return &GoEnv{
		gopath: gopath,
	}
}

func BuildEnv(vals map[string]string) []string {
	current := os.Environ()
	newenv := make([]string, 0, len(current))
	for _, v := range current {
		parts := strings.SplitN(v, "=", 2)
		k := parts[0]
		if val, ok := vals[k]; ok { // overwrite it
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, val))
			vals[k] = "" // mark it has deleted
		} else {
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, parts[1]))
		}
	}
	// now append the "new" ones
	for k, val := range vals {
		if len(val) != 0 {
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, val))
		}

	}

	return newenv
}

func Join(path string, elements ...string) string {
	if len(path) == 0 {
		return strings.Join(elements, string(os.PathListSeparator))
	} else {
		return path + string(os.PathListSeparator) + strings.Join(elements, string(os.PathListSeparator))
	}
	panic("unreachable statement")

}

func (g *GoEnv) Install(root string, all bool) {
	var cmd *exec.Cmd
	if all{
	cmd = exec.Command("go", "install", "-a",  "./src/...")
	} else {
		cmd = exec.Command("go", "install",  "./src/...")
	}
	
	locals := map[string]string{
		"GOPATH": Join(root, g.gopath),
	}

	cmd.Env = BuildEnv(locals)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = root // asbolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (g *GoEnv) Test(root string) {

	cmd := exec.Command("go", "test",  "./src/...")

	locals := map[string]string{
		"GOPATH": Join(g.gopath, root),
	}
	
	cmd.Env = BuildEnv(locals)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = root // asbolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (g *GoEnv) Get(pack string) {

	cmd := exec.Command("go", "get", pack)

	locals := map[string]string{
		"GOPATH": g.gopath,
	}
	//fmt.Printf("GOPATH = %v\n", g.gopath)
	cmd.Env = BuildEnv(locals)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = g.gopath // asbolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

// plan :
// script (on the server) that scans for godocs for instance ) to auto upload 1.0 of every product
// offer snapshot server instances, search the code (google code search), browse the doc

// tasks
// implements a server simple front end for artifacts  (name/ link to their source)
