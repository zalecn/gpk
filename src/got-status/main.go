package main 

import (
	"got.ericaro.net/got"
	"fmt"
)

func main() {
	p,_:= got.ReadProject()
	fmt.Printf("%v\n", p)
}

