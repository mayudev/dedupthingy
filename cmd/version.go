package cmd

import "fmt"

var Version = "0.0.1"

func runVersion() {
	fmt.Println("dedupthingy version", Version)
}
