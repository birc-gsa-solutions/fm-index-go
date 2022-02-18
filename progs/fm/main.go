package main

import (
	"fmt"
	"os"

	"birc.au.dk/gsa"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-p genome")
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "genome reads")
		os.Exit(1)
	}
	if os.Args[1] == "-p" {
		gsa.BwtPreproc(os.Args[2])
	} else {
		genome := gsa.ReadPreprocTables(os.Args[1])
		gsa.ScanFastq(os.Args[2], func(rec *gsa.FastqRecord) {
			for chrName, search := range genome {
				search(rec.Read, func(i int32) {
					cigar := fmt.Sprintf("%d%s", len(rec.Read), "M")
					gsa.PrintSam(rec.Name, chrName, i, cigar, rec.Read)
				})
			}
		})
	}
}
