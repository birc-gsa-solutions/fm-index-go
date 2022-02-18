package gsa

import (
	"fmt"
)

func printSam(rname, read, chrom string, pos int) {
	fmt.Printf("%s\t%s\t%d\t%dM\t%s\n", rname, chrom, pos+1, len(read), read)
}

type SearchAlgorithm func(x, p string, callback func(hit int))

func SearchGenome(genomeFname, readsFname string, alg SearchAlgorithm) {
	genome := LoadFasta(genomeFname)
	ScanFastq(readsFname, func(rec *FastqRecord) {
		for chromName, chromSeq := range genome {
			alg(chromSeq, rec.Read, func(pos int) {
				printSam(rec.Name, rec.Read, chromName, pos)
			})
		}
	})
}
