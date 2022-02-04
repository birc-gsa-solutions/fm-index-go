package main

import (
	"fmt"

	// Directories in the root of the repo can be imported
	// as long as we pretend that they sit relative to the
	// url birc.au.dk/gsa, like this for the example 'shared':
	"birc.au.dk/gsa/shared"
)

func main() {
	fmt.Println(shared.Hello())
}
