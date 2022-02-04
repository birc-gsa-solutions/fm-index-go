package main

import (
	"fmt"
	"os"

	// Directories in the root of the repo can be imported
	// as long as we pretend that they sit relative to the
	// url birc.au.dk/gsa, like this for the example 'shared':
	"birc.au.dk/gsa/shared"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-p genome")
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "genome reads")
		os.Exit(1)
	}
	if os.Args[1] == "-p" {
		// preprocess
		fmt.Println(shared.TodoPreprocess(os.Args[2]))
	} else {
		fmt.Println(shared.TodoMap(os.Args[1], os.Args[2]))
	}

}
