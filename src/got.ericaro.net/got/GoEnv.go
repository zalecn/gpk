package got

import (
	"os/exec"
	"os"
	"fmt"
)

// reflect some basic go operations

type GoEnv struct {
	gopath string
}




func NewGoEnv(gopath string) *GoEnv{
	return &GoEnv{
		gopath: gopath,
	}
	
}

func (g *GoEnv) Install(root string) {

	cmd := exec.Command("go", "install", "./src/...")
	
	env := []string{
	fmt.Sprintf("GOPATH=%s:%s", g.gopath, root),
	}
	
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir  = root // asbolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}
} 

// plan :
// script (on the server) that scans for godocs for instance ) to auto upload 1.0 of every product
// offer snapshot server instances, search the code (google code search), browse the doc

// tasks
// implements a server simple front end for artifacts  (name/ link to their source)

 


