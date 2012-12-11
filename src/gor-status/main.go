package main 

import (
	"go.ericaro.net/gor"
	"fmt"
)

func main() {
	p,_:= gor.ReadProject()
	fmt.Printf("%v\n", p)
}

