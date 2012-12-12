package cmd

import (
	"fmt"
	"got.ericaro.net/got"
)

// contain shared operations between commands

func PrintVersion() {
	fmt.Println("Version:", got.GotVersion)
}
