package gsa

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func BwtPreproc(genomeFileName string) {

	genome := LoadFasta(genomeFileName)
	preprocFileName := genomeFileName + ".fmidx"

	outFile, err := os.Create(preprocFileName)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s", preprocFileName, err)
	}

	processed := map[string]*FMIndexTables{}
	for name, seq := range genome {
		// We always preprocess the full set, even if it is
		// slower if we intent to do exact matching
		processed[name] = BuildFMIndexExactTables(seq)
	}

	enc := gob.NewEncoder(outFile)
	if err := enc.Encode(processed); err != nil {
		log.Fatal("encode error:", err)
	}

	if err := outFile.Close(); err != nil {
		log.Fatalf("Error closing file %s: %s", preprocFileName, err)
	}
}

func ReadPreprocTables(fname string) map[string]func(p string, cb func(i int32)) {
	infile, err := os.Open(fname + ".fmidx")
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Couldn't open file: %s, did you remember to preprocess?",
			fname)
		os.Exit(1)
	}

	genomeTables := map[string]*FMIndexTables{}
	dec := gob.NewDecoder(infile)

	if err := dec.Decode(&genomeTables); err != nil {
		log.Fatalf("Error decoding preprocessing file %s: %s",
			fname, err)
	}

	infile.Close()

	searchFuncs := map[string]func(p string, cb func(i int32)){}
	for rname, tbls := range genomeTables {
		searchFuncs[rname] = FMIndexExactFromTables(tbls)
	}

	return searchFuncs
}
