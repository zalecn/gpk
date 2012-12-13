package got

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

func (g *GoEnv) BuildEnv(vals map[string]string) []string {
	current := os.Environ()
	newenv := make([]string, 0, len(current))
	for _, v := range current {
		parts := strings.SplitN(v, "=", 2)
		k := parts[0]
		if val, ok := vals[k]; ok {
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, val))
		} else {
			newenv = append(newenv, fmt.Sprintf("%s=%s", k, parts[1]))
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

func (g *GoEnv) Install(root string) {

	cmd := exec.Command("go", "install", "./src/...")

	locals := map[string]string{
		"GOPATH": Join(g.gopath, root),
	}
	env := g.BuildEnv(locals)

	cmd.Env = env
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
	env := g.BuildEnv(locals)
	fmt.Println(env)
	cmd.Env = env
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
