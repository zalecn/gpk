package got

import (
	"os/exec"
	"os"
	"fmt"
)

// reflect some basic go operations

type GoEnv struct {
	env []string
}


func NewGoEnv(gopath string) *GoEnv{
	env := make([]string, 0,10)
	env = append(env, fmt.Sprintf("GOPATH=%s", gopath)  )
	return &GoEnv{
		env: env,
	}
	
}

func (g *GoEnv) Install(root string) {

	cmd := exec.Command("go", "install", "./src/...")
	cmd.Env = g.env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir  = root // asbolute path of the project
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}
} 

// plan :
// server for .gots database with .got table ( g, a , b, v, blob= tar.gz )
// deploy/download from server.
// script (on the server) that scans for godocs for instance ) to auto upload 1.0 of every product
// offer snapshot server instances, search the code (google code search), browse the doc

// tasks
// implement a download cmd ( http tar.gz -> install dir header contains the metadata ? right ?)
// implement an upload cmd ( reverse protocol )
// start  a server in go and google app engine, empty engine
// implements a server, and receive http -> blob
// implements a server simple front end for artifacts  (name/ link to their source)
// 

 


