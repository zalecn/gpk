package gor

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
	
	fmt.Printf("%v > %v %v  where %v\n", cmd.Dir, cmd.Path, cmd.Args, cmd.Env)
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v\n", err)
	}

} 

